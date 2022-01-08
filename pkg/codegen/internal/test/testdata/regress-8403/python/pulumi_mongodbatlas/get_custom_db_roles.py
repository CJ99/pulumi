# coding=utf-8
# *** WARNING: this file was generated by test. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from . import _utilities
from . import outputs

__all__ = [
    'GetCustomDbRolesResult',
    'AwaitableGetCustomDbRolesResult',
    'get_custom_db_roles',
]

@pulumi.output_type
class GetCustomDbRolesResult:
    def __init__(__self__, result=None):
        if result and not isinstance(result, dict):
            raise TypeError("Expected argument 'result' to be a dict")
        pulumi.set(__self__, "result", result)

    @property
    @pulumi.getter
    def result(self) -> Optional['outputs.GetCustomDbRolesResult']:
        return pulumi.get(self, "result")


class AwaitableGetCustomDbRolesResult(GetCustomDbRolesResult):
    # pylint: disable=using-constant-test
    def __await__(self):
        if False:
            yield self
        return GetCustomDbRolesResult(
            result=self.result)


def get_custom_db_roles(opts: Optional[pulumi.InvokeOptions] = None) -> AwaitableGetCustomDbRolesResult:
    """
    Use this data source to access information about an existing resource.
    """
    __args__ = dict()
    if opts is None:
        opts = pulumi.InvokeOptions()
    if opts.version is None:
        opts.version = _utilities.get_version()
    __ret__ = pulumi.runtime.invoke('mongodbatlas::getCustomDbRoles', __args__, opts=opts, typ=GetCustomDbRolesResult).value

    return AwaitableGetCustomDbRolesResult(
        result=__ret__.result)
