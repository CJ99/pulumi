// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package example

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"plain-object-disable-defaults/example/mod1"
	"plain-object-disable-defaults/example/mod2"
)

// BETA FEATURE - Options to configure the Helm Release resource.
type HelmReleaseSettings struct {
	// The backend storage driver for Helm. Values are: configmap, secret, memory, sql.
	Driver *string `pulumi:"driver"`
	// The path to the helm plugins directory.
	PluginsPath *string `pulumi:"pluginsPath"`
	// to test required args
	RequiredArg string `pulumi:"requiredArg"`
}

// HelmReleaseSettingsInput is an input type that accepts HelmReleaseSettingsArgs and HelmReleaseSettingsOutput values.
// You can construct a concrete instance of `HelmReleaseSettingsInput` via:
//
//          HelmReleaseSettingsArgs{...}
type HelmReleaseSettingsInput interface {
	pulumi.Input

	ToHelmReleaseSettingsOutput() HelmReleaseSettingsOutput
	ToHelmReleaseSettingsOutputWithContext(context.Context) HelmReleaseSettingsOutput
}

// BETA FEATURE - Options to configure the Helm Release resource.
type HelmReleaseSettingsArgs struct {
	// The backend storage driver for Helm. Values are: configmap, secret, memory, sql.
	Driver pulumi.StringPtrInput `pulumi:"driver"`
	// The path to the helm plugins directory.
	PluginsPath pulumi.StringPtrInput `pulumi:"pluginsPath"`
	// to test required args
	RequiredArg pulumi.StringInput `pulumi:"requiredArg"`
}

func (HelmReleaseSettingsArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*HelmReleaseSettings)(nil)).Elem()
}

func (i HelmReleaseSettingsArgs) ToHelmReleaseSettingsOutput() HelmReleaseSettingsOutput {
	return i.ToHelmReleaseSettingsOutputWithContext(context.Background())
}

func (i HelmReleaseSettingsArgs) ToHelmReleaseSettingsOutputWithContext(ctx context.Context) HelmReleaseSettingsOutput {
	return pulumi.ToOutputWithContext(ctx, i).(HelmReleaseSettingsOutput)
}

func (i HelmReleaseSettingsArgs) ToHelmReleaseSettingsPtrOutput() HelmReleaseSettingsPtrOutput {
	return i.ToHelmReleaseSettingsPtrOutputWithContext(context.Background())
}

func (i HelmReleaseSettingsArgs) ToHelmReleaseSettingsPtrOutputWithContext(ctx context.Context) HelmReleaseSettingsPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(HelmReleaseSettingsOutput).ToHelmReleaseSettingsPtrOutputWithContext(ctx)
}

// HelmReleaseSettingsPtrInput is an input type that accepts HelmReleaseSettingsArgs, HelmReleaseSettingsPtr and HelmReleaseSettingsPtrOutput values.
// You can construct a concrete instance of `HelmReleaseSettingsPtrInput` via:
//
//          HelmReleaseSettingsArgs{...}
//
//  or:
//
//          nil
type HelmReleaseSettingsPtrInput interface {
	pulumi.Input

	ToHelmReleaseSettingsPtrOutput() HelmReleaseSettingsPtrOutput
	ToHelmReleaseSettingsPtrOutputWithContext(context.Context) HelmReleaseSettingsPtrOutput
}

type helmReleaseSettingsPtrType HelmReleaseSettingsArgs

func HelmReleaseSettingsPtr(v *HelmReleaseSettingsArgs) HelmReleaseSettingsPtrInput {
	return (*helmReleaseSettingsPtrType)(v)
}

func (*helmReleaseSettingsPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**HelmReleaseSettings)(nil)).Elem()
}

func (i *helmReleaseSettingsPtrType) ToHelmReleaseSettingsPtrOutput() HelmReleaseSettingsPtrOutput {
	return i.ToHelmReleaseSettingsPtrOutputWithContext(context.Background())
}

func (i *helmReleaseSettingsPtrType) ToHelmReleaseSettingsPtrOutputWithContext(ctx context.Context) HelmReleaseSettingsPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(HelmReleaseSettingsPtrOutput)
}

