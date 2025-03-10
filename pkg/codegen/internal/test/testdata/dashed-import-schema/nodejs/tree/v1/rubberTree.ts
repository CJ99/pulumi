// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import { input as inputs, output as outputs, enums } from "../../types";
import * as utilities from "../../utilities";

export class RubberTree extends pulumi.CustomResource {
    /**
     * Get an existing RubberTree resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param state Any extra arguments used during the lookup.
     * @param opts Optional settings to control the behavior of the CustomResource.
     */
    public static get(name: string, id: pulumi.Input<pulumi.ID>, state?: RubberTreeState, opts?: pulumi.CustomResourceOptions): RubberTree {
        return new RubberTree(name, <any>state, { ...opts, id: id });
    }

    /** @internal */
    public static readonly __pulumiType = 'plant:tree/v1:RubberTree';

    /**
     * Returns true if the given object is an instance of RubberTree.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is RubberTree {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === RubberTree.__pulumiType;
    }

    public readonly container!: pulumi.Output<outputs.Container | undefined>;
    public readonly diameter!: pulumi.Output<enums.tree.v1.Diameter>;
    public readonly farm!: pulumi.Output<enums.tree.v1.Farm | string | undefined>;
    public readonly size!: pulumi.Output<enums.tree.v1.TreeSize | undefined>;
    public readonly type!: pulumi.Output<enums.tree.v1.RubberTreeVariety>;

    /**
     * Create a RubberTree resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args: RubberTreeArgs, opts?: pulumi.CustomResourceOptions)
    constructor(name: string, argsOrState?: RubberTreeArgs | RubberTreeState, opts?: pulumi.CustomResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (opts.id) {
            const state = argsOrState as RubberTreeState | undefined;
            resourceInputs["farm"] = state ? state.farm : undefined;
        } else {
            const args = argsOrState as RubberTreeArgs | undefined;
            if ((!args || args.diameter === undefined) && !opts.urn) {
                throw new Error("Missing required property 'diameter'");
            }
            if ((!args || args.type === undefined) && !opts.urn) {
                throw new Error("Missing required property 'type'");
            }
            resourceInputs["container"] = args ? (args.container ? pulumi.output(args.container).apply(inputs.containerArgsProvideDefaults) : undefined) : undefined;
            resourceInputs["diameter"] = (args ? args.diameter : undefined) ?? 6;
            resourceInputs["farm"] = (args ? args.farm : undefined) ?? "(unknown)";
            resourceInputs["size"] = (args ? args.size : undefined) ?? "medium";
            resourceInputs["type"] = (args ? args.type : undefined) ?? "Burgundy";
        }
        if (!opts.version) {
            opts = pulumi.mergeOptions(opts, { version: utilities.getVersion()});
        }
        super(RubberTree.__pulumiType, name, resourceInputs, opts);
    }
}

export interface RubberTreeState {
    farm?: pulumi.Input<enums.tree.v1.Farm | string>;
}

/**
 * The set of arguments for constructing a RubberTree resource.
 */
export interface RubberTreeArgs {
    container?: pulumi.Input<inputs.ContainerArgs>;
    diameter: pulumi.Input<enums.tree.v1.Diameter>;
    farm?: pulumi.Input<enums.tree.v1.Farm | string>;
    size?: pulumi.Input<enums.tree.v1.TreeSize>;
    type: pulumi.Input<enums.tree.v1.RubberTreeVariety>;
}
