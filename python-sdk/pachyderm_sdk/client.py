import contextlib
import os
import json
from base64 import b64decode
from dataclasses import dataclass
from pathlib import Path
from tempfile import NamedTemporaryFile
from typing import Dict, Optional, Union
from urllib.parse import urlparse

import grpc

from .api.admin.extension import ApiStub as _AdminStub
from .api.auth import ApiStub as _AuthStub
from .api.debug import DebugStub as _DebugStub
from .api.enterprise import ApiStub as _EnterpriseStub
from .api.identity import ApiStub as _IdentityStub
from .api.license import ApiStub as _LicenseStub
from .api.pfs.extension import ApiStub as _PfsStub
from .api.pps.extension import ApiStub as _PpsStub
from .api.transaction.extension import ApiStub as _TransactionStub
from .api.version import ApiStub as _VersionStub, Version
from .constants import (
    AUTH_TOKEN_ENV,
    GRPC_CHANNEL_OPTIONS,
    OIDC_TOKEN_ENV,
    PACHD_SERVICE_HOST_ENV,
    PACHD_SERVICE_PORT_ENV,
)
from .errors import AuthServiceNotActivated, BadClusterDeploymentID, ConfigError
from .interceptor import MetadataClientInterceptor, MetadataType

__all__ = ("Client", )


