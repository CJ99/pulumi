// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "./utilities";

import {Resource} from "./index";

export class BarResource extends pulumi.ComponentResource {
    /** @internal */
    public static readonly __pulumiType = 'bar::BarResource';

    /**
     * Returns true if the given object is an instance of BarResource.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is BarResource {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === BarResource.__pulumiType;
    }

    public readonly foo!: pulumi.Output<Resource | undefined>;

    /**
     * Create a BarResource resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args?: BarResourceArgs, opts?: pulumi.ComponentResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            resourceInputs["foo"] = args ? args.foo : undefined;
        } else {
            resourceInputs["foo"] = undefined /*out*/;
        }
        if (!opts.version) {
            opts = pulumi.mergeOptions(opts, { version: utilities.getVersion()});
        }
        super(BarResource.__pulumiType, name, resourceInputs, opts, true /*remote*/);
    }
}

/**
 * The set of arguments for constructing a BarResource resource.
 */
export interface BarResourceArgs {
    foo?: pulumi.Input<Resource>;
}
