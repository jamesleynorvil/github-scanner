package scanner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	display "github.com/jamesleynorvil/github-semgrep-scanner/pkg/display"
	download "github.com/jamesleynorvil/github-semgrep-scanner/pkg/download"
	scan "github.com/jamesleynorvil/github-semgrep-scanner/pkg/scan"
	search "github.com/jamesleynorvil/github-semgrep-scanner/pkg/search"
	log "github.com/sirupsen/logrus"
)

func initiateScanningProcess(query string, limit int, outputDir string) error {
	download_wg := &sync.WaitGroup{}
	matchedRepos, err := search.GetMatchedReposList(query, limit)
	if err != nil {
		return err
	}
	if len(matchedRepos) == 0 {
		errorMsg := fmt.Sprintf("The Github query [%s] did not return any result.", query)
		return errors.New(errorMsg)
	}
	for _, repo := range matchedRepos {
		download_wg.Add(1)
		go func(repo string, outputDir string, download_wg *sync.WaitGroup) {
			defer download_wg.Done()
			log.Infof("Dowloading: %s", repo)
			if err := download.DownloadGithubRepo(repo, outputDir); err != nil {
				log.Warnf("Failed to download: %s. %v", repo, err)
			}
		}(repo, outputDir, download_wg)
		download_wg.Wait()
	}
	return nil
}

func LaunchGithubSearch(query string, limit int) error {
	reposList, err := search.GetMatchedReposList(query, limit)
	if err != nil {
		return err
	}
	if len(reposList) == 0 {
		log.Errorf("The Github query [%s] did not return any result.", query)
		return nil
	}
	if err := display.ShowMatchedReposList(reposList); err != nil {
		log.Errorf("Unable to display the list of matched repositories. %v", err)
	}
	return nil
}

func LaunchScan(query string, limit int, semgrep_config string, keep_source bool, isAudit bool) error {
	outputDir, err := os.MkdirTemp("", "github-semgrep-scanner-output-")
	if err != nil {
		return err
	}
	if !keep_source {
		defer os.RemoveAll(outputDir)
	}
	if err = initiateScanningProcess(query, limit, outputDir); err != nil {
		return err
	}
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return err
	}
	scan_wg := &sync.WaitGroup{}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		scan_wg.Add(1)
		go func(outputDir string, dirName string, audit_wg *sync.WaitGroup) {
			defer audit_wg.Done()
			path := filepath.Join(outputDir, dirName)
			log.Infof("Scanning source code at: %s", path)
			if err := scan.ScanSourceCode(semgrep_config, outputDir, dirName); err != nil {
				log.Warnf("Failed to scan source code at: %s. %v", path, err)
			}
		}(outputDir, file.Name(), scan_wg)
		scan_wg.Wait()
	}
	if isAudit {
		if err := display.ShowAuditScanResults(outputDir); err != nil {
			log.Errorf("Unable to display audit scan results. %v", err)
		}
	} else {
		if err := display.ShowScanResults(outputDir); err != nil {
			log.Errorf("Unable to display scan results. %v", err)
		}
	}
	return nil
}
