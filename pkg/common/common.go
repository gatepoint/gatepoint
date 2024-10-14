package common

import "time"

// Argo CD application related constants
const (

	// GatepointAdminUsername is the username of the 'admin' user
	GatepointAdminUsername = "admin"
	// GatepointUserAgentName is the default user-agent name used by the gRPC API client library and grpc-gateway
	GatepointUserAgentName = "gatepoint-client"
	// GatepointSSAManager is the default Gatepoint manager name used by server-side apply syncs
	GatepointSSAManager = "gatepoint-controller"
	// AuthCookieName is the HTTP cookie name where we store our auth token
	AuthCookieName = "gatepoint.token"
	// StateCookieName is the HTTP cookie name that holds temporary nonce tokens for CSRF protection
	StateCookieName = "gatepoint.oauthstate"
	// StateCookieMaxAge is the maximum age of the oauth state cookie
	StateCookieMaxAge = time.Minute * 5

	// ChangePasswordSSOTokenMaxAge is the max token age for password change operation
	ChangePasswordSSOTokenMaxAge = time.Minute * 5
	// GithubAppCredsExpirationDuration is the default time used to cache the GitHub app credentials
	GithubAppCredsExpirationDuration = time.Minute * 60

	// PasswordPatten is the default password patten
	PasswordPatten = `^.{8,32}$`

	// LegacyShardingAlgorithm is the default value for Sharding Algorithm it uses an `uid` based distribution (non-uniform)
	LegacyShardingAlgorithm = "legacy"
	// RoundRobinShardingAlgorithm is a flag value that can be opted for Sharding Algorithm it uses an equal distribution across all shards
	RoundRobinShardingAlgorithm = "round-robin"
	// AppControllerHeartbeatUpdateRetryCount is the retry count for updating the Shard Mapping to the Shard Mapping ConfigMap used by Application Controller
	AppControllerHeartbeatUpdateRetryCount = 3

	// ConsistentHashingWithBoundedLoadsAlgorithm uses an algorithm that tries to use an equal distribution across
	// all shards but is optimised to handle sharding and/or cluster addition or removal. In case of sharding or
	// cluster changes, this algorithm minimises the changes between shard and clusters assignments.
	ConsistentHashingWithBoundedLoadsAlgorithm = "consistent-hashing"

	DefaultShardingAlgorithm = LegacyShardingAlgorithm
)

// Kubernetes ConfigMap and Secret resource names which hold Argo CD settings
const (
	GatepointConfigMapName              = "gatepoint-cm"
	GatepointSecretName                 = "gatepoint-secret"
	GatepointNotificationsConfigMapName = "gatepoint-notifications-cm"
	GatepointNotificationsSecretName    = "gatepoint-notifications-secret"
	GatepointRBACConfigMapName          = "gatepoint-rbac-cm"
	// GatepointKnownHostsConfigMapName contains SSH known hosts data for connecting repositories. Will get mounted as volume to pods
	GatepointKnownHostsConfigMapName = "gatepoint-ssh-known-hosts-cm"
	// GatepointTLSCertsConfigMapName contains TLS certificate data for connecting repositories. Will get mounted as volume to pods
	GatepointTLSCertsConfigMapName = "gatepoint-tls-certs-cm"
	GatepointGPGKeysConfigMapName  = "gatepoint-gpg-keys-cm"
	// GatepointAppControllerShardConfigMapName contains the application controller to shard mapping
	GatepointAppControllerShardConfigMapName = "gatepoint-app-controller-shard-cm"
	GatepointCmdParamsConfigMapName          = "gatepoint-cmd-params-cm"
)