// BETA FEATURE - Options to configure the Helm Release resource.
type HelmReleaseSettingsOutput struct{ *pulumi.OutputState }

func (HelmReleaseSettingsOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*HelmReleaseSettings)(nil)).Elem()
}

func (o HelmReleaseSettingsOutput) ToHelmReleaseSettingsOutput() HelmReleaseSettingsOutput {
	return o
}

func (o HelmReleaseSettingsOutput) ToHelmReleaseSettingsOutputWithContext(ctx context.Context) HelmReleaseSettingsOutput {
	return o
}

func (o HelmReleaseSettingsOutput) ToHelmReleaseSettingsPtrOutput() HelmReleaseSettingsPtrOutput {
	return o.ToHelmReleaseSettingsPtrOutputWithContext(context.Background())
}

func (o HelmReleaseSettingsOutput) ToHelmReleaseSettingsPtrOutputWithContext(ctx context.Context) HelmReleaseSettingsPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v HelmReleaseSettings) *HelmReleaseSettings {
		return &v
	}).(HelmReleaseSettingsPtrOutput)
}

// The backend storage driver for Helm. Values are: configmap, secret, memory, sql.
func (o HelmReleaseSettingsOutput) Driver() pulumi.StringPtrOutput {
	return o.ApplyT(func(v HelmReleaseSettings) *string { return v.Driver }).(pulumi.StringPtrOutput)
}

// The path to the helm plugins directory.
func (o HelmReleaseSettingsOutput) PluginsPath() pulumi.StringPtrOutput {
	return o.ApplyT(func(v HelmReleaseSettings) *string { return v.PluginsPath }).(pulumi.StringPtrOutput)
}

// to test required args
func (o HelmReleaseSettingsOutput) RequiredArg() pulumi.StringOutput {
	return o.ApplyT(func(v HelmReleaseSettings) string { return v.RequiredArg }).(pulumi.StringOutput)
}

type HelmReleaseSettingsPtrOutput struct{ *pulumi.OutputState }

func (HelmReleaseSettingsPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**HelmReleaseSettings)(nil)).Elem()
}

func (o HelmReleaseSettingsPtrOutput) ToHelmReleaseSettingsPtrOutput() HelmReleaseSettingsPtrOutput {
	return o
}

func (o HelmReleaseSettingsPtrOutput) ToHelmReleaseSettingsPtrOutputWithContext(ctx context.Context) HelmReleaseSettingsPtrOutput {
	return o
}

func (o HelmReleaseSettingsPtrOutput) Elem() HelmReleaseSettingsOutput {
	return o.ApplyT(func(v *HelmReleaseSettings) HelmReleaseSettings {
		if v != nil {
			return *v
		}
		var ret HelmReleaseSettings
		return ret
	}).(HelmReleaseSettingsOutput)
}

// The backend storage driver for Helm. Values are: configmap, secret, memory, sql.
func (o HelmReleaseSettingsPtrOutput) Driver() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *HelmReleaseSettings) *string {
		if v == nil {
			return nil
		}
		return v.Driver
	}).(pulumi.StringPtrOutput)
}

// The path to the helm plugins directory.
func (o HelmReleaseSettingsPtrOutput) PluginsPath() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *HelmReleaseSettings) *string {
		if v == nil {
			return nil
		}
		return v.PluginsPath
	}).(pulumi.StringPtrOutput)
}

// to test required args
func (o HelmReleaseSettingsPtrOutput) RequiredArg() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *HelmReleaseSettings) *string {
		if v == nil {
			return nil
		}
		return &v.RequiredArg
	}).(pulumi.StringPtrOutput)
}

// Options for tuning the Kubernetes client used by a Provider.
type KubeClientSettings struct {
	// Maximum burst for throttle. Default value is 10.
	Burst *int `pulumi:"burst"`
	// Maximum queries per second (QPS) to the API server from this client. Default value is 5.
	Qps     *float64            `pulumi:"qps"`
	RecTest *KubeClientSettings `pulumi:"recTest"`
}

