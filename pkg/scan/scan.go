package scan

import (
	"os"
	"os/exec"
	"path/filepath"
)

func ScanSourceCode(semgrep_config string, outputDir string, dirName string) error {
	path := filepath.Join(outputDir, dirName)
	output := filepath.Join(path, "semgrep_output.json")
	cmd := exec.Command("semgrep", "--config", semgrep_config, path, "--json")
	outputFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	cmd.Dir = outputDir
	cmd.Stdout = outputFile
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}
