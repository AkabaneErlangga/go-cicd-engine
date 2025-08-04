package runner

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"go-cicd-engine/internal/model"
)

func Execute(j model.Job) error {
	// 1. Clone repo
	cloneDir := filepath.Join(os.TempDir(), "cicd-"+j.ID)
	err := gitClone(j.RepoURL, j.Branch, cloneDir)
	if err != nil {
		return err
	}

	// 2. Run all commands
	for _, cmdStr := range j.Commands {
		log.Printf("🔧 Executing: %s", cmdStr)

		cmd := exec.Command("bash", "-c", cmdStr)
		cmd.Dir = cloneDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func gitClone(repoURL, branch, dest string) error {
	log.Println("📥 Cloning repo...")
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", branch, repoURL, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

