# Generated by the protocol buffer compiler.  DO NOT EDIT!
# sources: api/debug/debug.proto
# plugin: python-betterproto
# This file has been @generated
from dataclasses import dataclass
from datetime import timedelta
from typing import (
    TYPE_CHECKING,
    AsyncIterator,
    Dict,
    Iterator,
    Optional,
)

import betterproto
import betterproto.lib.google.protobuf as betterproto_lib_google_protobuf
import grpc
from betterproto.grpc.grpcio_server import ServicerBase

from .. import pps as _pps__


if TYPE_CHECKING:
    import grpc


@dataclass(eq=False, repr=False)
class ProfileRequest(betterproto.Message):
    profile: "Profile" = betterproto.message_field(1)
    filter: "Filter" = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class Profile(betterproto.Message):
    name: str = betterproto.string_field(1)
    duration: timedelta = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class Filter(betterproto.Message):
    pachd: bool = betterproto.bool_field(1, group="filter")
    pipeline: "_pps__.Pipeline" = betterproto.message_field(2, group="filter")
    worker: "Worker" = betterproto.message_field(3, group="filter")
    database: bool = betterproto.bool_field(4, group="filter")


@dataclass(eq=False, repr=False)
class Worker(betterproto.Message):
    pod: str = betterproto.string_field(1)
    redirected: bool = betterproto.bool_field(2)


@dataclass(eq=False, repr=False)
class BinaryRequest(betterproto.Message):
    filter: "Filter" = betterproto.message_field(1)


@dataclass(eq=False, repr=False)
class DumpRequest(betterproto.Message):
    filter: "Filter" = betterproto.message_field(1)
    limit: int = betterproto.int64_field(2)
    """
    Limit sets the limit for the number of commits / jobs that are returned for
    each repo / pipeline in the dump.
    """


class DebugStub:
    def __init__(self, channel: "grpc.Channel"):
        self.__rpc_profile = channel.unary_stream(
            "/debug_v2.Debug/Profile",
            request_serializer=ProfileRequest.SerializeToString,
            response_deserializer=betterproto_lib_google_protobuf.BytesValue.FromString,
        )
        self.__rpc_binary = channel.unary_stream(
            "/debug_v2.Debug/Binary",
            request_serializer=BinaryRequest.SerializeToString,
            response_deserializer=betterproto_lib_google_protobuf.BytesValue.FromString,
        )
        self.__rpc_dump = channel.unary_stream(
            "/debug_v2.Debug/Dump",
            request_serializer=DumpRequest.SerializeToString,
            response_deserializer=betterproto_lib_google_protobuf.BytesValue.FromString,
        )

    def profile(
        self, *, profile: "Profile" = None, filter: "Filter" = None
    ) -> Iterator["betterproto_lib_google_protobuf.BytesValue"]:
        request = ProfileRequest()
        if profile is not None:
            request.profile = profile
        if filter is not None:
            request.filter = filter

        for response in self.__rpc_profile(request):
            yield response

    def binary(
        self, *, filter: "Filter" = None
    ) -> Iterator["betterproto_lib_google_protobuf.BytesValue"]:
        request = BinaryRequest()
        if filter is not None:
            request.filter = filter

        for response in self.__rpc_binary(request):
            yield response

    def dump(
        self, *, filter: "Filter" = None, limit: int = 0
    ) -> Iterator["betterproto_lib_google_protobuf.BytesValue"]:
        request = DumpRequest()
        if filter is not None:
            request.filter = filter
        request.limit = limit

        for response in self.__rpc_dump(request):
            yield response


class DebugBase(ServicerBase):
    def profile(
        self, profile: "Profile", filter: "Filter", context: "grpc.ServicerContext"
    ) -> Iterator["betterproto_lib_google_protobuf.BytesValue"]:
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def binary(
        self, filter: "Filter", context: "grpc.ServicerContext"
    ) -> Iterator["betterproto_lib_google_protobuf.BytesValue"]:
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def dump(
        self, filter: "Filter", limit: int, context: "grpc.ServicerContext"
    ) -> Iterator["betterproto_lib_google_protobuf.BytesValue"]:
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    __proto_path__ = "debug_v2.Debug"

    @property
    def __rpc_methods__(self):
        return {
            "Profile": grpc.unary_stream_rpc_method_handler(
                self.profile,
                request_deserializer=ProfileRequest.FromString,
                response_serializer=ProfileRequest.SerializeToString,
            ),
            "Binary": grpc.unary_stream_rpc_method_handler(
                self.binary,
                request_deserializer=BinaryRequest.FromString,
                response_serializer=BinaryRequest.SerializeToString,
            ),
            "Dump": grpc.unary_stream_rpc_method_handler(
                self.dump,
                request_deserializer=DumpRequest.FromString,
                response_serializer=DumpRequest.SerializeToString,
            ),
        }
