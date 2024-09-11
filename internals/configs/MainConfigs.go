package configs

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var (
	Port                int
	UseTLS              bool
	TlsKeyFile          string
	TlsCertFile         string
	Environment         string
	AllowedOrigins      []string
	AllowedMethods      []string
	AllowedHeaders      []string
	ExposedHeaders      []string
	PostgreGormDBUrl    string
	AllowedCredentials  bool
	TokenType           string
	TokenSymmetricKey   string
	TokenAccessDuration time.Duration
)

// InitializeServerConfig initializes the server configuration using Viper.
func InitializeServerConfig(configFile string) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Load main configuration values
	Port = viper.GetInt("server.port")
	Environment = viper.GetString("application.config")
	UseTLS = viper.GetBool("tls.use_tls")
	TokenSymmetricKey = viper.GetString("token.symmetric_key")
	TokenType = viper.GetString("token.type")

	// Load environment-specific configurations
	loadEnvironmentConfig(Environment)

	// Parse the token access duration
	durationStr := viper.GetString("token.access_duration")
	var err error
	TokenAccessDuration, err = time.ParseDuration(durationStr)
	if err != nil {
		log.Fatalf("Invalid duration for token.access_duration: %v", err)
	}

	log.Printf("Token access duration set to: %s", TokenAccessDuration.String())
	log.Printf("Server is being initiated with %s environment", Environment)
}

// loadEnvironmentConfig loads the server configuration for the given environment.
func loadEnvironmentConfig(env string) {
	TlsKeyFile = viper.GetString(env + ".key_file")
	TlsCertFile = viper.GetString(env + ".cert_file")
	PostgreGormDBUrl = viper.GetString(env + ".postgres_gorm_url")
	AllowedOrigins = viper.GetStringSlice(env + ".cors.allowed_origins")
	AllowedMethods = viper.GetStringSlice(env + ".cors.allowed_methods")
	AllowedHeaders = viper.GetStringSlice(env + ".cors.allowed_headers")
	ExposedHeaders = viper.GetStringSlice(env + ".cors.exposed_headers")
	AllowedCredentials = viper.GetBool(env + ".cors.allow_credentials")
}
