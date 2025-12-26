// Package docker
package docker

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"horizonx-server/internal/logger"
)

type StreamHandler func(line string, isErr bool)

type Manager struct {
	log     logger.Logger
	workDir string
}

func NewManager(log logger.Logger, workDir string) *Manager {
	return &Manager{
		log:     log,
		workDir: workDir,
	}
}

func (m *Manager) Initialize() error {
	if err := os.MkdirAll(m.workDir, 0o755); err != nil {
		return fmt.Errorf("failed to create work directory: %w", err)
	}

	m.log.Info("docker manager initialized", "work_dir", m.workDir)
	return nil
}

func (m *Manager) GetAppDir(appID int64) string {
	return filepath.Join(m.workDir, fmt.Sprintf("app-%d", appID))
}

func (m *Manager) ComposeUp(ctx context.Context, appID int64, detached, build bool, onStream StreamHandler) (string, error) {
	appDir := m.GetAppDir(appID)

	args := []string{"compose", "up"}
	if detached {
		args = append(args, "-d")
	}
	if build {
		args = append(args, "--build")
	}

	return m.Run(ctx, appDir, onStream, args...)
}

func (m *Manager) ComposeDown(ctx context.Context, appID int64, removeVolumes bool, onStream StreamHandler) (string, error) {
	appDir := m.GetAppDir(appID)

	args := []string{"compose", "down"}
	if removeVolumes {
		args = append(args, "-v")
	}

	return m.Run(ctx, appDir, onStream, args...)
}

func (m *Manager) ComposeStop(ctx context.Context, appID int64, onStream StreamHandler) (string, error) {
	appDir := m.GetAppDir(appID)
	return m.Run(ctx, appDir, onStream, "compose", "stop")
}

func (m *Manager) ComposeStart(ctx context.Context, appID int64, onStream StreamHandler) (string, error) {
	appDir := m.GetAppDir(appID)
	return m.Run(ctx, appDir, onStream, "compose", "start")
}

func (m *Manager) ComposeRestart(ctx context.Context, appID int64, onStream StreamHandler) (string, error) {
	appDir := m.GetAppDir(appID)
	return m.Run(ctx, appDir, onStream, "compose", "restart")
}

func (m *Manager) ComposeLogs(ctx context.Context, appID int64, tail int, onStream StreamHandler) (string, error) {
	appDir := m.GetAppDir(appID)
	args := []string{"compose", "logs"}
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}

	return m.Run(ctx, appDir, onStream, args...)
}

func (m *Manager) ComposePs(ctx context.Context, appID int64, onStream StreamHandler) (string, error) {
	appDir := m.GetAppDir(appID)
	return m.Run(ctx, appDir, onStream, "compose", "ps")
}

func (m *Manager) ValidateDockerComposeFile(appID int64) error {
	appDir := m.GetAppDir(appID)
	files := []string{
		"docker-compose.yml",
		"docker-compose.yaml",
		"compose.yml",
		"compose.yaml",
	}

	for _, f := range files {
		if _, err := os.Stat(filepath.Join(appDir, f)); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no docker-compose file found")
}

func (m *Manager) WriteEnvFile(appID int64, envVars map[string]string) error {
	appDir := m.GetAppDir(appID)
	envPath := filepath.Join(appDir, ".env")

	var buf bytes.Buffer
	for k, v := range envVars {
		v = strings.ReplaceAll(v, "\n", "\\n")
		buf.WriteString(fmt.Sprintf("%s=\"%s\"\n", k, v))
	}

	return os.WriteFile(envPath, buf.Bytes(), 0o600)
}

func (m *Manager) IsDockerInstalled() bool {
	return exec.Command("docker", "--version").Run() == nil
}

func (m *Manager) IsDockerComposeAvailable() bool {
	return exec.Command("docker", "compose", "version").Run() == nil
}

func (m *Manager) Run(ctx context.Context, workDir string, onStream StreamHandler, args ...string) (string, error) {
	var buf bytes.Buffer

	err := m.RunStream(ctx, workDir, func(line string, isErr bool) {
		buf.WriteString(line + "\n")
		if onStream != nil {
			onStream(line, isErr)
		}
	}, args...)

	return buf.String(), err
}

func (m *Manager) RunStream(ctx context.Context, workDir string, onStream StreamHandler, args ...string) error {
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Dir = workDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	errCount := 2
	done := make(chan error, errCount)

	stream := func(r io.Reader, isErr bool) {
		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 1024), 1024*1024)

		for scanner.Scan() {
			onStream(scanner.Text(), isErr)
		}

		done <- scanner.Err()
	}

	go stream(stdout, false)
	go stream(stderr, true)

	for range errCount {
		if err := <-done; err != nil {
			return err
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("docker command failed: %w", err)
	}

	return nil
}