class Client:
    """The :class:`.Client` class that users will primarily interact with.
    Initialize an instance with ``python_pachyderm.Client()``.

    To see documentation on the methods :class:`.Client` can call, refer to the
    `mixins` module.
    """

    # Class variables for checking config
    env_config = "PACH_CONFIG"
    spout_config = "/pachctl/config.json"
    local_config = f"{Path.home()}/.pachyderm/config.json"

    def __init__(
        self,
        host: str = 'localhost',
        port: int = 30650,
        auth_token: Optional[str] = None,
        root_certs: Optional[bytes] = None,
        transaction_id: str = None,
        tls: bool = False,
    ):
        """
        Creates a Pachyderm client. If both files don't exist, a client
        with default settings is created.

        Parameters
        ----------
        host : str, optional
            The pachd host. Default is 'localhost', which is used with
            ``pachctl port-forward``.
        port : int, optional
            The port to connect to. Default is 30650.
        auth_token : str, optional
            The authentication token. Used if authentication is enabled on the
            cluster.
        root_certs : bytes, optional
            The PEM-encoded root certificates as byte string.
        transaction_id : str, optional
            The ID of the transaction to run operations on.
        tls : bool
            Whether TLS should be used. If `root_certs` are specified, they are
            used. Otherwise, we use the certs provided by certifi.
        """
        if auth_token is None:
            auth_token = os.environ.get(AUTH_TOKEN_ENV)

        tls = tls or (root_certs is not None)
        if tls and root_certs is None:
            # load default certs if none are specified
            import certifi

            with open(certifi.where(), "rb") as f:
                root_certs = f.read()

        self.address = "{}:{}".format(host, port)
        self.root_certs = root_certs
        channel = _create_channel(
            self.address, self.root_certs, options=GRPC_CHANNEL_OPTIONS
        )

        self._auth_token = auth_token
        self._transaction_id = transaction_id
        self._metadata = self._build_metadata()
        self._channel = _apply_metadata_interceptor(channel, self._metadata)

        # See implementation for api layout.
        self._init_api()

        if not auth_token and (oidc_token := os.environ.get(OIDC_TOKEN_ENV)):
            self.auth_token = self.auth.authenticate(id_token=oidc_token)

    def _init_api(self):
        self.admin = _AdminStub(self._channel)
        self.auth = _AuthStub(self._channel)
        self.debug = _DebugStub(self._channel)
        self.enterprise = _EnterpriseStub(self._channel)
        self.identity = _IdentityStub(self._channel)
        self.license = _LicenseStub(self._channel)
        self.pfs = _PfsStub(self._channel)
        self.pps = _PpsStub(self._channel)
        self.transaction = _TransactionStub(
            self._channel,
            get_transaction_id=lambda: self.transaction_id,
            set_transaction_id=lambda value: setattr(self, "transaction_id", value),
        )
        self._version_api = _VersionStub(self._channel)

    @classmethod
    def new_in_cluster(
        cls,
        auth_token: Optional[str] = None,
        transaction_id: Optional[str] = None
    ) -> "Client":
        """Creates a Pachyderm client that operates within a Pachyderm cluster.

        Parameters
        ----------
        auth_token : str, optional
            The authentication token. Used if authentication is enabled on the
            cluster.
        transaction_id : str, optional
            The ID of the transaction to run operations on.

        Returns
        -------
        Client
            A python_pachyderm client instance.
        """

        return cls(
            host=os.environ[PACHD_SERVICE_HOST_ENV],
            port=int(os.environ[PACHD_SERVICE_PORT_ENV]),
            auth_token=auth_token,
            transaction_id=transaction_id,
        )

    @classmethod
    def from_pachd_address(
        cls,
        pachd_address: str,
        auth_token: str = None,
        root_certs: bytes = None,
        transaction_id: str = None,
    ) -> "Client":
        """Creates a Pachyderm client from a given pachd address.

        Parameters
        ----------
        pachd_address : str
            The address of pachd server
        auth_token : str, optional
            The authentication token. Used if authentication is enabled on the
            cluster.
        root_certs : bytes, optional
            The PEM-encoded root certificates as byte string. If unspecified,
            this will load default certs from certifi.
        transaction_id : str, optional
            The ID of the transaction to run operations on.

        Returns
        -------
        Client
            A python_pachyderm client instance.
        """
        if "://" not in pachd_address:
            pachd_address = "grpc://{}".format(pachd_address)

        u = urlparse(pachd_address)

        if u.scheme not in ("grpc", "http", "grpcs", "https"):
            raise ValueError("unrecognized pachd address scheme: {}".format(u.scheme))
        if u.path or u.params or u.query or u.fragment or u.username or u.password:
            raise ValueError("invalid pachd address")

        return cls(
            host=u.hostname,
            port=u.port,
            auth_token=auth_token,
            root_certs=root_certs,
            transaction_id=transaction_id,
            tls=u.scheme == "grpcs" or u.scheme == "https",
        )

    @classmethod
    def from_config(cls, config_file: Union[Path, str]) -> "Client":
        """Creates a Pachyderm client from a config file.

        Parameters
        ----------
        config_file : Union[Path, str]
            The path to a config json file.

        Returns
        -------
        Client
            A properly configured Client.
        """
        config = _ConfigFile(config_file)
        active_context = config.active_context
        client = cls.from_pachd_address(
            active_context.active_pachd_address,
            auth_token=active_context.session_token,
            root_certs=active_context.server_cas_decoded,
            transaction_id=active_context.active_transaction,
        )

        # Verify the deployment ID of the active context with the cluster.
        expected_deployment_id = active_context.cluster_deployment_id
        if expected_deployment_id:
            cluster_info = client.admin.inspect_cluster()
            if cluster_info.deployment_id != expected_deployment_id:
                raise BadClusterDeploymentID(
                    expected_deployment_id, cluster_info.deployment_id
                )

        return client

    @property
    def auth_token(self):
        return self._auth_token

    @auth_token.setter
    def auth_token(self, value):
        self._auth_token = value
        self._metadata = self._build_metadata()
        self._channel = _apply_metadata_interceptor(
            channel=_create_channel(
                self.address, self.root_certs, options=GRPC_CHANNEL_OPTIONS
            ),
            metadata=self._metadata,
        )
        self._init_api()

    @property
    def transaction_id(self):
        return self._transaction_id

    @transaction_id.setter
    def transaction_id(self, value):
        self._transaction_id = value
        self._metadata = self._build_metadata()
        self._channel = _apply_metadata_interceptor(
            channel=_create_channel(
                self.address, self.root_certs, options=GRPC_CHANNEL_OPTIONS
            ),
            metadata=self._metadata,
        )
        self._init_api()

    def _build_metadata(self):
        metadata = []
        if self._auth_token is not None:
            metadata.append(("authn-token", self._auth_token))
        if self._transaction_id is not None:
            metadata.append(("pach-transaction", self._transaction_id))
        return metadata

    def delete_all(self) -> None:
        """Delete all repos, commits, files, pipelines, and jobs.
        This resets the cluster to its initial state.
        """
        # Try removing all identities if auth is activated.
        with contextlib.suppress(AuthServiceNotActivated):
            self.identity.delete_all()

        # Try deactivating auth if activated.
        with contextlib.suppress(AuthServiceNotActivated):
            self.auth.deactivate()

        # Try removing all licenses if auth is activated.
        with contextlib.suppress(AuthServiceNotActivated):
            self.license.delete_all()

        self.pps.delete_all()
        self.pfs.delete_all()
        self.transaction.delete_all()

    def get_version(self) -> Version:
        return self._version_api.get_version()


