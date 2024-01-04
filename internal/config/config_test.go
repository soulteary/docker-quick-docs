package config_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/soulteary/docker-quick-docs/internal/config"
	"github.com/stretchr/testify/assert"
)

// TestReadConfigFileSuccess tests if ReadConfigFile can read the default config file.
func TestReadConfigFileSuccess(t *testing.T) {
	const configFileContent = `[
		{
			"from": "https://ecomfe.github.io/san/",
			"to": "/san/",
			"type": "html",
			"dir": "/san/"
		}
	]`

	// Create a temporary file simulating the config.json
	tmpfile, err := os.CreateTemp(".", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	// Write data to the tmp config file
	if _, err := tmpfile.Write([]byte(configFileContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test ReadConfigFile
	data, err := config.ReadConfigFile(tmpfile.Name())

	assert.NoError(t, err)
	assert.Equal(t, configFileContent, string(data))
	config.PostRules = []config.PostRule{}
}

// TestReadConfigFileEnvVar tests if ReadConfigFile correctly reads from CONFIG environment variable.
func TestReadConfigFileEnvVar(t *testing.T) {
	const configFileContent = `{"from":"env_from","to":"env_to"}`
	// Create a temporary file simulating the config.json
	tmpfile, err := os.CreateTemp(".", "env_config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write data to the tmp config file
	if _, err := tmpfile.Write([]byte(configFileContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Set the environment variable to point to our temp file
	os.Setenv("CONFIG", tmpfile.Name())
	defer os.Unsetenv("CONFIG")

	data, err := config.ReadConfigFile(config.DOCS_DEFAULT_CONFIG)
	assert.NoError(t, err)
	assert.Equal(t, configFileContent, string(data))
	config.PostRules = []config.PostRule{}
}

// TestReadConfigFileNoFile tests ReadConfigFile's response when no file is present.
func TestReadConfigFileNoFile(t *testing.T) {
	os.Unsetenv("CONFIG")
	data, err := config.ReadConfigFile("not-exist.json")
	assert.Error(t, err)
	assert.Equal(t, []byte(""), data)
}

// TestGetConfigSuccess tests if GetConfig successfully parses the configuration.
func TestGetConfigSuccess(t *testing.T) {
	const configFileContent = `[
		{
			"from": "/source",
			"to": "/destination"
		}
	]`

	// Create a temporary file simulating the config.json
	tmpfile, err := os.CreateTemp(".", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	// Write data to the tmp config file
	if _, err := tmpfile.Write([]byte(configFileContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	os.Setenv("CONFIG", tmpfile.Name())
	defer os.Unsetenv("CONFIG")

	config.GetConfig()

	var expected []config.PostRule
	err = json.Unmarshal([]byte(configFileContent), &expected)
	expected[0].Type = "text/html"
	expected[0].Dir = "*"
	assert.NoError(t, err)
	assert.Equal(t, expected, config.PostRules)
	config.PostRules = []config.PostRule{}
}

// TestGetConfigParseError tests if GetConfig handles JSON unmarshalling errors.
func TestGetConfigParseError(t *testing.T) {
	const configFileContent = `invalid JSON content`

	// Create a temporary file simulating the config.json
	tmpfile, err := os.CreateTemp(".", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	// Write data to the tmp config file
	if _, err := tmpfile.Write([]byte(configFileContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	os.Setenv("CONFIG", tmpfile.Name())
	defer os.Unsetenv("CONFIG")

	config.GetConfig()

	assert.Equal(t, len(config.PostRules), 0)
}
