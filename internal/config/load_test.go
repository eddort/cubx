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
programs:
  - name: testProgram
    aliases: ["testalias"]
    image: testimage
    serializer: testhandler
    description: "Test Program description"
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

	if len(config.Programs) != 1 {
		t.Fatalf("Expected 1 Program, got %d", len(config.Programs))
	}

	cmd := config.Programs[0]
	if cmd.Name != "testProgram" {
		t.Errorf("Expected Program name 'testProgram', got '%s'", cmd.Name)
	}
	if cmd.Serializer != "testhandler" {
		t.Errorf("Expected Program handler 'testhandler', got '%s'", cmd.Serializer)
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
programs:
  - aliases: ["testalias"]
    image: testimage
    serializer: testhandler
    description: "Test Program description without name"
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
programs:
  - name: localProgram
    aliases: ["localalias"]
    image: localimage
    serializer: localhandler
    description: "Local Program description"
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
programs:
  - name: homeProgram
    aliases: ["homealias"]
    image: homeimage
    serializer: homehandler
    description: "Home Program description"
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

	if len(config.Programs) != 2 {
		t.Fatalf("Expected 2 Programs, got %d", len(config.Programs))
	}

	cmd1 := config.Programs[0]
	if cmd1.Name != "localProgram" {
		t.Errorf("Expected Program name 'localProgram', got '%s'", cmd1.Name)
	}
	if cmd1.Serializer != "localhandler" {
		t.Errorf("Expected Program handler 'localhandler', got '%s'", cmd1.Serializer)
	}

	cmd2 := config.Programs[1]
	if cmd2.Name != "homeProgram" {
		t.Errorf("Expected Program name 'homeProgram', got '%s'", cmd2.Name)
	}
	if cmd2.Serializer != "homehandler" {
		t.Errorf("Expected Program handler 'homehandler', got '%s'", cmd2.Serializer)
	}
}
