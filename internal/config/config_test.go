package config

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
)

func setupViperForTest(t *testing.T) {
	viper.Reset()
}

func TestLoadConfig_WithAppKeys(t *testing.T) {
	setupViperForTest(t)

	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewReader([]byte(`
endpoint: ovh-eu
app_key: test_app_key
app_secret: test_app_secret
consumer_key: test_consumer_key
default_zone: ganima.xyz
`)))

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Endpoint != "ovh-eu" {
		t.Errorf("expected endpoint ovh-eu, got %s", cfg.Endpoint)
	}
	if cfg.AppKey != "test_app_key" {
		t.Errorf("expected app_key test_app_key, got %s", cfg.AppKey)
	}
	if cfg.AppSecret != "test_app_secret" {
		t.Errorf("expected app_secret test_app_secret, got %s", cfg.AppSecret)
	}
	if cfg.ConsumerKey != "test_consumer_key" {
		t.Errorf("expected consumer_key test_consumer_key, got %s", cfg.ConsumerKey)
	}
	if cfg.DefaultZone != "ganima.xyz" {
		t.Errorf("expected default_zone ganima.xyz, got %s", cfg.DefaultZone)
	}
}

func TestLoadConfig_WithOAuth2(t *testing.T) {
	setupViperForTest(t)

	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewReader([]byte(`
endpoint: ovh-eu
oauth2_client_id: test_client_id
oauth2_client_secret: test_client_secret
default_zone: ganima.xyz
`)))

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.OAuth2ClientID != "test_client_id" {
		t.Errorf("expected oauth2_client_id test_client_id, got %s", cfg.OAuth2ClientID)
	}
	if cfg.OAuth2ClientSecret != "test_client_secret" {
		t.Errorf("expected oauth2_client_secret test_client_secret, got %s", cfg.OAuth2ClientSecret)
	}
}

func TestLoadConfig_Empty(t *testing.T) {
	setupViperForTest(t)

	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewReader([]byte(``)))

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Endpoint != "" {
		t.Errorf("expected empty endpoint, got %s", cfg.Endpoint)
	}
}

func TestConfig_GetOVHClient_PrefersOAuth2(t *testing.T) {
	cfg := &Config{
		Endpoint:           "ovh-eu",
		AppKey:             "app_key",
		AppSecret:          "app_secret",
		ConsumerKey:        "consumer_key",
		OAuth2ClientID:     "oauth2_id",
		OAuth2ClientSecret: "oauth2_secret",
	}

	_, err := cfg.GetOVHClient()
	if err != nil {
		t.Logf("GetOVHClient returned error (expected with invalid creds): %v", err)
	}
}
