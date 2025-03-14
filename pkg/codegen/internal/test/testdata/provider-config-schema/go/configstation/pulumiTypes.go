// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package configstation

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Child struct {
	Age  *int    `pulumi:"age"`
	Name *string `pulumi:"name"`
}

// ChildInput is an input type that accepts ChildArgs and ChildOutput values.
// You can construct a concrete instance of `ChildInput` via:
//
//          ChildArgs{...}
type ChildInput interface {
	pulumi.Input

	ToChildOutput() ChildOutput
	ToChildOutputWithContext(context.Context) ChildOutput
}

type ChildArgs struct {
	Age  pulumi.IntPtrInput    `pulumi:"age"`
	Name pulumi.StringPtrInput `pulumi:"name"`
}

func (ChildArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*Child)(nil)).Elem()
}

func (i ChildArgs) ToChildOutput() ChildOutput {
	return i.ToChildOutputWithContext(context.Background())
}

func (i ChildArgs) ToChildOutputWithContext(ctx context.Context) ChildOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ChildOutput)
}

type ChildOutput struct{ *pulumi.OutputState }

func (ChildOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Child)(nil)).Elem()
}

func (o ChildOutput) ToChildOutput() ChildOutput {
	return o
}

func (o ChildOutput) ToChildOutputWithContext(ctx context.Context) ChildOutput {
	return o
}

func (o ChildOutput) Age() pulumi.IntPtrOutput {
	return o.ApplyT(func(v Child) *int { return v.Age }).(pulumi.IntPtrOutput)
}

func (o ChildOutput) Name() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Child) *string { return v.Name }).(pulumi.StringPtrOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ChildInput)(nil)).Elem(), ChildArgs{})
	pulumi.RegisterOutputType(ChildOutput{})
}
