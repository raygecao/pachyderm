# Generated by the protocol buffer compiler.  DO NOT EDIT!
# sources: api/logs/logs.proto
# plugin: python-betterproto
# This file has been @generated
from dataclasses import dataclass
from datetime import datetime
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

from .. import pps as _pps__


if TYPE_CHECKING:
    import grpc


class LogLevel(betterproto.Enum):
    LOG_LEVEL_DEBUG = 0
    LOG_LEVEL_INFO = 1
    LOG_LEVEL_ERROR = 2


class LogFormat(betterproto.Enum):
    LOG_FORMAT_UNKNOWN = 0
    """error"""

    LOG_FORMAT_VERBATIM_WITH_TIMESTAMP = 1
    LOG_FORMAT_PARSED_JSON = 2
    LOG_FORMAT_PPS_LOGMESSAGE = 3


@dataclass(eq=False, repr=False)
class LogQuery(betterproto.Message):
    user: "UserLogQuery" = betterproto.message_field(1, group="query_type")
    admin: "AdminLogQuery" = betterproto.message_field(2, group="query_type")


@dataclass(eq=False, repr=False)
class AdminLogQuery(betterproto.Message):
    logql: str = betterproto.string_field(1, group="admin_type")
    """Arbitrary LogQL query"""

    pod: str = betterproto.string_field(2, group="admin_type")
    """A pod's logs (all containers)"""

    pod_container: "PodContainer" = betterproto.message_field(3, group="admin_type")
    """One container"""

    app: str = betterproto.string_field(4, group="admin_type")
    """One "app" (logql -> {app=X})"""

    master: "PipelineLogQuery" = betterproto.message_field(5, group="admin_type")
    """All master worker lines from a pipeline"""

    storage: "PipelineLogQuery" = betterproto.message_field(6, group="admin_type")
    """All storage container lines from a pipeline"""

    user: "UserLogQuery" = betterproto.message_field(7, group="admin_type")
    """All worker lines from a pipeline/job"""


@dataclass(eq=False, repr=False)
class PodContainer(betterproto.Message):
    pod: str = betterproto.string_field(1)
    container: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class UserLogQuery(betterproto.Message):
    """Only returns "user" logs"""

    project: str = betterproto.string_field(1, group="user_type")
    """All pipelines in the project"""

    pipeline: "PipelineLogQuery" = betterproto.message_field(2, group="user_type")
    """One pipeline in a project"""

    datum: str = betterproto.string_field(3, group="user_type")
    """One datum."""

    job: str = betterproto.string_field(4, group="user_type")
    """One job, across pipelines and projects"""

    pipeline_job: "PipelineJobLogQuery" = betterproto.message_field(
        5, group="user_type"
    )
    """One job in one pipeline"""


@dataclass(eq=False, repr=False)
class PipelineLogQuery(betterproto.Message):
    project: str = betterproto.string_field(1)
    pipeline: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class PipelineJobLogQuery(betterproto.Message):
    pipeline: "PipelineLogQuery" = betterproto.message_field(1)
    job: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class PipelineDatumLogQuery(betterproto.Message):
    pipeline: "PipelineLogQuery" = betterproto.message_field(1)
    datum: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class LogFilter(betterproto.Message):
    time_range: "TimeRangeLogFilter" = betterproto.message_field(1)
    limit: int = betterproto.uint64_field(2)
    regex: "RegexLogFilter" = betterproto.message_field(3)
    level: "LogLevel" = betterproto.enum_field(4)
    """
    Minimum log level to return; worker will always run at level debug, but
    setting INFO here restores original behavior
    """


@dataclass(eq=False, repr=False)
class TimeRangeLogFilter(betterproto.Message):
    from_: datetime = betterproto.message_field(1)
    """Can be null"""

    until: datetime = betterproto.message_field(2)
    """Can be null"""


@dataclass(eq=False, repr=False)
class RegexLogFilter(betterproto.Message):
    pattern: str = betterproto.string_field(1)
    negate: bool = betterproto.bool_field(2)


@dataclass(eq=False, repr=False)
class GetLogsRequest(betterproto.Message):
    query: "LogQuery" = betterproto.message_field(1)
    filter: "LogFilter" = betterproto.message_field(2)
    tail: bool = betterproto.bool_field(3)
    want_paging_hint: bool = betterproto.bool_field(4)
    log_format: "LogFormat" = betterproto.enum_field(5)


@dataclass(eq=False, repr=False)
class GetLogsResponse(betterproto.Message):
    paging_hint: "PagingHint" = betterproto.message_field(1, group="response_type")
    log: "LogMessage" = betterproto.message_field(2, group="response_type")


@dataclass(eq=False, repr=False)
class PagingHint(betterproto.Message):
    older: "GetLogsRequest" = betterproto.message_field(1)
    newer: "GetLogsRequest" = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class LogMessage(betterproto.Message):
    verbatim: "VerbatimLogMessage" = betterproto.message_field(1, group="log_type")
    json: "ParsedJsonLogMessage" = betterproto.message_field(2, group="log_type")
    pps_log_message: "_pps__.LogMessage" = betterproto.message_field(
        3, group="log_type"
    )


@dataclass(eq=False, repr=False)
class VerbatimLogMessage(betterproto.Message):
    line: bytes = betterproto.bytes_field(1)
    timestamp: datetime = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class ParsedJsonLogMessage(betterproto.Message):
    verbatim: "VerbatimLogMessage" = betterproto.message_field(1)
    """The verbatim line from Loki"""

    object: "betterproto_lib_google_protobuf.Struct" = betterproto.message_field(2)
    """A raw JSON parse of the entire line"""

    native_timestamp: datetime = betterproto.message_field(3)
    """If a parseable timestamp was found in `fields`"""

    pps_log_message: "_pps__.LogMessage" = betterproto.message_field(4)
    """For code that wants to filter on pipeline/job/etc"""


class ApiStub:

    def __init__(self, channel: "grpc.Channel"):
        self.__rpc_get_logs = channel.unary_stream(
            "/logs.API/GetLogs",
            request_serializer=GetLogsRequest.SerializeToString,
            response_deserializer=GetLogsResponse.FromString,
        )

    def get_logs(
        self,
        *,
        query: "LogQuery" = None,
        filter: "LogFilter" = None,
        tail: bool = False,
        want_paging_hint: bool = False,
        log_format: "LogFormat" = None
    ) -> Iterator["GetLogsResponse"]:

        request = GetLogsRequest()
        if query is not None:
            request.query = query
        if filter is not None:
            request.filter = filter
        request.tail = tail
        request.want_paging_hint = want_paging_hint
        request.log_format = log_format

        for response in self.__rpc_get_logs(request):
            yield response


class ApiBase:

    def get_logs(
        self,
        query: "LogQuery",
        filter: "LogFilter",
        tail: bool,
        want_paging_hint: bool,
        log_format: "LogFormat",
        context: "grpc.ServicerContext",
    ) -> Iterator["GetLogsResponse"]:
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    __proto_path__ = "logs.API"

    @property
    def __rpc_methods__(self):
        return {
            "GetLogs": grpc.unary_stream_rpc_method_handler(
                self.get_logs,
                request_deserializer=GetLogsRequest.FromString,
                response_serializer=GetLogsRequest.SerializeToString,
            ),
        }