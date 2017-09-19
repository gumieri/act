package editor

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func Open(editorCommand string, fileName, helperText string) (text string, err error) {
	filePath := fmt.Sprintf("%s/%s", os.TempDir(), fileName)

	tmpFile, err := os.Create(filePath)
	if err != nil {
		return
	}

	_, err = tmpFile.WriteString(helperText)
	if err != nil {
		return
	}

	tmpFile.Close()

	editorCmd := exec.Command(editorCommand, filePath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	err = editorCmd.Start()
	if err != nil {
		return
	}

	err = editorCmd.Wait()
	if err != nil {
		return
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	re := regexp.MustCompile("(?m)[\r\n]+^#.*$")
	text = strings.Trim(re.ReplaceAllString(string(content), ""), "\n")

	return
}
