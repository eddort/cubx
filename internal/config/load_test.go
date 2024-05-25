package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func setupTempDir(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	cleanup := func() { os.RemoveAll(tempDir) }
	return tempDir, cleanup
}

func TestLoadValidConfig(t *testing.T) {
	tempDir, cleanup := setupTempDir(t)
	defer cleanup()

	cubxDir := filepath.Join(tempDir, ".cubx")
	if err := os.Mkdir(cubxDir, 0755); err != nil {
		t.Fatalf("Failed to create .cubx dir: %v", err)
	}

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

	homeDir := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", homeDir)

	config, loadedConfigs, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expectedLoadedConfigs := []string{
		validConfigPath,
	}

	if len(loadedConfigs) != len(expectedLoadedConfigs) || loadedConfigs[0] != expectedLoadedConfigs[0] {
		t.Fatalf("Expected config file used: %v, got: %v", expectedLoadedConfigs, loadedConfigs)
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
}

func TestLoadInvalidConfig(t *testing.T) {
	tempDir, cleanup := setupTempDir(t)
	defer cleanup()

	cubxDir := filepath.Join(tempDir, ".cubx")
	if err := os.Mkdir(cubxDir, 0755); err != nil {
		t.Fatalf("Failed to create .cubx dir: %v", err)
	}

	invalidConfigPath := filepath.Join(cubxDir, "config.yaml")
	invalidConfigContent := []byte(`
commands:
  - aliases: ["testalias"]
    image: testimage
    handler: testhandler
    description: "Test command description without name"
`)
	if err := os.WriteFile(invalidConfigPath, invalidConfigContent, 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	homeDir := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", homeDir)

	_, _, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for invalid config, got nil")
	}
}

func TestMergeConfig(t *testing.T) {
	tempDir, cleanup := setupTempDir(t)
	defer cleanup()

	// Create local configuration in .cubx
	localCubxDir := filepath.Join(tempDir, ".cubx")
	if err := os.Mkdir(localCubxDir, 0755); err != nil {
		t.Fatalf("Failed to create .cubx dir: %v", err)
	}

	localConfigPath := filepath.Join(localCubxDir, "config.yaml")
	localConfigContent := []byte(`
commands:
  - name: localcommand
    aliases: ["localalias"]
    image: localimage
    handler: localhandler
    description: "Local command description"
`)
	if err := os.WriteFile(localConfigPath, localConfigContent, 0644); err != nil {
		t.Fatalf("Failed to write local config file: %v", err)
	}

	// Create a configuration in the home directory
	homeCubxDir := filepath.Join(tempDir, "home", ".cubx")
	if err := os.MkdirAll(homeCubxDir, 0755); err != nil {
		t.Fatalf("Failed to create home .cubx dir: %v", err)
	}

	homeConfigPath := filepath.Join(homeCubxDir, "config.yaml")
	homeConfigContent := []byte(`
commands:
  - name: homecommand
    aliases: ["homealias"]
    image: homeimage
    handler: homehandler
    description: "Home command description"
`)
	if err := os.WriteFile(homeConfigPath, homeConfigContent, 0644); err != nil {
		t.Fatalf("Failed to write home config file: %v", err)
	}

	// Set up a temporary directory for the home configuration
	homeDir := os.Getenv("HOME")
	os.Setenv("HOME", filepath.Join(tempDir, "home"))
	fmt.Println("TMP HOME", os.Getenv("HOME"), homeConfigPath)
	defer os.Setenv("HOME", homeDir)

	// Go to tempDir to emulate the current directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	config, loadedConfigs, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expectedLoadedConfigs := []string{
		homeConfigPath,
		localConfigPath,
	}

	if len(loadedConfigs) != len(expectedLoadedConfigs) {
		t.Fatalf("Expected %d loaded config files, got %d: %v", len(expectedLoadedConfigs), len(loadedConfigs), loadedConfigs)
	}

	if len(config.Commands) != 2 {
		t.Fatalf("Expected 2 commands, got %d", len(config.Commands))
	}

	cmd1 := config.Commands[0]
	if cmd1.Name != "localcommand" {
		t.Errorf("Expected command name 'localcommand', got '%s'", cmd1.Name)
	}
	if cmd1.Handler != "localhandler" {
		t.Errorf("Expected command handler 'localhandler', got '%s'", cmd1.Handler)
	}

	cmd2 := config.Commands[1]
	if cmd2.Name != "homecommand" {
		t.Errorf("Expected command name 'homecommand', got '%s'", cmd2.Name)
	}
	if cmd2.Handler != "homehandler" {
		t.Errorf("Expected command handler 'homehandler', got '%s'", cmd2.Handler)
	}
}
