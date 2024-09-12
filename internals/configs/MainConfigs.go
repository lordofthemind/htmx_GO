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
	UseCORS             bool
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
	UseCORS = viper.GetBool("server.use_cors")
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
	TlsKeyFile = viper.GetString(fmt.Sprintf("%s.key_file", env))
	TlsCertFile = viper.GetString(fmt.Sprintf("%s.cert_file", env))
	MongoDBUrl = viper.GetString(fmt.Sprintf("%s.mongoDB_url", env))
	AllowedOrigins = viper.GetStringSlice(fmt.Sprintf("%s.cors.allowed_origins", env))
	AllowedMethods = viper.GetStringSlice(fmt.Sprintf("%s.cors.allowed_methods", env))
	AllowedHeaders = viper.GetStringSlice(fmt.Sprintf("%s.cors.allowed_headers", env))
	ExposedHeaders = viper.GetStringSlice(fmt.Sprintf("%s.cors.exposed_headers", env))
	AllowedCredentials = viper.GetBool(fmt.Sprintf("%s.cors.allow_credentials", env))
}
