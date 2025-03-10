// Copyright 2016-2021, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate go run generate.go

package pulumi

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"

	structpb "github.com/golang/protobuf/ptypes/struct"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"google.golang.org/grpc"
)

var disableResourceReferences = cmdutil.IsTruthy(os.Getenv("PULUMI_DISABLE_RESOURCE_REFERENCES"))

// Context handles registration of resources and exposes metadata about the current deployment context.
type Context struct {
	ctx         context.Context
	info        RunInfo
	stack       Resource
	exports     map[string]Input
	monitor     pulumirpc.ResourceMonitorClient
	monitorConn *grpc.ClientConn
	engine      pulumirpc.EngineClient
	engineConn  *grpc.ClientConn

	keepResources    bool // true if resources should be marshaled as strongly-typed references.
	keepOutputValues bool // true if outputs should be marshaled as strongly-type output values.

	rpcs     int        // the number of outstanding RPC requests.
	rpcsDone *sync.Cond // an event signaling completion of RPCs.
	rpcsLock sync.Mutex // a lock protecting the RPC count and event.
	rpcError error      // the first error (if any) encountered during an RPC.

	join workGroup // the waitgroup for non-RPC async work associated with this context

	Log Log // the logging interface for the Pulumi log stream.
}

// NewContext creates a fresh run context out of the given metadata.
func NewContext(ctx context.Context, info RunInfo) (*Context, error) {
	// Connect to the gRPC endpoints if we have addresses for them.
	var monitorConn *grpc.ClientConn
	var monitor pulumirpc.ResourceMonitorClient
	if addr := info.MonitorAddr; addr != "" {
		conn, err := grpc.Dial(
			info.MonitorAddr,
			grpc.WithInsecure(),
			rpcutil.GrpcChannelOptions(),
		)
		if err != nil {
			return nil, fmt.Errorf("connecting to resource monitor over RPC: %w", err)
		}
		monitorConn = conn
		monitor = pulumirpc.NewResourceMonitorClient(monitorConn)
	}

	var engineConn *grpc.ClientConn
	var engine pulumirpc.EngineClient
	if info.engineConn != nil {
		engineConn = info.engineConn
		engine = pulumirpc.NewEngineClient(engineConn)
	} else if addr := info.EngineAddr; addr != "" {
		conn, err := grpc.Dial(
			info.EngineAddr,
			grpc.WithInsecure(),
			rpcutil.GrpcChannelOptions(),
		)
		if err != nil {
			return nil, fmt.Errorf("connecting to engine over RPC: %w", err)
		}
		engineConn = conn
		engine = pulumirpc.NewEngineClient(engineConn)
	}

	if info.Mocks != nil {
		monitor = &mockMonitor{project: info.Project, stack: info.Stack, mocks: info.Mocks}
		engine = &mockEngine{}
	}

	supportsFeature := func(id string) (bool, error) {
		if monitor != nil {
			resp, err := monitor.SupportsFeature(ctx, &pulumirpc.SupportsFeatureRequest{Id: id})
			if err != nil {
				return false, fmt.Errorf("checking monitor features: %w", err)
			}
			return resp.GetHasSupport(), nil
		}
		return false, nil
	}

	keepResources, err := supportsFeature("resourceReferences")
	if err != nil {
		return nil, err
	}

	keepOutputValues, err := supportsFeature("outputValues")
	if err != nil {
		return nil, err
	}

	context := &Context{
		ctx:              ctx,
		info:             info,
		exports:          make(map[string]Input),
		monitorConn:      monitorConn,
		monitor:          monitor,
		engineConn:       engineConn,
		engine:           engine,
		keepResources:    keepResources,
		keepOutputValues: keepOutputValues,
	}
	context.rpcsDone = sync.NewCond(&context.rpcsLock)
	context.Log = &logState{
		engine: engine,
		ctx:    ctx,
		join:   &context.join,
	}
	return context, nil
}

// Close implements io.Closer and relinquishes any outstanding resources held by the context.
func (ctx *Context) Close() error {
	if ctx.engineConn != nil {
		if err := ctx.engineConn.Close(); err != nil {
			return err
		}
	}
	if ctx.monitorConn != nil {
		if err := ctx.monitorConn.Close(); err != nil {
			return err
		}
	}
	return nil
}

// wait waits for all asynchronous work associated with this context to drain. RPCs may not be queued once wait
// returns.
func (ctx *Context) wait() error {
	// Wait for async work to flush.
	ctx.join.Wait()

	// Ensure all outstanding RPCs have completed before proceeding. Also, prevent any new RPCs from happening.
	ctx.rpcsLock.Lock()
	defer ctx.rpcsLock.Unlock()

	// Wait until the RPC count hits zero.
	for ctx.rpcs > 0 {
		ctx.rpcsDone.Wait()
	}

	// Mark the RPCs flag so that no more RPCs are permitted.
	ctx.rpcs = noMoreRPCs

	if ctx.rpcError != nil {
		return fmt.Errorf("waiting for RPCs: %w", ctx.rpcError)
	}

	return nil
}

// Project returns the current project name.
func (ctx *Context) Project() string { return ctx.info.Project }

// Stack returns the current stack name being deployed into.
func (ctx *Context) Stack() string { return ctx.info.Stack }

// Parallel returns the degree of parallelism currently being used by the engine (1 being entirely serial).
func (ctx *Context) Parallel() int { return ctx.info.Parallel }

// DryRun is true when evaluating a program for purposes of planning, instead of performing a true deployment.
func (ctx *Context) DryRun() bool { return ctx.info.DryRun }

// GetConfig returns the config value, as a string, and a bool indicating whether it exists or not.
func (ctx *Context) GetConfig(key string) (string, bool) {
	v, ok := ctx.info.Config[key]
	return v, ok
}

// IsConfigSecret returns true if the config value is a secret.
func (ctx *Context) IsConfigSecret(key string) bool {
	for _, secretKey := range ctx.info.ConfigSecretKeys {
		if key == secretKey {
			return true
		}
	}
	return false
}

