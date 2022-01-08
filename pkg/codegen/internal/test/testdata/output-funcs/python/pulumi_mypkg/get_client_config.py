# coding=utf-8
# *** WARNING: this file was generated by test. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from . import _utilities

__all__ = [
    'GetClientConfigResult',
    'AwaitableGetClientConfigResult',
    'get_client_config',
]

@pulumi.output_type
class GetClientConfigResult:
    """
    Configuration values returned by getClientConfig.
    """
    def __init__(__self__, client_id=None, object_id=None, subscription_id=None, tenant_id=None):
        if client_id and not isinstance(client_id, str):
            raise TypeError("Expected argument 'client_id' to be a str")
        pulumi.set(__self__, "client_id", client_id)
        if object_id and not isinstance(object_id, str):
            raise TypeError("Expected argument 'object_id' to be a str")
        pulumi.set(__self__, "object_id", object_id)
        if subscription_id and not isinstance(subscription_id, str):
            raise TypeError("Expected argument 'subscription_id' to be a str")
        pulumi.set(__self__, "subscription_id", subscription_id)
        if tenant_id and not isinstance(tenant_id, str):
            raise TypeError("Expected argument 'tenant_id' to be a str")
        pulumi.set(__self__, "tenant_id", tenant_id)

    @property
    @pulumi.getter(name="clientId")
    def client_id(self) -> str:
        """
        Azure Client ID (Application Object ID).
        """
        return pulumi.get(self, "client_id")

    @property
    @pulumi.getter(name="objectId")
    def object_id(self) -> str:
        """
        Azure Object ID of the current user or service principal.
        """
        return pulumi.get(self, "object_id")

    @property
    @pulumi.getter(name="subscriptionId")
    def subscription_id(self) -> str:
        """
        Azure Subscription ID
        """
        return pulumi.get(self, "subscription_id")

    @property
    @pulumi.getter(name="tenantId")
    def tenant_id(self) -> str:
        """
        Azure Tenant ID
        """
        return pulumi.get(self, "tenant_id")


class AwaitableGetClientConfigResult(GetClientConfigResult):
    # pylint: disable=using-constant-test
    def __await__(self):
        if False:
            yield self
        return GetClientConfigResult(
            client_id=self.client_id,
            object_id=self.object_id,
            subscription_id=self.subscription_id,
            tenant_id=self.tenant_id)


def get_client_config(opts: Optional[pulumi.InvokeOptions] = None) -> AwaitableGetClientConfigResult:
    """
    Failing example taken from azure-native. Original doc: Use this function to access the current configuration of the native Azure provider.
    """
    __args__ = dict()
    if opts is None:
        opts = pulumi.InvokeOptions()
    if opts.version is None:
        opts.version = _utilities.get_version()
    __ret__ = pulumi.runtime.invoke('mypkg::getClientConfig', __args__, opts=opts, typ=GetClientConfigResult).value

    return AwaitableGetClientConfigResult(
        client_id=__ret__.client_id,
        object_id=__ret__.object_id,
        subscription_id=__ret__.subscription_id,
        tenant_id=__ret__.tenant_id)