// KubeClientSettingsInput is an input type that accepts KubeClientSettingsArgs and KubeClientSettingsOutput values.
// You can construct a concrete instance of `KubeClientSettingsInput` via:
//
//          KubeClientSettingsArgs{...}
type KubeClientSettingsInput interface {
	pulumi.Input

	ToKubeClientSettingsOutput() KubeClientSettingsOutput
	ToKubeClientSettingsOutputWithContext(context.Context) KubeClientSettingsOutput
}

// Options for tuning the Kubernetes client used by a Provider.
type KubeClientSettingsArgs struct {
	// Maximum burst for throttle. Default value is 10.
	Burst pulumi.IntPtrInput `pulumi:"burst"`
	// Maximum queries per second (QPS) to the API server from this client. Default value is 5.
	Qps     pulumi.Float64PtrInput     `pulumi:"qps"`
	RecTest KubeClientSettingsPtrInput `pulumi:"recTest"`
}

func (KubeClientSettingsArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*KubeClientSettings)(nil)).Elem()
}

func (i KubeClientSettingsArgs) ToKubeClientSettingsOutput() KubeClientSettingsOutput {
	return i.ToKubeClientSettingsOutputWithContext(context.Background())
}

func (i KubeClientSettingsArgs) ToKubeClientSettingsOutputWithContext(ctx context.Context) KubeClientSettingsOutput {
	return pulumi.ToOutputWithContext(ctx, i).(KubeClientSettingsOutput)
}

func (i KubeClientSettingsArgs) ToKubeClientSettingsPtrOutput() KubeClientSettingsPtrOutput {
	return i.ToKubeClientSettingsPtrOutputWithContext(context.Background())
}

func (i KubeClientSettingsArgs) ToKubeClientSettingsPtrOutputWithContext(ctx context.Context) KubeClientSettingsPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(KubeClientSettingsOutput).ToKubeClientSettingsPtrOutputWithContext(ctx)
}

// KubeClientSettingsPtrInput is an input type that accepts KubeClientSettingsArgs, KubeClientSettingsPtr and KubeClientSettingsPtrOutput values.
// You can construct a concrete instance of `KubeClientSettingsPtrInput` via:
//
//          KubeClientSettingsArgs{...}
//
//  or:
//
//          nil
type KubeClientSettingsPtrInput interface {
	pulumi.Input

	ToKubeClientSettingsPtrOutput() KubeClientSettingsPtrOutput
	ToKubeClientSettingsPtrOutputWithContext(context.Context) KubeClientSettingsPtrOutput
}

type kubeClientSettingsPtrType KubeClientSettingsArgs

func KubeClientSettingsPtr(v *KubeClientSettingsArgs) KubeClientSettingsPtrInput {
	return (*kubeClientSettingsPtrType)(v)
}

func (*kubeClientSettingsPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**KubeClientSettings)(nil)).Elem()
}

func (i *kubeClientSettingsPtrType) ToKubeClientSettingsPtrOutput() KubeClientSettingsPtrOutput {
	return i.ToKubeClientSettingsPtrOutputWithContext(context.Background())
}

func (i *kubeClientSettingsPtrType) ToKubeClientSettingsPtrOutputWithContext(ctx context.Context) KubeClientSettingsPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(KubeClientSettingsPtrOutput)
}

// Options for tuning the Kubernetes client used by a Provider.
type KubeClientSettingsOutput struct{ *pulumi.OutputState }

func (KubeClientSettingsOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*KubeClientSettings)(nil)).Elem()
}

func (o KubeClientSettingsOutput) ToKubeClientSettingsOutput() KubeClientSettingsOutput {
	return o
}

func (o KubeClientSettingsOutput) ToKubeClientSettingsOutputWithContext(ctx context.Context) KubeClientSettingsOutput {
	return o
}

func (o KubeClientSettingsOutput) ToKubeClientSettingsPtrOutput() KubeClientSettingsPtrOutput {
	return o.ToKubeClientSettingsPtrOutputWithContext(context.Background())
}

func (o KubeClientSettingsOutput) ToKubeClientSettingsPtrOutputWithContext(ctx context.Context) KubeClientSettingsPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v KubeClientSettings) *KubeClientSettings {
		return &v
	}).(KubeClientSettingsPtrOutput)
}

