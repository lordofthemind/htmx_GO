package configs

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

var (
	Port                int
	UseTLS              bool
	AllowedCredentials  bool
	TokenType           string
	StaticPath          string
	MongoDBUrl          string
	TlsKeyFile          string
	Environment         string
	TlsCertFile         string
	TemplatePath        string
	TokenSymmetricKey   string
	AllowedOrigins      []string
	AllowedMethods      []string
	AllowedHeaders      []string
	ExposedHeaders      []string
	TokenAccessDuration time.Duration
)

// InitializeServerConfig initializes the server configuration using Viper.
// It returns an error if any issues occur during the configuration setup.
func InitializeServerConfig(configFile string) error {
	viper.SetConfigFile(configFile)

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Load main configuration values
	Port = viper.GetInt("server.port")
	Environment = viper.GetString("application.config")
	TemplatePath = viper.GetString("application.template_path")
	StaticPath = viper.GetString("application.static_path")
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
		return fmt.Errorf("invalid duration for token.access_duration: %w", err)
	}

	log.Printf("Token access duration set to: %s", TokenAccessDuration.String())
	log.Printf("Server is being initiated with %s environment", Environment)
	return nil
}

// loadEnvironmentConfig loads the server configuration for the given environment.
func loadEnvironmentConfig(env string) {
	TlsKeyFile = viper.GetString(env + ".key_file")
	TlsCertFile = viper.GetString(env + ".cert_file")
	MongoDBUrl = viper.GetString(env + ".mongoDB_url")
	AllowedOrigins = viper.GetStringSlice(env + ".cors.allowed_origins")
	AllowedMethods = viper.GetStringSlice(env + ".cors.allowed_methods")
	AllowedHeaders = viper.GetStringSlice(env + ".cors.allowed_headers")
	ExposedHeaders = viper.GetStringSlice(env + ".cors.exposed_headers")
	AllowedCredentials = viper.GetBool(env + ".cors.allow_credentials")
}