def _apply_metadata_interceptor(
    channel: grpc.Channel, metadata: MetadataType
) -> grpc.Channel:
    metadata_interceptor = MetadataClientInterceptor(metadata)
    return grpc.intercept_channel(channel, metadata_interceptor)


def _create_channel(
    address: str,
    root_certs: Optional[bytes],
    options: MetadataType,
) -> grpc.Channel:
    if root_certs is not None:
        ssl = grpc.ssl_channel_credentials(root_certificates=root_certs)
        return grpc.secure_channel(address, ssl, options=options)
    return grpc.insecure_channel(address, options=options)


class _ConfigFile:

    def __init__(self, config_file: Union[Path, str]):
        config_file = Path(os.path.expanduser(config_file)).resolve()
        self._config_file_data = json.loads(config_file.read_bytes())

    @classmethod
    def from_bytes(cls, config_file_data: bytes):
        with NamedTemporaryFile() as temp_config_file:
            temp_config_file.write(config_file_data)
            return cls(temp_config_file.name)

    @property
    def user_id(self) -> str:
        return self._config_file_data["user_id"]

    @property
    def active_context(self) -> "_Context":
        active_context_name = self._config_file_data["v2"]["active_context"]
        contexts = self._config_file_data["v2"]["contexts"]
        if active_context_name not in contexts:
            raise ConfigError(f"active context not found: {active_context_name}")
        return _Context(**contexts[active_context_name])

    @property
    def active_enterprise_context(self) -> "_Context":
        context_name = self._config_file_data["v2"].get("active_enterprise_context")
        if context_name is None:
            raise ConfigError("active enterprise context is not specified")
        contexts = self._config_file_data["v2"]["contexts"]
        if context_name not in contexts:
            raise ConfigError(f"active enterprise context not found: {context_name}")
        return _Context(**contexts[context_name])


@dataclass
class _Context:
    source: Optional[int] = None
    """An integer that specifies where the config came from. 
    This parameter is for internal use only and should not be modified."""

    pachd_address: Optional[str] = None
    """A host:port specification for connecting to pachd."""

    server_cas: Optional[str] = None
    """Trusted root certificates for the cluster, formatted as a 
    base64-encoded PEM. This is only set when TLS is enabled."""

    session_token: Optional[str] = None
    """A secret token identifying the current user within their pachyderm
    cluster. This is included in all RPCs and used to determine if a user's
    actions are authorized. This is only set when auth is enabled."""

    active_transaction: Optional[str] = None
    """The currently active transaction for batching together commands."""

    cluster_name: Optional[str] = None
    """The name of the underlying Kubernetes cluster."""

    auth_info: Optional[str] = None
    """The name of the underlying Kubernetes cluster’s auth credentials"""

    namespace: Optional[str] = None
    """The underlying Kubernetes cluster’s namespace"""

    cluster_deployment_id: Optional[str] = None
    """The pachyderm cluster deployment ID that is used to ensure the
    operations run on the expected cluster."""

    project: Optional[str] = None

    enterprise_server: bool = False
    """Whether the context represents an enterprise server."""

    port_forwarders: Dict[str, int] = None
    """A mapping of service name -> local port."""

    @property
    def active_pachd_address(self) -> str:
        """This pachd factors in port-forwarding. """
        if self.pachd_address is None:
            port = 30650
            if self.port_forwarders:
                port = self.port_forwarders.get('pachd', 30650)
            return f"grpc://localhost:{port}"
        return self.pachd_address

    @property
    def server_cas_decoded(self) -> Optional[bytes]:
        """The base64 decoded root certificates in PEM format, if they exist."""
        if self.server_cas:
            return b64decode(bytes(self.server_cas, "utf-8"))