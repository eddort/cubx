package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create temporary subdirectory for .cubx
	cubxDir := filepath.Join(tempDir, ".cubx")
	if err := os.Mkdir(cubxDir, 0755); err != nil {
		t.Fatalf("Failed to create .cubx dir: %v", err)
	}

	// Create temporary file with valid configuration in the current directory
	validConfigPath := filepath.Join(cubxDir, "config.yaml")
	validConfigContent := []byte(`
commands:
  - name: testcommand
    aliases: ["testalias"]
    image: testimage
    handler: testhandler
    description: "Test command description"
`)
	if err := os.WriteFile(validConfigPath, validConfigContent, 0644); err != nil {
		t.Fatalf("Failed to write valid config file: %v", err)
	}

	// Create temporary file with invalid configuration
	invalidConfigContent := []byte(`
commands:
  - aliases: ["testalias"]
    image: testimage
    handler: testhandler
    description: "Test command description without name"
`)
	if err := os.WriteFile(filepath.Join(cubxDir, "config_invalid.yaml"), invalidConfigContent, 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	// Set the environment variable for the home directory
	homeDir := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", homeDir)

	// Test loading valid configuration
	t.Run("valid config", func(t *testing.T) {
		config, loadedConfigs, err := LoadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if len(loadedConfigs) == 0 || loadedConfigs[0] != validConfigPath {
			t.Fatalf("Expected config file used: %s, got: %v", validConfigPath, loadedConfigs)
		}

		if len(config.Commands) != 1 {
			t.Fatalf("Expected 1 command, got %d", len(config.Commands))
		}

		cmd := config.Commands[0]
		if cmd.Name != "testcommand" {
			t.Errorf("Expected command name 'testcommand', got '%s'", cmd.Name)
		}
		if cmd.Handler != "testhandler" {
			t.Errorf("Expected command handler 'testhandler', got '%s'", cmd.Handler)
		}
	})

	// Test loading invalid configuration
	t.Run("invalid config", func(t *testing.T) {
		if err := os.WriteFile(filepath.Join(cubxDir, "config.yaml"), invalidConfigContent, 0644); err != nil {
			t.Fatalf("Failed to write invalid config file: %v", err)
		}

		_, _, err := LoadConfig()
		if err == nil {
			t.Fatal("Expected error for invalid config, got nil")
		}
	})
}
