package config

import (
	"github.com/ovh/go-ovh/ovh"
	"github.com/spf13/viper"
)

type Config struct {
	Endpoint    string `mapstructure:"endpoint"`
	DefaultZone string `mapstructure:"default_zone"`

	AppKey      string `mapstructure:"app_key"`
	AppSecret   string `mapstructure:"app_secret"`
	ConsumerKey string `mapstructure:"consumer_key"`

	OAuth2ClientID     string `mapstructure:"oauth2_client_id"`
	OAuth2ClientSecret string `mapstructure:"oauth2_client_secret"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) GetOVHClient() (*ovh.Client, error) {
	if c.OAuth2ClientID != "" && c.OAuth2ClientSecret != "" {
		return c.getOAuth2Client()
	}
	return c.getLegacyClient()
}

func (c *Config) getLegacyClient() (*ovh.Client, error) {
	return ovh.NewClient(
		c.Endpoint,
		c.AppKey,
		c.AppSecret,
		c.ConsumerKey,
	)
}

func (c *Config) getOAuth2Client() (*ovh.Client, error) {
	return ovh.NewEndpointClient(c.Endpoint)
}
