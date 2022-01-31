package download

import (
	"os/exec"
)

func DownloadGithubRepo(repoUrl string, outputDir string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", repoUrl)
	cmd.Dir = outputDir
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}