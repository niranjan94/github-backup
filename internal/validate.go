package internal

import (
	"regexp"
)

// ValidateGithubName validates a GitHub username/organization name
func ValidateGithubName(name string) bool {
	match, _ := regexp.MatchString("^[\\-\\.\\w]*$", name)
	if !match {
		return false
	}
	return true
}
