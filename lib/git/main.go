package git

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func TopLevelPath(gitPath string) (path string, err error) {
	out, err := exec.Command(gitPath, "rev-parse", "--show-toplevel").Output()
	path = strings.Trim(string(out), "\n")
	return
}

func IssueIdFromBranch(gitPath string, gitRegex string) (issueId int, err error) {
	out, _ := exec.Command(gitPath, "rev-parse", "--abbrev-ref", "HEAD").Output()

	regexC, err := regexp.Compile(gitRegex)

	if err != nil {
		return
	}

	issueId, err = strconv.Atoi(regexC.FindString(string(out)))

	return
}
