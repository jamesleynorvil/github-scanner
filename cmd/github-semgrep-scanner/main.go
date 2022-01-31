package main

import (
	"flag"
	"os"

	scanner "github.com/jamesleynorvil/github-semgrep-scanner/internal/scanner"
	log "github.com/sirupsen/logrus"
)

const (
	SCAN   = "scan"
	AUDIT  = "audit"
	SEARCH = "search"
)

var (
	limit          int
	query          string
	semgrep_config string
	keep_source    bool
	scanFlagSet    *flag.FlagSet
	auditFlagSet   *flag.FlagSet
	searchFlagSet  *flag.FlagSet
)

func init() {
	searchFlagSet = flag.NewFlagSet(SEARCH, flag.ExitOnError)
	searchFlagSet.StringVar(&query, "query", "", "GitHub search query.")
	searchFlagSet.IntVar(&limit, "limit", 0, "Maximum number of repository to return from the search.")
	auditFlagSet = flag.NewFlagSet(AUDIT, flag.ExitOnError)
	auditFlagSet.StringVar(&query, "query", "", "GitHub search query.")
	auditFlagSet.IntVar(&limit, "limit", 0, "Maximum number of repository to return from the search.")
	auditFlagSet.StringVar(&semgrep_config, "semgrep-config", "", "Location of the semgrep rules.")
	auditFlagSet.BoolVar(&keep_source, "keep-source", false, "Whether to keep the source code after the audit.")
	scanFlagSet = flag.NewFlagSet(SCAN, flag.ExitOnError)
	scanFlagSet.StringVar(&query, "query", "", "GitHub search query.")
	scanFlagSet.IntVar(&limit, "limit", 0, "Maximum number of repository to return from the search.")
	scanFlagSet.StringVar(&semgrep_config, "semgrep-config", "", "Location of the semgrep rules.")
	scanFlagSet.BoolVar(&keep_source, "keep-source", false, "Whether to keep the source code after the scan.")
}

func handleSearchCmd() {
	searchFlagSet.Parse(os.Args[2:])
	if query == "" || limit < 1 {
		searchFlagSet.Usage()
		os.Exit(1)
	}
	if err := scanner.LaunchGithubSearch(query, limit); err != nil {
		log.Errorf("Unable to launch github search. %v", err)
	}
}

func handleScanCmd() {
	scanFlagSet.Parse(os.Args[2:])
	if query == "" || semgrep_config == "" || limit < 1 {
		scanFlagSet.Usage()
		os.Exit(1)
	}
	if err := scanner.LaunchScan(query, limit, semgrep_config, keep_source, false); err != nil {
		log.Errorf("Unable to launch scan. %v", err)
	}
}

func handleAuditCmd() {
	auditFlagSet.Parse(os.Args[2:])
	if query == "" || semgrep_config == "" || limit < 1 {
		auditFlagSet.Usage()
		os.Exit(1)
	}
	if err := scanner.LaunchScan(query, limit, semgrep_config, keep_source, true); err != nil {
		log.Errorf("Unable to launch audit scan. %v", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Errorf("Expected a command: {%s | %s | %s}", SEARCH, AUDIT, SCAN)
		flag.PrintDefaults()
		os.Exit(1)
	}
	switch os.Args[1] {
	case SEARCH:
		handleSearchCmd()
	case SCAN:
		handleScanCmd()
	case AUDIT:
		handleAuditCmd()
	default:
		log.Errorf("Unkown Command: %v. Expected any of {%s | %s | %s}", os.Args[1], SEARCH, AUDIT, SCAN)
		flag.PrintDefaults()
		os.Exit(1)
	}
}