// Invoke will invoke a provider's function, identified by its token tok. This function call is synchronous.
//
// args and result must be pointers to struct values fields and appropriately tagged and typed for use with Pulumi.
func (ctx *Context) Invoke(tok string, args interface{}, result interface{}, opts ...InvokeOption) (err error) {
	if tok == "" {
		return errors.New("invoke token must not be empty")
	}

	resultV := reflect.ValueOf(result)
	if !(resultV.Kind() == reflect.Ptr &&
		(resultV.Elem().Kind() == reflect.Struct ||
			(resultV.Elem().Kind() == reflect.Map && resultV.Elem().Type().Key().Kind() == reflect.String))) {
		return errors.New("result must be a pointer to a struct or map value")
	}

	options := &invokeOptions{}
	for _, o := range opts {
		if o != nil {
			o.applyInvokeOption(options)
		}
	}

	// Note that we're about to make an outstanding RPC request, so that we can rendezvous during shutdown.
	if err = ctx.beginRPC(); err != nil {
		return err
	}
	defer ctx.endRPC(err)

	var providerRef string
	if provider := mergeProviders(tok, options.Parent, options.Provider, nil)[getPackage(tok)]; provider != nil {
		pr, err := ctx.resolveProviderReference(provider)
		if err != nil {
			return err
		}
		providerRef = pr
	}

	// Serialize arguments. Outputs will not be awaited: instead, an error will be returned if any Outputs are present.
	if args == nil {
		args = struct{}{}
	}
	resolvedArgs, _, err := marshalInput(args, anyType, false)
	if err != nil {
		return fmt.Errorf("marshaling arguments: %w", err)
	}

	resolvedArgsMap := resource.PropertyMap{}
	if resolvedArgs.IsObject() {
		resolvedArgsMap = resolvedArgs.ObjectValue()
	}

	rpcArgs, err := plugin.MarshalProperties(
		resolvedArgsMap,
		ctx.withKeepOrRejectUnknowns(plugin.MarshalOptions{
			KeepSecrets:   true,
			KeepResources: ctx.keepResources,
		}),
	)
	if err != nil {
		return fmt.Errorf("marshaling arguments: %w", err)
	}

	// Now, invoke the RPC to the provider synchronously.
	logging.V(9).Infof("Invoke(%s, #args=%d): RPC call being made synchronously", tok, len(resolvedArgsMap))
	resp, err := ctx.monitor.Invoke(ctx.ctx, &pulumirpc.InvokeRequest{
		Tok:               tok,
		Args:              rpcArgs,
		Provider:          providerRef,
		Version:           options.Version,
		PluginDownloadURL: options.PluginDownloadURL,
		AcceptResources:   !disableResourceReferences,
	})
	if err != nil {
		logging.V(9).Infof("Invoke(%s, ...): error: %v", tok, err)
		return err
	}

	// If there were any failures from the provider, return them.
	if len(resp.Failures) > 0 {
		logging.V(9).Infof("Invoke(%s, ...): success: w/ %d failures", tok, len(resp.Failures))
		var ferr error
		for _, failure := range resp.Failures {
			ferr = multierror.Append(ferr,
				fmt.Errorf("%s invoke failed: %s (%s)", tok, failure.Reason, failure.Property))
		}
		return ferr
	}

	// Otherwise, simply unmarshal the output properties and return the result.
	outProps, err := plugin.UnmarshalProperties(
		resp.Return,
		ctx.withKeepOrRejectUnknowns(plugin.MarshalOptions{
			KeepSecrets:   true,
			KeepResources: true,
		}),
	)
	if err != nil {
		return err
	}

	// fail if there are secrets returned from the invoke
	hasSecret, err := unmarshalOutput(ctx, resource.NewObjectProperty(outProps), resultV.Elem())
	if err != nil {
		return err
	}
	if hasSecret {
		return errors.New("unexpected secret result returned to invoke call")
	}
	logging.V(9).Infof("Invoke(%s, ...): success: w/ %d outs (err=%v)", tok, len(outProps), err)
	return nil
}

