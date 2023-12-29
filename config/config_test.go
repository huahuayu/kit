package config

import (
	"os"
	"testing"
)

type sampleConfig struct {
	Database struct {
		Host string `yaml:"host" json:"host" toml:"host" env:"DATABASE_HOST"`
		Port int    `yaml:"port" json:"port" toml:"port" env:"DATABASE_PORT"`
		User string `yaml:"user" json:"user" toml:"user" env:"DATABASE_USER"`
	} `yaml:"database" json:"database" toml:"database"`
	Logging struct {
		Level string `yaml:"level" json:"level" toml:"level" env:"LOGGING_LEVEL"`
	} `yaml:"logging" json:"logging" toml:"logging"`
	FeatureFlags struct {
		BetaFeatures bool `yaml:"betaFeatures" json:"betaFeatures" toml:"betaFeatures" env:"FEATUREFLAGS_BETAFEATURES"`
	} `yaml:"featureFlags" json:"featureFlags" toml:"featureFlags"`
	ApiVersion []string          `yaml:"apiVersion" json:"apiVersion" toml:"apiVersion" env:"APIVERSION"`
	Mapping    map[string]string `yaml:"mapping" json:"mapping" toml:"mapping" env:"MAPPING"`
}

func TestLoadYAML(t *testing.T) {
	loader := NewConfigLoader[sampleConfig]()
	cfg, err := loader.Load("testdata/sample.yml", YAML)
	if err != nil {
		t.Fatalf("Failed to load YAML: %v", err)
	}
	validateConfig(t, cfg)
}

func TestLoadTOML(t *testing.T) {
	loader := NewConfigLoader[sampleConfig]()
	cfg, err := loader.Load("testdata/sample.toml", TOML)
	if err != nil {
		t.Fatalf("Failed to load TOML: %v", err)
	}
	validateConfig(t, cfg)
}

func TestLoadJSON(t *testing.T) {
	loader := NewConfigLoader[sampleConfig]()
	cfg, err := loader.Load("testdata/sample.json", JSON)
	if err != nil {
		t.Fatalf("Failed to load JSON: %v", err)
	}
	validateConfig(t, cfg)
}

func TestLoadDotEnv(t *testing.T) {
	loader := NewConfigLoader[sampleConfig]()
	cfg, err := loader.Load("testdata/sample.env", DOTENV)
	if err != nil {
		t.Fatalf("Failed to load .env: %v", err)
	}
	validateConfig(t, cfg)
}

func validateConfig(t *testing.T, cfg sampleConfig) {
	// Validate database configuration
	if cfg.Database.Host != "dbserver" || cfg.Database.Port != 5432 || cfg.Database.User != "admin" {
		t.Errorf("Database configuration did not load correctly. Got: %+v", cfg.Database)
	}

	// Validate logging level
	if cfg.Logging.Level != "debug" {
		t.Errorf("Logging level did not load correctly. Got: %s", cfg.Logging.Level)
	}

	// Validate feature flags
	if !cfg.FeatureFlags.BetaFeatures {
		t.Errorf("Feature flags did not load correctly. BetaFeatures should be true.")
	}

	// Validate API versions
	expectedApiVersions := []string{"v1", "v2", "v3"}
	for i, v := range cfg.ApiVersion {
		if v != expectedApiVersions[i] {
			t.Errorf("API version mismatch. Expected %s, got %s", expectedApiVersions[i], v)
		}
	}

	// Validate mappings
	expectedMapping := map[string]string{"foo": "bar", "baz": "qux"}
	for k, v := range expectedMapping {
		if cfg.Mapping[k] != v {
			t.Errorf("Mapping mismatch. Expected %s for key %s, got %s", v, k, cfg.Mapping[k])
		}
	}
}

func TestLoadYAMLWithEnvOverride(t *testing.T) {
	// Set environment variables to override config values
	os.Setenv("DATABASE_HOST", "env_dbserver")
	os.Setenv("DATABASE_PORT", "5433") // Changed for test
	os.Setenv("DATABASE_USER", "env_admin")
	os.Setenv("LOGGING_LEVEL", "info")
	os.Setenv("FEATUREFLAGS_BETAFEATURES", "false")
	os.Setenv("APIVERSION", "v4,v5")                // Comma-separated for slice
	os.Setenv("MAPPING", "foo:env_bar,baz:env_qux") // Comma-separated for map

	defer func() {
		// Clean up environment variables after test
		os.Unsetenv("DATABASE_HOST")
		os.Unsetenv("DATABASE_PORT")
		os.Unsetenv("DATABASE_USER")
		os.Unsetenv("LOGGING_LEVEL")
		os.Unsetenv("FEATUREFLAGS_BETAFEATURES")
		os.Unsetenv("APIVERSION")
		os.Unsetenv("MAPPING")
	}()

	loader := NewConfigLoader[sampleConfig]()
	cfg, err := loader.Load("testdata/sample.yml", YAML)
	if err != nil {
		t.Fatalf("Failed to load YAML: %v", err)
	}

	// Validate overridden database configuration
	if cfg.Database.Host != "env_dbserver" || cfg.Database.Port != 5433 || cfg.Database.User != "env_admin" {
		t.Errorf("Database configuration did not override correctly. Got: %+v", cfg.Database)
	}

	// Validate overridden logging level
	if cfg.Logging.Level != "info" {
		t.Errorf("Logging level did not override correctly. Got: %s", cfg.Logging.Level)
	}

	// Validate overridden feature flags
	if cfg.FeatureFlags.BetaFeatures {
		t.Errorf("Feature flags did not override correctly. BetaFeatures should be false.")
	}

	// Validate overridden API versions
	expectedApiVersions := []string{"v4", "v5"}
	for i, v := range cfg.ApiVersion {
		if v != expectedApiVersions[i] {
			t.Errorf("API version override mismatch. Expected %s, got %s", expectedApiVersions[i], v)
		}
	}

	// Validate overridden mappings
	expectedMapping := map[string]string{"foo": "env_bar", "baz": "env_qux"}
	for k, v := range expectedMapping {
		if cfg.Mapping[k] != v {
			t.Errorf("Mapping override mismatch. Expected %s for key %s, got %s", v, k, cfg.Mapping[k])
		}
	}
}
