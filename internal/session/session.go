package session

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Run(settings Settings) error {
	shellPath := os.Getenv("SHELL")
	shell := filepath.Base(shellPath)
	if shell != "bash" && shell != "zsh" {
		return fmt.Errorf("unsupported shell %s: this script only supports 'bash' or 'zsh'", shell)
	}

	tempDir, err := os.MkdirTemp("", "shellrc")
	if err != nil {
		return fmt.Errorf("failed to create a temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	localRc := filepath.Join(os.Getenv("HOME"), "."+shell+"rc")
	tempRcPath := filepath.Join(tempDir, "."+shell+"rc")

	if _, err := os.Stat(localRc); err == nil {
		input, err := os.ReadFile(localRc)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", localRc, err)
		}
		if err = os.WriteFile(tempRcPath, input, 0644); err != nil {
			return fmt.Errorf("failed to write to %s: %w", tempRcPath, err)
		}
	} else {
		if _, err := os.Create(tempRcPath); err != nil {
			return fmt.Errorf("failed to create a file %s: %w", tempRcPath, err)
		}
	}

	file, err := os.OpenFile(tempRcPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s for appending: %w", tempRcPath, err)
	}
	defer file.Close()

	for _, setting := range settings.SettingsToStrings() {
		if _, err := file.WriteString(setting); err != nil {
			return fmt.Errorf("failed to write setting to file %s: %w", tempRcPath, err)
		}
	}

	cmd := exec.Command(shell, "-i")
	env := os.Environ()

	if shell == "zsh" {
		env = append(env, "ZDOTDIR="+tempDir)
	} else {
		env = append(env, "HOME="+tempDir)
	}

	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run the shell %s: %w", shell, err)
	}

	return nil
}
