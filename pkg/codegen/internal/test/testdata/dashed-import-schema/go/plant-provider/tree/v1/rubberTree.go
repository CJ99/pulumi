// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package v1

import (
	"context"
	"reflect"

	"dashed-import-schema/plant-provider"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type RubberTree struct {
	pulumi.CustomResourceState

	Container plantprovider.ContainerPtrOutput `pulumi:"container"`
	Diameter  DiameterOutput                   `pulumi:"diameter"`
	Farm      pulumi.StringPtrOutput           `pulumi:"farm"`
	Size      TreeSizePtrOutput                `pulumi:"size"`
	Type      RubberTreeVarietyOutput          `pulumi:"type"`
}

// NewRubberTree registers a new resource with the given unique name, arguments, and options.
func NewRubberTree(ctx *pulumi.Context,
	name string, args *RubberTreeArgs, opts ...pulumi.ResourceOption) (*RubberTree, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Container != nil {
		args.Container = args.Container.ToContainerPtrOutput().ApplyT(func(v *plantprovider.Container) *plantprovider.Container { return v.Defaults() }).(plantprovider.ContainerPtrOutput)
	}
	if isZero(args.Diameter) {
		args.Diameter = Diameter(6.0)
	}
	if isZero(args.Farm) {
		args.Farm = pulumi.StringPtr("(unknown)")
	}
	if isZero(args.Size) {
		args.Size = TreeSize("medium")
	}
	if isZero(args.Type) {
		args.Type = RubberTreeVariety("Burgundy")
	}
	var resource RubberTree
	err := ctx.RegisterResource("plant:tree/v1:RubberTree", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetRubberTree gets an existing RubberTree resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetRubberTree(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *RubberTreeState, opts ...pulumi.ResourceOption) (*RubberTree, error) {
	var resource RubberTree
	err := ctx.ReadResource("plant:tree/v1:RubberTree", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering RubberTree resources.
type rubberTreeState struct {
	Farm *string `pulumi:"farm"`
}

type RubberTreeState struct {
	Farm pulumi.StringPtrInput
}

func (RubberTreeState) ElementType() reflect.Type {
	return reflect.TypeOf((*rubberTreeState)(nil)).Elem()
}

type rubberTreeArgs struct {
	Container *plantprovider.Container `pulumi:"container"`
	Diameter  Diameter                 `pulumi:"diameter"`
	Farm      *string                  `pulumi:"farm"`
	Size      *TreeSize                `pulumi:"size"`
	Type      RubberTreeVariety        `pulumi:"type"`
}

// The set of arguments for constructing a RubberTree resource.
type RubberTreeArgs struct {
	Container plantprovider.ContainerPtrInput
	Diameter  DiameterInput
	Farm      pulumi.StringPtrInput
	Size      TreeSizePtrInput
	Type      RubberTreeVarietyInput
}

func (RubberTreeArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*rubberTreeArgs)(nil)).Elem()
}

type RubberTreeInput interface {
	pulumi.Input

	ToRubberTreeOutput() RubberTreeOutput
	ToRubberTreeOutputWithContext(ctx context.Context) RubberTreeOutput
}

func (*RubberTree) ElementType() reflect.Type {
	return reflect.TypeOf((**RubberTree)(nil)).Elem()
}

func (i *RubberTree) ToRubberTreeOutput() RubberTreeOutput {
	return i.ToRubberTreeOutputWithContext(context.Background())
}

func (i *RubberTree) ToRubberTreeOutputWithContext(ctx context.Context) RubberTreeOutput {
	return pulumi.ToOutputWithContext(ctx, i).(RubberTreeOutput)
}

type RubberTreeOutput struct{ *pulumi.OutputState }

func (RubberTreeOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**RubberTree)(nil)).Elem()
}

func (o RubberTreeOutput) ToRubberTreeOutput() RubberTreeOutput {
	return o
}

func (o RubberTreeOutput) ToRubberTreeOutputWithContext(ctx context.Context) RubberTreeOutput {
	return o
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*RubberTreeInput)(nil)).Elem(), &RubberTree{})
	pulumi.RegisterOutputType(RubberTreeOutput{})
}