// Maximum burst for throttle. Default value is 10.
func (o KubeClientSettingsOutput) Burst() pulumi.IntPtrOutput {
	return o.ApplyT(func(v KubeClientSettings) *int { return v.Burst }).(pulumi.IntPtrOutput)
}

// Maximum queries per second (QPS) to the API server from this client. Default value is 5.
func (o KubeClientSettingsOutput) Qps() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v KubeClientSettings) *float64 { return v.Qps }).(pulumi.Float64PtrOutput)
}

func (o KubeClientSettingsOutput) RecTest() KubeClientSettingsPtrOutput {
	return o.ApplyT(func(v KubeClientSettings) *KubeClientSettings { return v.RecTest }).(KubeClientSettingsPtrOutput)
}

type KubeClientSettingsPtrOutput struct{ *pulumi.OutputState }

func (KubeClientSettingsPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**KubeClientSettings)(nil)).Elem()
}

func (o KubeClientSettingsPtrOutput) ToKubeClientSettingsPtrOutput() KubeClientSettingsPtrOutput {
	return o
}

func (o KubeClientSettingsPtrOutput) ToKubeClientSettingsPtrOutputWithContext(ctx context.Context) KubeClientSettingsPtrOutput {
	return o
}

func (o KubeClientSettingsPtrOutput) Elem() KubeClientSettingsOutput {
	return o.ApplyT(func(v *KubeClientSettings) KubeClientSettings {
		if v != nil {
			return *v
		}
		var ret KubeClientSettings
		return ret
	}).(KubeClientSettingsOutput)
}

// Maximum burst for throttle. Default value is 10.
func (o KubeClientSettingsPtrOutput) Burst() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *KubeClientSettings) *int {
		if v == nil {
			return nil
		}
		return v.Burst
	}).(pulumi.IntPtrOutput)
}

// Maximum queries per second (QPS) to the API server from this client. Default value is 5.
func (o KubeClientSettingsPtrOutput) Qps() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v *KubeClientSettings) *float64 {
		if v == nil {
			return nil
		}
		return v.Qps
	}).(pulumi.Float64PtrOutput)
}

func (o KubeClientSettingsPtrOutput) RecTest() KubeClientSettingsPtrOutput {
	return o.ApplyT(func(v *KubeClientSettings) *KubeClientSettings {
		if v == nil {
			return nil
		}
		return v.RecTest
	}).(KubeClientSettingsPtrOutput)
}

// Make sure that defaults propagate through types
type LayeredType struct {
	// The answer to the question
	Answer *float64            `pulumi:"answer"`
	Other  HelmReleaseSettings `pulumi:"other"`
	// Test how plain types interact
	PlainOther *HelmReleaseSettings `pulumi:"plainOther"`
	// The question already answered
	Question  *string      `pulumi:"question"`
	Recursive *LayeredType `pulumi:"recursive"`
	// To ask and answer
	Thinker string `pulumi:"thinker"`
}

// LayeredTypeInput is an input type that accepts LayeredTypeArgs and LayeredTypeOutput values.
// You can construct a concrete instance of `LayeredTypeInput` via:
//
//          LayeredTypeArgs{...}
type LayeredTypeInput interface {
	pulumi.Input

	ToLayeredTypeOutput() LayeredTypeOutput
	ToLayeredTypeOutputWithContext(context.Context) LayeredTypeOutput
}

// Make sure that defaults propagate through types
type LayeredTypeArgs struct {
	// The answer to the question
	Answer pulumi.Float64PtrInput   `pulumi:"answer"`
	Other  HelmReleaseSettingsInput `pulumi:"other"`
	// Test how plain types interact
	PlainOther *HelmReleaseSettingsArgs `pulumi:"plainOther"`
	// The question already answered
	Question  pulumi.StringPtrInput `pulumi:"question"`
	Recursive LayeredTypePtrInput   `pulumi:"recursive"`
	// To ask and answer
	Thinker pulumi.StringInput `pulumi:"thinker"`
}

func (LayeredTypeArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LayeredType)(nil)).Elem()
}

