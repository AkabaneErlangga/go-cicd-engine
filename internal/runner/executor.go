package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"go-cicd-engine/internal/config"
	"go-cicd-engine/internal/model"
)

var ErrTimeout = errors.New("job timed out")
func Execute(j model.Job, out io.Writer) error {
	// 1. Clone repo
	cloneDir := filepath.Join(os.TempDir(), "cicd-"+j.ID)
	err := gitClone(j.RepoURL, j.Branch, cloneDir, out)
	if err != nil {
		return err
	}

	// 2. Load commands from .cicd.yaml
	log.Println("🔍 Loading .cicd.yaml...")
	commands, err := config.LoadConfig(cloneDir)
	if err != nil {	
  	log.Printf("❌ Failed to load config: %v", err)
		return err
	}
	if len(commands) == 0 {
		log.Println("⚠️  No .cicd.yaml found or no steps to run.")
		return nil
	}

	// 3. Run all commands
	for _, cmdStr := range commands {
		log.Printf("🔧 Executing: %s", cmdStr)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)

		cmd.Dir = cloneDir
		cmd.Stdout = out
		cmd.Stderr = out

		err := cmd.Run()

		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("❗ Command timed out: %s", cmdStr)
			return ErrTimeout
		}

		if err != nil {
			fmt.Printf("❌ Command failed: %s\nError: %v\n", cmdStr, err)
			return err
		}
	}

	return nil
}

func gitClone(repoURL, branch, dest string, out io.Writer) error {
	log.Println("📥 Cloning repo...")
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", branch, repoURL, dest)
	cmd.Stdout = out
	cmd.Stderr = out
	return cmd.Run()
}

