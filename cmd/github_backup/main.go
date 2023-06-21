package main

import (
	"github.com/niranjan94/github-backup/internal"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func main() {
	rawOwners := strings.TrimSpace(os.Getenv("GITHUB_BACKUP_OWNERS"))
	rawToken := strings.TrimSpace(os.Getenv("GITHUB_BACKUP_TOKEN"))
	rawDirectory := strings.TrimSpace(os.Getenv("GITHUB_BACKUP_DIRECTORY"))
	rawConcurrency := strings.TrimSpace(os.Getenv("GITHUB_BACKUP_CONCURRENCY"))
	rawPeriodSeconds := strings.TrimSpace(os.Getenv("GITHUB_BACKUP_PERIOD_SECONDS"))

	if rawOwners == "" {
		logrus.Fatalf("GITHUB_BACKUP_OWNERS is not set")
	}

	if rawToken == "" {
		logrus.Fatalf("GITHUB_BACKUP_TOKEN is not set")
	}

	if rawDirectory == "" {
		logrus.Fatalf("GITHUB_BACKUP_DIRECTORY is not set")
	}

	if rawConcurrency == "" {
		rawConcurrency = "10"
	}

	if rawPeriodSeconds == "" {
		rawPeriodSeconds = "86400" // 24 hours
	}

	config := &internal.Config{
		Token:     rawToken,
		Directory: rawDirectory,
		Owners:    []string{},
	}

	if periodSeconds, err := strconv.Atoi(rawPeriodSeconds); err != nil {
		logrus.Fatalf("GITHUB_BACKUP_PERIOD_SECONDS is not a number")
	} else {
		config.PeriodSeconds = periodSeconds
	}

	if concurrency, err := strconv.Atoi(rawConcurrency); err != nil {
		logrus.Fatalf("GITHUB_BACKUP_CONCURRENCY is not a number")
	} else {
		config.Concurrency = concurrency
	}

	for _, rawOwner := range strings.Split(rawOwners, ",") {
		rawOwner := strings.TrimSpace(rawOwner)
		if !internal.ValidateGithubName(rawOwner) {
			logrus.Warnf("Ignoring invalid owner: %s\n", rawOwner)
			continue
		}
		config.Owners = append(config.Owners, strings.ToLower(rawOwner))
	}

	if len(config.Owners) == 0 {
		logrus.Fatalf("No owners set")
	}

	internal.NewBackupSchedule(config).Start()
}