func (i LayeredTypeArgs) ToLayeredTypeOutput() LayeredTypeOutput {
	return i.ToLayeredTypeOutputWithContext(context.Background())
}

func (i LayeredTypeArgs) ToLayeredTypeOutputWithContext(ctx context.Context) LayeredTypeOutput {
	return pulumi.ToOutputWithContext(ctx, i).(LayeredTypeOutput)
}

func (i LayeredTypeArgs) ToLayeredTypePtrOutput() LayeredTypePtrOutput {
	return i.ToLayeredTypePtrOutputWithContext(context.Background())
}

func (i LayeredTypeArgs) ToLayeredTypePtrOutputWithContext(ctx context.Context) LayeredTypePtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(LayeredTypeOutput).ToLayeredTypePtrOutputWithContext(ctx)
}

// LayeredTypePtrInput is an input type that accepts LayeredTypeArgs, LayeredTypePtr and LayeredTypePtrOutput values.
// You can construct a concrete instance of `LayeredTypePtrInput` via:
//
//          LayeredTypeArgs{...}
//
//  or:
//
//          nil
type LayeredTypePtrInput interface {
	pulumi.Input

	ToLayeredTypePtrOutput() LayeredTypePtrOutput
	ToLayeredTypePtrOutputWithContext(context.Context) LayeredTypePtrOutput
}

type layeredTypePtrType LayeredTypeArgs

func LayeredTypePtr(v *LayeredTypeArgs) LayeredTypePtrInput {
	return (*layeredTypePtrType)(v)
}

func (*layeredTypePtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**LayeredType)(nil)).Elem()
}

func (i *layeredTypePtrType) ToLayeredTypePtrOutput() LayeredTypePtrOutput {
	return i.ToLayeredTypePtrOutputWithContext(context.Background())
}

func (i *layeredTypePtrType) ToLayeredTypePtrOutputWithContext(ctx context.Context) LayeredTypePtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(LayeredTypePtrOutput)
}

// Make sure that defaults propagate through types
type LayeredTypeOutput struct{ *pulumi.OutputState }

func (LayeredTypeOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LayeredType)(nil)).Elem()
}

func (o LayeredTypeOutput) ToLayeredTypeOutput() LayeredTypeOutput {
	return o
}

func (o LayeredTypeOutput) ToLayeredTypeOutputWithContext(ctx context.Context) LayeredTypeOutput {
	return o
}

func (o LayeredTypeOutput) ToLayeredTypePtrOutput() LayeredTypePtrOutput {
	return o.ToLayeredTypePtrOutputWithContext(context.Background())
}

func (o LayeredTypeOutput) ToLayeredTypePtrOutputWithContext(ctx context.Context) LayeredTypePtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v LayeredType) *LayeredType {
		return &v
	}).(LayeredTypePtrOutput)
}

// The answer to the question
func (o LayeredTypeOutput) Answer() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v LayeredType) *float64 { return v.Answer }).(pulumi.Float64PtrOutput)
}

func (o LayeredTypeOutput) Other() HelmReleaseSettingsOutput {
	return o.ApplyT(func(v LayeredType) HelmReleaseSettings { return v.Other }).(HelmReleaseSettingsOutput)
}

// Test how plain types interact
func (o LayeredTypeOutput) PlainOther() HelmReleaseSettingsPtrOutput {
	return o.ApplyT(func(v LayeredType) *HelmReleaseSettings { return v.PlainOther }).(HelmReleaseSettingsPtrOutput)
}

// The question already answered
func (o LayeredTypeOutput) Question() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LayeredType) *string { return v.Question }).(pulumi.StringPtrOutput)
}

func (o LayeredTypeOutput) Recursive() LayeredTypePtrOutput {
	return o.ApplyT(func(v LayeredType) *LayeredType { return v.Recursive }).(LayeredTypePtrOutput)
}

// To ask and answer
func (o LayeredTypeOutput) Thinker() pulumi.StringOutput {
	return o.ApplyT(func(v LayeredType) string { return v.Thinker }).(pulumi.StringOutput)
}

type LayeredTypePtrOutput struct{ *pulumi.OutputState }

