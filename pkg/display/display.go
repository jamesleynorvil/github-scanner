package display

import (
	log "github.com/sirupsen/logrus"
)

func ShowMatchedReposList(repos []string) error {
	for _, repo := range repos {
		log.Info(repo)
	}
	return nil
}

func ShowScanResults(outputDir string) error {
	return nil
}

func ShowAuditScanResults(outputDir string) error {
	return nil
}
