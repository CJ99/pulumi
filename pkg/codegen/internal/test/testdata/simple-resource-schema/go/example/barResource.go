// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package example

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type BarResource struct {
	pulumi.ResourceState

	Foo ResourceOutput `pulumi:"foo"`
}

// NewBarResource registers a new resource with the given unique name, arguments, and options.
func NewBarResource(ctx *pulumi.Context,
	name string, args *BarResourceArgs, opts ...pulumi.ResourceOption) (*BarResource, error) {
	if args == nil {
		args = &BarResourceArgs{}
	}

	var resource BarResource
	err := ctx.RegisterRemoteComponentResource("bar::BarResource", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

type barResourceArgs struct {
	Foo *Resource `pulumi:"foo"`
}

// The set of arguments for constructing a BarResource resource.
type BarResourceArgs struct {
	Foo ResourceInput
}

func (BarResourceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*barResourceArgs)(nil)).Elem()
}

type BarResourceInput interface {
	pulumi.Input

	ToBarResourceOutput() BarResourceOutput
	ToBarResourceOutputWithContext(ctx context.Context) BarResourceOutput
}

func (*BarResource) ElementType() reflect.Type {
	return reflect.TypeOf((**BarResource)(nil)).Elem()
}

func (i *BarResource) ToBarResourceOutput() BarResourceOutput {
	return i.ToBarResourceOutputWithContext(context.Background())
}

func (i *BarResource) ToBarResourceOutputWithContext(ctx context.Context) BarResourceOutput {
	return pulumi.ToOutputWithContext(ctx, i).(BarResourceOutput)
}

type BarResourceOutput struct{ *pulumi.OutputState }

func (BarResourceOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**BarResource)(nil)).Elem()
}

func (o BarResourceOutput) ToBarResourceOutput() BarResourceOutput {
	return o
}

func (o BarResourceOutput) ToBarResourceOutputWithContext(ctx context.Context) BarResourceOutput {
	return o
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*BarResourceInput)(nil)).Elem(), &BarResource{})
	pulumi.RegisterOutputType(BarResourceOutput{})
}