func (LayeredTypePtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**LayeredType)(nil)).Elem()
}

func (o LayeredTypePtrOutput) ToLayeredTypePtrOutput() LayeredTypePtrOutput {
	return o
}

func (o LayeredTypePtrOutput) ToLayeredTypePtrOutputWithContext(ctx context.Context) LayeredTypePtrOutput {
	return o
}

func (o LayeredTypePtrOutput) Elem() LayeredTypeOutput {
	return o.ApplyT(func(v *LayeredType) LayeredType {
		if v != nil {
			return *v
		}
		var ret LayeredType
		return ret
	}).(LayeredTypeOutput)
}

// The answer to the question
func (o LayeredTypePtrOutput) Answer() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v *LayeredType) *float64 {
		if v == nil {
			return nil
		}
		return v.Answer
	}).(pulumi.Float64PtrOutput)
}

func (o LayeredTypePtrOutput) Other() HelmReleaseSettingsPtrOutput {
	return o.ApplyT(func(v *LayeredType) *HelmReleaseSettings {
		if v == nil {
			return nil
		}
		return &v.Other
	}).(HelmReleaseSettingsPtrOutput)
}

// Test how plain types interact
func (o LayeredTypePtrOutput) PlainOther() HelmReleaseSettingsPtrOutput {
	return o.ApplyT(func(v *LayeredType) *HelmReleaseSettings {
		if v == nil {
			return nil
		}
		return v.PlainOther
	}).(HelmReleaseSettingsPtrOutput)
}

// The question already answered
func (o LayeredTypePtrOutput) Question() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *LayeredType) *string {
		if v == nil {
			return nil
		}
		return v.Question
	}).(pulumi.StringPtrOutput)
}

func (o LayeredTypePtrOutput) Recursive() LayeredTypePtrOutput {
	return o.ApplyT(func(v *LayeredType) *LayeredType {
		if v == nil {
			return nil
		}
		return v.Recursive
	}).(LayeredTypePtrOutput)
}

// To ask and answer
func (o LayeredTypePtrOutput) Thinker() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *LayeredType) *string {
		if v == nil {
			return nil
		}
		return &v.Thinker
	}).(pulumi.StringPtrOutput)
}

// A test for namespaces (mod main)
type Typ struct {
	Mod1 *mod1.Typ `pulumi:"mod1"`
	Mod2 *mod2.Typ `pulumi:"mod2"`
	Val  *string   `pulumi:"val"`
}

// TypInput is an input type that accepts TypArgs and TypOutput values.
// You can construct a concrete instance of `TypInput` via:
//
//          TypArgs{...}
type TypInput interface {
	pulumi.Input

	ToTypOutput() TypOutput
	ToTypOutputWithContext(context.Context) TypOutput
}

// A test for namespaces (mod main)
type TypArgs struct {
	Mod1 mod1.TypPtrInput      `pulumi:"mod1"`
	Mod2 mod2.TypPtrInput      `pulumi:"mod2"`
	Val  pulumi.StringPtrInput `pulumi:"val"`
}

func (TypArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*Typ)(nil)).Elem()
}

func (i TypArgs) ToTypOutput() TypOutput {
	return i.ToTypOutputWithContext(context.Background())
}

func (i TypArgs) ToTypOutputWithContext(ctx context.Context) TypOutput {
	return pulumi.ToOutputWithContext(ctx, i).(TypOutput)
}

func (i TypArgs) ToTypPtrOutput() TypPtrOutput {
	return i.ToTypPtrOutputWithContext(context.Background())
}

func (i TypArgs) ToTypPtrOutputWithContext(ctx context.Context) TypPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(TypOutput).ToTypPtrOutputWithContext(ctx)
}

// TypPtrInput is an input type that accepts TypArgs, TypPtr and TypPtrOutput values.
// You can construct a concrete instance of `TypPtrInput` via:
//
//          TypArgs{...}
//
//  or:
//
//          nil
type TypPtrInput interface {
	pulumi.Input

	ToTypPtrOutput() TypPtrOutput
	ToTypPtrOutputWithContext(context.Context) TypPtrOutput
}

type typPtrType TypArgs

