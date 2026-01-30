package gh

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var (
	sshRe   = regexp.MustCompile(`git@github\.com:(.+?)(?:\.git)?$`)
	httpsRe = regexp.MustCompile(`github\.com/(.+?)(?:\.git)?$`)
)

// DetectRepo extracts "owner/repo" from the git origin remote URL.
func DetectRepo() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", fmt.Errorf("no git remote found: %w", err)
	}
	url := strings.TrimSpace(string(out))

	if m := sshRe.FindStringSubmatch(url); m != nil {
		return m[1], nil
	}
	if m := httpsRe.FindStringSubmatch(url); m != nil {
		return m[1], nil
	}
	return "", fmt.Errorf("could not parse repo from remote URL: %s", url)
}
