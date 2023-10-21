package utils

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

// EnvDefaults : EnvDefaults Structure to hold Env default Variable values
type EnvDefaults struct {
	// ServerHost api server host
	ServerHost string

	// ServerPort api server port
	ServerPort string

	// HostURL api host URL
	HostURL string

	// SessionSecret secret key for session cookie
	SessionSecret string

	// JWTTimeoutMinutes JWT timeout
	JWTTimeoutMinutes string

	// JWTAuthSecret JWT timeout
	JWTAuthSecret *rsa.PrivateKey

	// JWTSigningSecret secret word to sign jwt tokens
	JWTSigningSecret string

	// AuthenticityTokenTTLMinutes authenticity token time to live in minutes
	AuthenticityTokenTTLMinutes string

	// ScryptSecret secret key for scrypt hash of passwords
	ScryptSecret string

	// CronConfig cron initial configuration
	CronConfig string

	// FluentHost fluentd host
	FluentHost string

	// FluentPort fluentd port
	FluentPort string

	// ElasticURL elastic search url
	ElasticURL string

	// ElasticUsername elastic search server username
	ElasticUsername string

	// ElasticPassword elastic search server password
	ElasticPassword string

	// ElasticBearer elastic search server token
	ElasticBearer string

	// CollectionPrefix database collection prefix to all collections
	CollectionPrefix string

	OpaClientUrl    string
	OpaAppName      string
	SpiceDBURL      string
	SpiceDBToken    string
	SpiceDBInsecure string
}

// GetEnv return values of the environment
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getJWTTokenKey(jwtSecretBase64 string) *rsa.PrivateKey {
	jwtSecretPem, err := base64.StdEncoding.DecodeString(jwtSecretBase64)
	if err != nil {
		panic(fmt.Errorf("no valid JWT secret (JWT_AUTH_SECRET) in base64 format: %s", err.Error()))
	}
	jwtSecret, err := jwt.ParseRSAPrivateKeyFromPEM(jwtSecretPem)
	if err != nil {
		panic(fmt.Errorf("no valid JWT secret (JWT_AUTH_SECRET): %s", err.Error()))
	}

	return jwtSecret
}

// Env : Env variable to hold Env default Variable values
var Env EnvDefaults

// SetupEnvDefaults : Initialize EnvDefaults
func SetupEnvDefaults() {
	Env = EnvDefaults{
		// Server configurations
		ServerHost:    GetEnv("SERVER_HOST", ""),
		ServerPort:    GetEnv("SERVER_PORT", ""),
		HostURL:       GetEnv("HOST_URL", ""),
		SessionSecret: GetEnv("SESSION_SECRET", ""),

		// JWT configuration
		JWTTimeoutMinutes: GetEnv("JWT_TIMEOUT_MINUTES", "1440"),
		JWTAuthSecret:     getJWTTokenKey(GetEnv("JWT_AUTH_SECRET", "")),
		ScryptSecret:      GetEnv("SCRYPT_SECRET", ""),

		CronConfig: GetEnv("CRON_CONFIG", "@every 30s"),

		JWTSigningSecret:            GetEnv("JWT_SIGNING_SECRET", ""),
		AuthenticityTokenTTLMinutes: GetEnv("AUTHENTICITY_TOKEN_TTL_MINUTES", "5"),

		FluentHost: GetEnv("FLUENT_HOST", ""),
		FluentPort: GetEnv("FLUENT_PORT", ""),

		ElasticURL:      GetEnv("ELASTIC_URL", ""),
		ElasticUsername: GetEnv("ELASTIC_USERNAME", ""),
		ElasticPassword: GetEnv("ELASTIC_PASSWORD", ""),

		CollectionPrefix: GetEnv("MONGO_COLLECTION_PREFIX", "papelito_"),

		OpaClientUrl: GetEnv("OPA_CLIENT_URL", ""),
		OpaAppName:   GetEnv("OPA_APP_NAME", ""),

		SpiceDBURL:      GetEnv("SPICEDB_GRPC_URL", ""),
		SpiceDBToken:    GetEnv("SPICEDB_TOKEN", ""),
		SpiceDBInsecure: GetEnv("SPICEDB_INSECURE", ""),
	}
}