func TypPtr(v *TypArgs) TypPtrInput {
	return (*typPtrType)(v)
}

func (*typPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**Typ)(nil)).Elem()
}

func (i *typPtrType) ToTypPtrOutput() TypPtrOutput {
	return i.ToTypPtrOutputWithContext(context.Background())
}

func (i *typPtrType) ToTypPtrOutputWithContext(ctx context.Context) TypPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(TypPtrOutput)
}

// A test for namespaces (mod main)
type TypOutput struct{ *pulumi.OutputState }

func (TypOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Typ)(nil)).Elem()
}

func (o TypOutput) ToTypOutput() TypOutput {
	return o
}

func (o TypOutput) ToTypOutputWithContext(ctx context.Context) TypOutput {
	return o
}

func (o TypOutput) ToTypPtrOutput() TypPtrOutput {
	return o.ToTypPtrOutputWithContext(context.Background())
}

func (o TypOutput) ToTypPtrOutputWithContext(ctx context.Context) TypPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v Typ) *Typ {
		return &v
	}).(TypPtrOutput)
}

func (o TypOutput) Mod1() mod1.TypPtrOutput {
	return o.ApplyT(func(v Typ) *mod1.Typ { return v.Mod1 }).(mod1.TypPtrOutput)
}

func (o TypOutput) Mod2() mod2.TypPtrOutput {
	return o.ApplyT(func(v Typ) *mod2.Typ { return v.Mod2 }).(mod2.TypPtrOutput)
}

func (o TypOutput) Val() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Typ) *string { return v.Val }).(pulumi.StringPtrOutput)
}

type TypPtrOutput struct{ *pulumi.OutputState }

func (TypPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**Typ)(nil)).Elem()
}

func (o TypPtrOutput) ToTypPtrOutput() TypPtrOutput {
	return o
}

func (o TypPtrOutput) ToTypPtrOutputWithContext(ctx context.Context) TypPtrOutput {
	return o
}

func (o TypPtrOutput) Elem() TypOutput {
	return o.ApplyT(func(v *Typ) Typ {
		if v != nil {
			return *v
		}
		var ret Typ
		return ret
	}).(TypOutput)
}

func (o TypPtrOutput) Mod1() mod1.TypPtrOutput {
	return o.ApplyT(func(v *Typ) *mod1.Typ {
		if v == nil {
			return nil
		}
		return v.Mod1
	}).(mod1.TypPtrOutput)
}

func (o TypPtrOutput) Mod2() mod2.TypPtrOutput {
	return o.ApplyT(func(v *Typ) *mod2.Typ {
		if v == nil {
			return nil
		}
		return v.Mod2
	}).(mod2.TypPtrOutput)
}

func (o TypPtrOutput) Val() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Typ) *string {
		if v == nil {
			return nil
		}
		return v.Val
	}).(pulumi.StringPtrOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*HelmReleaseSettingsInput)(nil)).Elem(), HelmReleaseSettingsArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*HelmReleaseSettingsPtrInput)(nil)).Elem(), HelmReleaseSettingsArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*KubeClientSettingsInput)(nil)).Elem(), KubeClientSettingsArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*KubeClientSettingsPtrInput)(nil)).Elem(), KubeClientSettingsArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*LayeredTypeInput)(nil)).Elem(), LayeredTypeArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*LayeredTypePtrInput)(nil)).Elem(), LayeredTypeArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*TypInput)(nil)).Elem(), TypArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*TypPtrInput)(nil)).Elem(), TypArgs{})
	pulumi.RegisterOutputType(HelmReleaseSettingsOutput{})
	pulumi.RegisterOutputType(HelmReleaseSettingsPtrOutput{})
	pulumi.RegisterOutputType(KubeClientSettingsOutput{})
	pulumi.RegisterOutputType(KubeClientSettingsPtrOutput{})
	pulumi.RegisterOutputType(LayeredTypeOutput{})
	pulumi.RegisterOutputType(LayeredTypePtrOutput{})
	pulumi.RegisterOutputType(TypOutput{})
	pulumi.RegisterOutputType(TypPtrOutput{})
}