// Call will invoke a provider call function, identified by its token tok.
//
// output is used to determine the output type to return; self is optional for methods.
func (ctx *Context) Call(tok string, args Input, output Output, self Resource, opts ...InvokeOption) (Output, error) {
	if tok == "" {
		return nil, errors.New("call token must not be empty")
	}

	output = ctx.newOutput(reflect.TypeOf(output))

	options := &invokeOptions{}
	for _, o := range opts {
		if o != nil {
			o.applyInvokeOption(options)
		}
	}

	// Note that we're about to make an outstanding RPC request, so that we can rendezvous during shutdown.
	if err := ctx.beginRPC(); err != nil {
		return nil, err
	}

	prepareCallRequest := func() (*pulumirpc.CallRequest, error) {
		// Determine the provider, version and url to use.
		var provider ProviderResource
		var version string
		var pluginURL string
		if self != nil {
			provider = self.getProvider()
			version = self.getVersion()
			pluginURL = self.getPluginDownloadURL()
		} else {
			provider = mergeProviders(tok, options.Parent, options.Provider, nil)[getPackage(tok)]
			version = options.Version
			pluginURL = options.PluginDownloadURL
		}
		var providerRef string
		if provider != nil {
			pr, err := ctx.resolveProviderReference(provider)
			if err != nil {
				return nil, err
			}
			providerRef = pr
		}

		// Serialize all args, first by awaiting them, and then marshaling them to the requisite gRPC values.
		resolvedArgs, argDeps, _, err := marshalInputs(args)
		if err != nil {
			return nil, fmt.Errorf("marshaling args: %w", err)
		}

		// If we have a value for self, add it to the arguments.
		if self != nil {
			var deps []URN
			resolvedSelf, selfDeps, err := marshalInput(self, reflect.TypeOf(self), true)
			if err != nil {
				return nil, fmt.Errorf("marshaling __self__: %w", err)
			}
			for _, dep := range selfDeps {
				depURN, _, _, err := dep.URN().awaitURN(context.TODO())
				if err != nil {
					return nil, err
				}
				deps = append(deps, depURN)
			}
			resolvedArgs["__self__"] = resolvedSelf
			argDeps["__self__"] = deps
		}

		// Marshal all properties for the RPC call.
		rpcArgs, err := plugin.MarshalProperties(
			resolvedArgs,
			ctx.withKeepOrRejectUnknowns(plugin.MarshalOptions{
				KeepSecrets:      true,
				KeepResources:    ctx.keepResources,
				KeepOutputValues: ctx.keepOutputValues,
			}))
		if err != nil {
			return nil, fmt.Errorf("marshaling args: %w", err)
		}

		// Convert the arg dependencies map for RPC and remove duplicates.
		rpcArgDeps := make(map[string]*pulumirpc.CallRequest_ArgumentDependencies)
		for k, deps := range argDeps {
			sort.Slice(deps, func(i, j int) bool { return deps[i] < deps[j] })

			urns := make([]string, 0, len(deps))
			for i, d := range deps {
				if i > 0 && urns[i-1] == string(d) {
					continue
				}
				urns = append(urns, string(d))
			}

			rpcArgDeps[k] = &pulumirpc.CallRequest_ArgumentDependencies{
				Urns: urns,
			}
		}

		return &pulumirpc.CallRequest{
			Tok:               tok,
			Args:              rpcArgs,
			ArgDependencies:   rpcArgDeps,
			Provider:          providerRef,
			Version:           version,
			PluginDownloadURL: pluginURL,
		}, nil
	}

	// Kick off the call.
	go func() {
		var ret *structpb.Struct
		var deps []Resource
		var err error
		defer func() {
			defer ctx.endRPC(err)

			var outprops resource.PropertyMap
			if err == nil {
				outprops, err = plugin.UnmarshalProperties(ret, ctx.withKeepOrRejectUnknowns(plugin.MarshalOptions{
					KeepSecrets:   true,
					KeepResources: true,
				}))
			}
			if err != nil {
				logging.V(9).Infof("Call(%s, ...): success: w/ unmarshal error: %v", tok, err)
				output.getState().reject(err)
				return
			}

			// Allocate storage for the unmarshalled output.
			var secret bool
			dest := reflect.New(output.ElementType()).Elem()
			known := !outprops.ContainsUnknowns()
			secret, err = unmarshalOutput(ctx, resource.NewObjectProperty(outprops), dest)
			if err != nil {
				output.getState().reject(err)
			} else {
				output.getState().resolve(dest.Interface(), known, secret, deps)
			}

			logging.V(9).Infof("Call(%s, ...): success: w/ %d outs (err=%v)", tok, len(outprops), err)
		}()

		// Prepare the RPC request.
		var req *pulumirpc.CallRequest
		req, err = prepareCallRequest()
		if err != nil {
			return
		}

		// Now, call the RPC.
		var resp *pulumirpc.CallResponse
		logging.V(9).Infof("Call(%s): Goroutine spawned, RPC call being made", tok)
		resp, err = ctx.monitor.Call(ctx.ctx, req)
		if err != nil {
			logging.V(9).Infof("Call(%s, ...): error: %v", tok, err)
		} else if len(resp.Failures) > 0 {
			logging.V(9).Infof("Call(%s, ...): success: w/ %d failures", tok, len(resp.Failures))
			for _, failure := range resp.Failures {
				err = multierror.Append(err,
					fmt.Errorf("%s call failed: %s (%s)", tok, failure.Reason, failure.Property))
			}
		}

		if resp != nil {
			ret = resp.Return

			// Combine the individual dependencies into a single set of dependency resources.
			urns := make(map[string]struct{})
			for _, returnDependencies := range resp.GetReturnDependencies() {
				for _, urn := range returnDependencies.GetUrns() {
					urns[urn] = struct{}{}
				}
			}
			for urn := range urns {
				deps = append(deps, ctx.newDependencyResource(URN(urn)))
			}
		}
	}()

	return output, nil
}

// ReadResource reads an existing custom resource's state from the resource monitor. t is the fully qualified type
// token and name is the "name" part to use in creating a stable and globally unique URN for the object. id is the ID
// of the resource to read, and props contains any state necessary to perform the read (typically props will be nil).
// opts contains optional settings that govern the way the resource is managed.
//
// The value passed to resource must be a pointer to a struct. The fields of this struct that correspond to output
// properties of the resource must have types that are assignable from Output, and must have a `pulumi` tag that
// records the name of the corresponding output property. The struct must embed the CustomResourceState type.
//
// For example, given a custom resource with an int-typed output "foo" and a string-typed output "bar", one would
// define the following CustomResource type:
//
//     type MyResource struct {
//         pulumi.CustomResourceState
//
//         Foo pulumi.IntOutput    `pulumi:"foo"`
//         Bar pulumi.StringOutput `pulumi:"bar"`
//     }
//
// And invoke ReadResource like so:
//
//     var resource MyResource
//     err := ctx.ReadResource(tok, name, id, nil, &resource, opts...)
//
func (ctx *Context) ReadResource(
	t, name string, id IDInput, props Input, resource CustomResource, opts ...ResourceOption) error {
	if t == "" {
		return errors.New("resource type argument cannot be empty")
	} else if name == "" {
		return errors.New("resource name argument (for URN creation) cannot be empty")
	} else if id == nil {
		return errors.New("resource ID is required for lookup and cannot be empty")
	}

	if props != nil {
		propsType := reflect.TypeOf(props)
		if propsType.Kind() == reflect.Ptr {
			propsType = propsType.Elem()
		}
		if !(propsType.Kind() == reflect.Struct ||
			(propsType.Kind() == reflect.Map && propsType.Key().Kind() == reflect.String)) {
			return errors.New("props must be a struct or map or a pointer to a struct or map")
		}
	}

	options := merge(opts...)
	aliasParent := options.Parent
	if options.Parent == nil {
		options.Parent = ctx.stack
	}

	// Before anything else, if there are transformations registered, give them a chance to run to modify the
	// user-provided properties and options assigned to this resource.
	props, options, transformations, err := applyTransformations(t, name, props, resource, opts, options)
	if err != nil {
		return err
	}

	// Collapse aliases to URNs.
	aliasURNs, err := ctx.collapseAliases(options.Aliases, t, name, aliasParent)
	if err != nil {
		return err
	}

	// Note that we're about to make an outstanding RPC request, so that we can rendezvous during shutdown.
	if err := ctx.beginRPC(); err != nil {
		return err
	}

	// Merge providers.
	providers := mergeProviders(t, options.Parent, options.Provider, options.Providers)

	// Get the provider for the resource.
	provider := getProvider(t, options.Provider, providers)

	// Create resolvers for the resource's outputs.
	res := ctx.makeResourceState(t, name, resource, providers, provider,
		options.Version, options.PluginDownloadURL, aliasURNs, transformations)

	// Kick off the resource read operation.  This will happen asynchronously and resolve the above properties.
	go func() {
		// No matter the outcome, make sure all promises are resolved and that we've signaled completion of this RPC.
		var urn, resID string
		var inputs *resourceInputs
		var state *structpb.Struct
		var err error
		defer func() {
			res.resolve(ctx, err, inputs, urn, resID, state, nil)
			ctx.endRPC(err)
		}()

		idToRead, known, _, err := id.ToIDOutput().awaitID(context.TODO())
		if !known || err != nil {
			return
		}

		// Prepare the inputs for an impending operation.
		inputs, err = ctx.prepareResourceInputs(resource, props, t, options, res, false)
		if err != nil {
			return
		}

		logging.V(9).Infof("ReadResource(%s, %s): Goroutine spawned, RPC call being made", t, name)
		resp, err := ctx.monitor.ReadResource(ctx.ctx, &pulumirpc.ReadResourceRequest{
			Type:                    t,
			Name:                    name,
			Parent:                  inputs.parent,
			Properties:              inputs.rpcProps,
			Provider:                inputs.provider,
			Id:                      string(idToRead),
			Aliases:                 inputs.aliases,
			AcceptSecrets:           true,
			AcceptResources:         !disableResourceReferences,
			AdditionalSecretOutputs: inputs.additionalSecretOutputs,
		})
		if err != nil {
			logging.V(9).Infof("ReadResource(%s, %s): error: %v", t, name, err)
		} else {
			logging.V(9).Infof("ReadResource(%s, %s): success: %s %s ...", t, name, resp.Urn, id)
		}
		if resp != nil {
			urn, resID = resp.Urn, string(idToRead)
			state = resp.Properties
		}
	}()

	return nil
}

// RegisterResource creates and registers a new resource object. t is the fully qualified type token and name is
// the "name" part to use in creating a stable and globally unique URN for the object. props contains the goal state
// for the resource object and opts contains optional settings that govern the way the resource is created.
//
// The value passed to resource must be a pointer to a struct. The fields of this struct that correspond to output
// properties of the resource must have types that are assignable from Output, and must have a `pulumi` tag that
// records the name of the corresponding output property. The struct must embed either the ResourceState or the
// CustomResourceState type.
//
// For example, given a custom resource with an int-typed output "foo" and a string-typed output "bar", one would
// define the following CustomResource type:
//
//     type MyResource struct {
//         pulumi.CustomResourceState
//
//         Foo pulumi.IntOutput    `pulumi:"foo"`
//         Bar pulumi.StringOutput `pulumi:"bar"`
//     }
//
// And invoke RegisterResource like so:
//
//     var resource MyResource
//     err := ctx.RegisterResource(tok, name, props, &resource, opts...)
//
func (ctx *Context) RegisterResource(
	t, name string, props Input, resource Resource, opts ...ResourceOption) error {

	return ctx.registerResource(t, name, props, resource, false /*remote*/, opts...)
}

func (ctx *Context) getResource(urn string) (*pulumirpc.RegisterResourceResponse, error) {
	// This is a resource that already exists. Read its state from the engine.
	resolvedArgsMap := resource.NewPropertyMapFromMap(map[string]interface{}{
		"urn": urn,
	})

	rpcArgs, err := plugin.MarshalProperties(
		resolvedArgsMap,
		ctx.withKeepOrRejectUnknowns(plugin.MarshalOptions{
			KeepSecrets:   true,
			KeepResources: ctx.keepResources,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("marshaling arguments: %w", err)
	}

	tok := "pulumi:pulumi:getResource"
	logging.V(9).Infof("Invoke(%s, #args=%d): RPC call being made synchronously", tok, len(resolvedArgsMap))
	resp, err := ctx.monitor.Invoke(ctx.ctx, &pulumirpc.InvokeRequest{
		Tok:  "pulumi:pulumi:getResource",
		Args: rpcArgs,
	})
	if err != nil {
		return nil, fmt.Errorf("Invoke(%s, ...): error: %v", tok, err)
	}

	// If there were any failures from the provider, return them.
	if len(resp.Failures) > 0 {
		logging.V(9).Infof("Invoke(%s, ...): success: w/ %d failures", tok, len(resp.Failures))
		var ferr error
		for _, failure := range resp.Failures {
			ferr = multierror.Append(ferr,
				fmt.Errorf("%s invoke failed: %s (%s)", tok, failure.Reason, failure.Property))
		}
		return nil, ferr
	}

	return &pulumirpc.RegisterResourceResponse{
		Urn:    resp.Return.Fields["urn"].GetStringValue(),
		Id:     resp.Return.Fields["id"].GetStringValue(),
		Object: resp.Return.Fields["state"].GetStructValue(),
	}, nil
}

func (ctx *Context) registerResource(
	t, name string, props Input, resource Resource, remote bool, opts ...ResourceOption) error {
	if t == "" {
		return errors.New("resource type argument cannot be empty")
	} else if name == "" {
		return errors.New("resource name argument (for URN creation) cannot be empty")
	}

	_, custom := resource.(CustomResource)
	if !custom && remote {
		resource.markRemoteComponent()
	}

	if _, isProvider := resource.(ProviderResource); isProvider && !strings.HasPrefix(t, "pulumi:providers:") {
		return errors.New("provider resource type must begin with \"pulumi:providers:\"")
	}

	if props != nil {
		propsType := reflect.TypeOf(props)
		if propsType.Kind() == reflect.Ptr {
			propsType = propsType.Elem()
		}
		if !(propsType.Kind() == reflect.Struct ||
			(propsType.Kind() == reflect.Map && propsType.Key().Kind() == reflect.String)) {
			return errors.New("props must be a struct or map or a pointer to a struct or map")
		}
	}

	options := merge(opts...)
	parent := options.Parent
	if options.Parent == nil {
		options.Parent = ctx.stack
	}

	// Before anything else, if there are transformations registered, give them a chance to run to modify the
	// user-provided properties and options assigned to this resource.
	props, options, transformations, err := applyTransformations(t, name, props, resource, opts, options)
	if err != nil {
		return err
	}

	// Collapse aliases to URNs.
	aliasURNs, err := ctx.collapseAliases(options.Aliases, t, name, parent)
	if err != nil {
		return err
	}

	// Note that we're about to make an outstanding RPC request, so that we can rendezvous during shutdown.
	if err := ctx.beginRPC(); err != nil {
		return err
	}

	// Merge providers.
	providers := mergeProviders(t, options.Parent, options.Provider, options.Providers)

	// Get the provider for the resource.
	provider := getProvider(t, options.Provider, providers)

	// Create resolvers for the resource's outputs.
	resState := ctx.makeResourceState(t, name, resource, providers, provider,
		options.Version, options.PluginDownloadURL, aliasURNs, transformations)

	// Kick off the resource registration.  If we are actually performing a deployment, the resulting properties
	// will be resolved asynchronously as the RPC operation completes.  If we're just planning, values won't resolve.
	go func() {
		// No matter the outcome, make sure all promises are resolved and that we've signaled completion of this RPC.
		var urn, resID string
		var inputs *resourceInputs
		var state *structpb.Struct
		deps := make(map[string][]Resource)
		var err error
		defer func() {
			resState.resolve(ctx, err, inputs, urn, resID, state, deps)
			ctx.endRPC(err)
		}()

		// Prepare the inputs for an impending operation.
		inputs, err = ctx.prepareResourceInputs(resource, props, t, options, resState, remote)
		if err != nil {
			return
		}

		var resp *pulumirpc.RegisterResourceResponse
		if len(options.URN) > 0 {
			resp, err = ctx.getResource(options.URN)
			if err != nil {
				logging.V(9).Infof("getResource(%s, %s): error: %v", t, name, err)
			} else {
				logging.V(9).Infof("getResource(%s, %s): success: %s %s ...", t, name, resp.Urn, resp.Id)
			}
		} else {
			logging.V(9).Infof("RegisterResource(%s, %s): Goroutine spawned, RPC call being made", t, name)
			resp, err = ctx.monitor.RegisterResource(ctx.ctx, &pulumirpc.RegisterResourceRequest{
				Type:                    t,
				Name:                    name,
				Parent:                  inputs.parent,
				Object:                  inputs.rpcProps,
				Custom:                  custom,
				Protect:                 inputs.protect,
				Dependencies:            inputs.deps,
				Provider:                inputs.provider,
				Providers:               inputs.providers,
				PropertyDependencies:    inputs.rpcPropertyDeps,
				DeleteBeforeReplace:     inputs.deleteBeforeReplace,
				ImportId:                inputs.importID,
				CustomTimeouts:          inputs.customTimeouts,
				IgnoreChanges:           inputs.ignoreChanges,
				Aliases:                 inputs.aliases,
				AcceptSecrets:           true,
				AcceptResources:         !disableResourceReferences,
				AdditionalSecretOutputs: inputs.additionalSecretOutputs,
				Version:                 inputs.version,
				PluginDownloadURL:       inputs.pluginDownloadURL,
				Remote:                  remote,
				ReplaceOnChanges:        inputs.replaceOnChanges,
			})
			if err != nil {
				logging.V(9).Infof("RegisterResource(%s, %s): error: %v", t, name, err)
			} else {
				logging.V(9).Infof("RegisterResource(%s, %s): success: %s %s ...", t, name, resp.Urn, resp.Id)
			}
		}

		if resp != nil {
			urn, resID = resp.Urn, resp.Id
			state = resp.Object
			for key, propertyDependencies := range resp.GetPropertyDependencies() {
				var resources []Resource
				for _, urn := range propertyDependencies.GetUrns() {
					resources = append(resources, &ResourceState{urn: URNInput(URN(urn)).ToURNOutput()})
				}
				deps[key] = resources
			}
		}
	}()

	return nil
}

func (ctx *Context) RegisterComponentResource(
	t, name string, resource ComponentResource, opts ...ResourceOption) error {

	return ctx.RegisterResource(t, name, nil /*props*/, resource, opts...)
}

func (ctx *Context) RegisterRemoteComponentResource(
	t, name string, props Input, resource ComponentResource, opts ...ResourceOption) error {

	return ctx.registerResource(t, name, props, resource, true /*remote*/, opts...)
}

// resourceState contains the results of a resource registration operation.
type resourceState struct {
	outputs           map[string]Output
	providers         map[string]ProviderResource
	provider          ProviderResource
	version           string
	pluginDownloadURL string
	aliases           []URNOutput
	name              string
	transformations   []ResourceTransformation
}

// Apply transformations and return the transformations themselves, as well as the transformed props and opts.
func applyTransformations(t, name string, props Input, resource Resource, opts []ResourceOption,
	options *resourceOptions) (Input, *resourceOptions, []ResourceTransformation, error) {

	transformations := options.Transformations
	if options.Parent != nil {
		transformations = append(transformations, options.Parent.getTransformations()...)
	}

	for _, transformation := range transformations {
		args := &ResourceTransformationArgs{
			Resource: resource,
			Type:     t,
			Name:     name,
			Props:    props,
			Opts:     opts,
		}

		res := transformation(args)
		if res != nil {
			resOptions := merge(res.Opts...)

			if resOptions.Parent != nil && resOptions.Parent.URN() != options.Parent.URN() {
				return nil, nil, nil, errors.New("transformations cannot currently be used to change the `parent` of a resource")
			}
			props = res.Props
			options = resOptions
		}
	}

	return props, options, transformations, nil
}

// checks all possible sources of providers and merges them with preference given to the most specific
func mergeProviders(t string, parent Resource, provider ProviderResource,
	providers map[string]ProviderResource) map[string]ProviderResource {

	// copy parent providers
	result := make(map[string]ProviderResource)
	if parent != nil {
		for k, v := range parent.getProviders() {
			result[k] = v
		}
	}

	// copy provider map
	for k, v := range providers {
		result[k] = v
	}

	// copy specific provider, if any
	if provider != nil {
		pkg := getPackage(t)
		result[pkg] = provider
	}

	return result
}

// getProvider gets the provider for the resource.
func getProvider(t string, provider ProviderResource, providers map[string]ProviderResource) ProviderResource {
	if provider == nil {
		pkg := getPackage(t)
		provider = providers[pkg]
	}
	return provider
}

// getPackage takes in a type and returns the pkg
func getPackage(t string) string {
	components := strings.Split(t, ":")
	if len(components) != 3 {
		return ""
	}
	return components[0]
}

// collapseAliases collapses a list of Aliases into a list of URNs. Parent aliases
// are also included. If there are N child aliases, and M parent aliases, there will
// be (M+1)*(N+1)-1 total aliases, or, as calculated in the logic below, N+(M*(1+N)).
func (ctx *Context) collapseAliases(aliases []Alias, t, name string, parent Resource) ([]URNOutput, error) {
	project, stack := ctx.Project(), ctx.Stack()

	var aliasURNs []URNOutput

	for _, alias := range aliases {
		urn, err := alias.collapseToURN(name, t, parent, project, stack)
		if err != nil {
			return nil, fmt.Errorf("error collapsing alias to URN: %w", err)
		}
		aliasURNs = append(aliasURNs, urn)
	}

	if parent != nil {
		parentAliases := parent.getAliases()
		for i := range parentAliases {
			parentAlias := parentAliases[i]
			urn := inheritedChildAlias(name, parent.getName(), t, project, stack, parentAlias)
			aliasURNs = append(aliasURNs, urn)
			for j := range aliases {
				childAlias := aliases[j]
				urn, err := childAlias.collapseToURN(name, t, parent, project, stack)
				if err != nil {
					return nil, fmt.Errorf("error collapsing alias to URN: %w", err)
				}
				inheritedAlias := urn.ApplyT(func(urn URN) URNOutput {
					aliasedChildName := string(resource.URN(urn).Name())
					aliasedChildType := string(resource.URN(urn).Type())
					return inheritedChildAlias(aliasedChildName, parent.getName(), aliasedChildType, project, stack, parentAlias)
				}).ApplyT(func(urn interface{}) URN {
					return urn.(URN)
				}).(URNOutput)
				aliasURNs = append(aliasURNs, inheritedAlias)
			}
		}
	}

	return aliasURNs, nil
}

var mapOutputType = reflect.TypeOf((*MapOutput)(nil)).Elem()

// makeResourceState creates a set of resolvers that we'll use to finalize state, for URNs, IDs, and output
// properties.
func (ctx *Context) makeResourceState(t, name string, resourceV Resource, providers map[string]ProviderResource,
	provider ProviderResource, version, pluginDownloadURL string, aliases []URNOutput,
	transformations []ResourceTransformation) *resourceState {

	// Ensure that the input resource is a pointer to a struct. Note that we don't fail if it is not, and we probably
	// ought to.
	resource := reflect.ValueOf(resourceV)
	typ := resource.Type()
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return &resourceState{}
	}
	resource, typ = resource.Elem(), typ.Elem()

	var rs *ResourceState
	var crs *CustomResourceState
	var prs *ProviderResourceState

	// Check to see if a value of exactly `*ResourceState`, `*CustomResourceState`, or `*ProviderResourceState` was
	// provided.
	switch r := resourceV.(type) {
	case *ResourceState:
		rs = r
	case *CustomResourceState:
		crs = r
	case *ProviderResourceState:
		prs = r
	}

	// Find the particular Resource implementation and the settable, `pulumi`-tagged fields in the input type. The
	// former is used for any URN or ID fields; the latter are used to determine the expected outputs of the resource
	// after its RegisterResource call completes. For each of those fields, create an appropriately-typed Output and
	// map the Output to its property name so we can resolve it later.
	state := &resourceState{outputs: map[string]Output{}}
	for i := 0; i < typ.NumField(); i++ {
		fieldV := resource.Field(i)
		if !fieldV.CanSet() {
			continue
		}

		field := typ.Field(i)
		switch {
		case field.Anonymous && field.Type == resourceStateType:
			rs = fieldV.Addr().Interface().(*ResourceState)
		case field.Anonymous && field.Type == customResourceStateType:
			crs = fieldV.Addr().Interface().(*CustomResourceState)
		case field.Anonymous && field.Type == providerResourceStateType:
			prs = fieldV.Addr().Interface().(*ProviderResourceState)
		case field.Type.Implements(outputType):
			tag, has := typ.Field(i).Tag.Lookup("pulumi")
			if !has {
				continue
			}

			output := ctx.newOutput(field.Type, resourceV)
			fieldV.Set(reflect.ValueOf(output))

			if tag == "" && field.Type != mapOutputType {
				output.getState().reject(fmt.Errorf("the field %v must be a MapOutput or its tag must be non-empty", field.Name))
			}

			state.outputs[tag] = output
		}
	}

	// Create provider- and custom resource-specific state/resolvers.
	if prs != nil {
		crs = &prs.CustomResourceState
		prs.pkg = t[len("pulumi:providers:"):]
	}
	if crs != nil {
		rs = &crs.ResourceState
		crs.id = IDOutput{ctx.newOutputState(idType, resourceV)}
		state.outputs["id"] = crs.id
	}

	// Populate ResourceState resolvers. (Pulled into function to keep the nil-ness linter check happy).
	populateResourceStateResolvers := func() {
		contract.Assert(rs != nil)
		state.providers = providers
		rs.providers = providers
		state.provider = provider
		rs.provider = provider
		state.version = version
		rs.version = version
		state.pluginDownloadURL = pluginDownloadURL
		rs.pluginDownloadURL = pluginDownloadURL
		rs.urn = URNOutput{ctx.newOutputState(urnType, resourceV)}
		state.outputs["urn"] = rs.urn
		state.name = name
		rs.name = name
		state.aliases = aliases
		rs.aliases = aliases
		state.transformations = transformations
		rs.transformations = transformations
	}
	populateResourceStateResolvers()

	return state
}

// resolve resolves the resource outputs using the given error and/or values.
func (state *resourceState) resolve(ctx *Context, err error, inputs *resourceInputs, urn, id string,
	result *structpb.Struct, deps map[string][]Resource) {

	dryrun := ctx.DryRun()

	var inprops resource.PropertyMap
	if inputs != nil {
		inprops = inputs.resolvedProps
	}

	var outprops resource.PropertyMap
	if err == nil {
		outprops, err = plugin.UnmarshalProperties(
			result,
			ctx.withKeepOrRejectUnknowns(plugin.MarshalOptions{
				KeepSecrets:   true,
				KeepResources: true,
			}),
		)
	}
	if err != nil {
		// If there was an error, we must reject everything.
		for _, output := range state.outputs {
			output.getState().reject(err)
		}
		return
	}

	outprops["urn"] = resource.NewStringProperty(urn)
	if id != "" || !dryrun {
		outprops["id"] = resource.NewStringProperty(id)
	} else {
		outprops["id"] = resource.MakeComputed(resource.PropertyValue{})
	}

	if _, hasRemainingOutput := state.outputs[""]; hasRemainingOutput {
		remaining, known := resource.PropertyMap{}, true
		for k, v := range outprops {
			if v.IsNull() || v.IsComputed() || v.IsOutput() {
				known = !dryrun
			}
			if _, ok := state.outputs[string(k)]; !ok {
				remaining[k] = v
			}
		}
		if !known {
			outprops[""] = resource.MakeComputed(resource.NewStringProperty(""))
		} else {
			outprops[""] = resource.NewObjectProperty(remaining)
		}
	}

	for k, output := range state.outputs {
		// If this is an unknown or missing value during a dry run, do nothing.
		v, ok := outprops[resource.PropertyKey(k)]
		if !ok && !dryrun {
			v = inprops[resource.PropertyKey(k)]
		}

		known := true
		if v.IsNull() || v.IsComputed() || v.IsOutput() {
			known = !dryrun
		}

		// Allocate storage for the unmarshalled output.
		dest := reflect.New(output.ElementType()).Elem()
		secret, err := unmarshalOutput(ctx, v, dest)
		if err != nil {
			output.getState().reject(err)
		} else {
			output.getState().resolve(dest.Interface(), known, secret, deps[k])
		}
	}
}

// resourceInputs reflects all of the inputs necessary to perform core resource RPC operations.
type resourceInputs struct {
	parent                  string
	deps                    []string
	protect                 bool
	provider                string
	providers               map[string]string
	resolvedProps           resource.PropertyMap
	rpcProps                *structpb.Struct
	rpcPropertyDeps         map[string]*pulumirpc.RegisterResourceRequest_PropertyDependencies
	deleteBeforeReplace     bool
	importID                string
	customTimeouts          *pulumirpc.RegisterResourceRequest_CustomTimeouts
	ignoreChanges           []string
	aliases                 []string
	additionalSecretOutputs []string
	version                 string
	pluginDownloadURL       string
	replaceOnChanges        []string
}

// prepareResourceInputs prepares the inputs for a resource operation, shared between read and register.
func (ctx *Context) prepareResourceInputs(res Resource, props Input, t string, opts *resourceOptions,
	state *resourceState, remote bool) (*resourceInputs, error) {

	// Get the parent and dependency URNs from the options, in addition to the protection bit.  If there wasn't an
	// explicit parent, and a root stack resource exists, we will automatically parent to that.
	resOpts, err := ctx.getOpts(res, t, state.provider, opts, remote)
	if err != nil {
		return nil, fmt.Errorf("resolving options: %w", err)
	}

	// Serialize all properties, first by awaiting them, and then marshaling them to the requisite gRPC values.
	resolvedProps, propertyDeps, rpcDeps, err := marshalInputs(props)
	if err != nil {
		return nil, fmt.Errorf("marshaling properties: %w", err)
	}

	// Marshal all properties for the RPC call.
	rpcProps, err := plugin.MarshalProperties(
		resolvedProps,
		ctx.withKeepOrRejectUnknowns(plugin.MarshalOptions{
			KeepSecrets:   true,
			KeepResources: ctx.keepResources,
			// To initially scope the use of this new feature, we only keep output values when
			// remote is true (for multi-lang components).
			KeepOutputValues: remote && ctx.keepOutputValues,
		}))
	if err != nil {
		return nil, fmt.Errorf("marshaling properties: %w", err)
	}

	// Convert the property dependencies map for RPC and remove duplicates.
	rpcPropertyDeps := make(map[string]*pulumirpc.RegisterResourceRequest_PropertyDependencies)
	for k, deps := range propertyDeps {
		urns := make([]string, len(deps))
		for i, d := range deps {
			urns[i] = string(d)
		}
		sort.Strings(urns)

		rpcPropertyDeps[k] = &pulumirpc.RegisterResourceRequest_PropertyDependencies{
			Urns: urns,
		}
	}

	// Merge all dependencies with what we got earlier from property marshaling, and remove duplicates.
	var deps []string
	depSet := urnSet{}
	for _, dep := range append(resOpts.depURNs, rpcDeps...) {
		if !depSet.has(dep) {
			deps = append(deps, string(dep))
			depSet.add(dep)
		}
	}
	sort.Strings(deps)

	// Await alias URNs
	aliases := make([]string, len(state.aliases))
	for i, alias := range state.aliases {
		urn, _, _, err := alias.awaitURN(context.Background())
		if err != nil {
			return nil, fmt.Errorf("error waiting for alias URN to resolve: %w", err)
		}
		aliases[i] = string(urn)
	}

	return &resourceInputs{
		parent:                  string(resOpts.parentURN),
		deps:                    deps,
		protect:                 resOpts.protect,
		provider:                resOpts.providerRef,
		providers:               resOpts.providerRefs,
		resolvedProps:           resolvedProps,
		rpcProps:                rpcProps,
		rpcPropertyDeps:         rpcPropertyDeps,
		deleteBeforeReplace:     resOpts.deleteBeforeReplace,
		importID:                string(resOpts.importID),
		customTimeouts:          getTimeouts(opts.CustomTimeouts),
		ignoreChanges:           resOpts.ignoreChanges,
		aliases:                 aliases,
		additionalSecretOutputs: resOpts.additionalSecretOutputs,
		version:                 state.version,
		pluginDownloadURL:       state.pluginDownloadURL,
		replaceOnChanges:        resOpts.replaceOnChanges,
	}, nil
}

func getTimeouts(custom *CustomTimeouts) *pulumirpc.RegisterResourceRequest_CustomTimeouts {
	var timeouts pulumirpc.RegisterResourceRequest_CustomTimeouts
	if custom != nil {
		timeouts.Update = custom.Update
		timeouts.Create = custom.Create
		timeouts.Delete = custom.Delete
	}
	return &timeouts
}

// Helper struct for the return type of `getOpts`.
type resourceOpts struct {
	parentURN               URN
	depURNs                 []URN
	protect                 bool
	providerRef             string
	providerRefs            map[string]string
	deleteBeforeReplace     bool
	importID                ID
	ignoreChanges           []string
	additionalSecretOutputs []string
	replaceOnChanges        []string
}

// getOpts returns a set of resource options from an array of them. This includes the parent URN, any dependency URNs,
// a boolean indicating whether the resource is to be protected, and the URN and ID of the resource's provider, if any.
func (ctx *Context) getOpts(res Resource, t string, provider ProviderResource, opts *resourceOptions, remote bool,
) (resourceOpts, error) {

	var importID ID
	if opts.Import != nil {
		id, _, _, err := opts.Import.ToIDOutput().awaitID(context.TODO())
		if err != nil {
			return resourceOpts{}, err
		}
		importID = id
	}

	var parentURN URN
	if opts.Parent != nil {
		opts.Parent.addChild(res)

		urn, _, _, err := opts.Parent.URN().awaitURN(context.TODO())
		if err != nil {
			return resourceOpts{}, err
		}
		parentURN = urn
	}

	var depURNs []URN
	if opts.DependsOn != nil {
		depSet := urnSet{}
		for _, r := range opts.DependsOn {
			dependsOn, err := r(ctx.ctx)
			if err != nil {
				return resourceOpts{}, err
			}
			depSet.union(dependsOn)
		}
		depURNs = depSet.values()
	}

	var providerRef string
	if provider != nil {
		pr, err := ctx.resolveProviderReference(provider)
		if err != nil {
			return resourceOpts{}, err
		}
		providerRef = pr
	}

	var providerRefs map[string]string
	if remote {
		if opts.Providers != nil {
			providerRefs = make(map[string]string, len(opts.Providers))
			for name, provider := range opts.Providers {
				pr, err := ctx.resolveProviderReference(provider)
				if err != nil {
					return resourceOpts{}, err
				}
				providerRefs[name] = pr
			}
		}
	}

	return resourceOpts{
		parentURN:               parentURN,
		depURNs:                 depURNs,
		protect:                 opts.Protect,
		providerRef:             providerRef,
		providerRefs:            providerRefs,
		deleteBeforeReplace:     opts.DeleteBeforeReplace,
		importID:                importID,
		ignoreChanges:           opts.IgnoreChanges,
		additionalSecretOutputs: opts.AdditionalSecretOutputs,
		replaceOnChanges:        opts.ReplaceOnChanges,
	}, nil
}

func (ctx *Context) resolveProviderReference(provider ProviderResource) (string, error) {
	urn, _, _, err := provider.URN().awaitURN(context.TODO())
	if err != nil {
		return "", err
	}
	id, known, _, err := provider.ID().awaitID(context.TODO())
	if err != nil {
		return "", err
	}
	if !known {
		id = rpcTokenUnknownValue
	}
	return string(urn) + "::" + string(id), nil
}

// noMoreRPCs is a sentinel value used to stop subsequent RPCs from occurring.
const noMoreRPCs = -1

// beginRPC attempts to start a new RPC request, returning a non-nil error if no more RPCs are permitted
// (usually because the program is shutting down).
func (ctx *Context) beginRPC() error {
	ctx.rpcsLock.Lock()
	defer ctx.rpcsLock.Unlock()

	// If we're done with RPCs, return an error.
	if ctx.rpcs == noMoreRPCs {
		return errors.New("attempted illegal RPC after program completion")
	}

	ctx.rpcs++
	return nil
}

// endRPC signals the completion of an RPC and notifies any potential awaiters when outstanding RPCs hit zero.
func (ctx *Context) endRPC(err error) {
	ctx.rpcsLock.Lock()
	defer ctx.rpcsLock.Unlock()

	if err != nil && ctx.rpcError == nil {
		ctx.rpcError = err
	}

	ctx.rpcs--
	if ctx.rpcs == 0 {
		ctx.rpcsDone.Broadcast()
	}
}

// RegisterResourceOutputs completes the resource registration, attaching an optional set of computed outputs.
func (ctx *Context) RegisterResourceOutputs(resource Resource, outs Map) error {
	// Note that we're about to make an outstanding RPC request, so that we can rendezvous during shutdown.
	if err := ctx.beginRPC(); err != nil {
		return err
	}

	go func() {
		// No matter the outcome, make sure all promises are resolved and that we've signaled completion of this RPC.
		var err error
		defer func() {
			// Signal the completion of this RPC and notify any potential awaiters.
			ctx.endRPC(err)
		}()

		urn, _, _, err := resource.URN().awaitURN(context.TODO())
		if err != nil {
			return
		}

		outsResolved, _, err := marshalInput(outs, anyType, true)
		if err != nil {
			return
		}

		outsMarshalled, err := plugin.MarshalProperties(
			outsResolved.ObjectValue(),
			ctx.withKeepOrRejectUnknowns(plugin.MarshalOptions{
				KeepSecrets:   true,
				KeepResources: ctx.keepResources,
			}))
		if err != nil {
			return
		}

		// Register the outputs
		logging.V(9).Infof("RegisterResourceOutputs(%s): RPC call being made", urn)
		_, err = ctx.monitor.RegisterResourceOutputs(ctx.ctx, &pulumirpc.RegisterResourceOutputsRequest{
			Urn:     string(urn),
			Outputs: outsMarshalled,
		})

		logging.V(9).Infof("RegisterResourceOutputs(%s): %v", urn, err)
	}()

	return nil
}

// Export registers a key and value pair with the current context's stack.
func (ctx *Context) Export(name string, value Input) {
	ctx.exports[name] = value
}

// RegisterStackTransformation adds a transformation to all future resources constructed in this Pulumi stack.
func (ctx *Context) RegisterStackTransformation(t ResourceTransformation) error {
	ctx.stack.addTransformation(t)
	return nil
}

func (ctx *Context) newOutputState(elementType reflect.Type, deps ...Resource) *OutputState {
	return newOutputState(&ctx.join, elementType, deps...)
}

func (ctx *Context) newOutput(typ reflect.Type, deps ...Resource) Output {
	return newOutput(&ctx.join, typ, deps...)
}

// NewOutput creates a new output associated with this context.
func (ctx *Context) NewOutput() (Output, func(interface{}), func(error)) {
	return newAnyOutput(&ctx.join)
}

// Sets marshalling flags based on `ctx.DryRun()`: we will either
// preserve unknowns as-is or fail strictly with an exception if any
// unkowns are found. The third option, filtering out unknown values
// from the data structure being marshalled, is never used.
func (ctx *Context) withKeepOrRejectUnknowns(options plugin.MarshalOptions) plugin.MarshalOptions {
	if ctx.DryRun() {
		options.KeepUnknowns = true
	} else {
		options.RejectUnknowns = true
	}
	return options
}
