package editor

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Open the specified editor
/* Parameters:
 * editorCommand
 * fileName: the name of the temporary file
 * helperText: a text to be inserted on the file to be edited
 * removeSharps: a boolean to remove or not any line starting with #
 */
func Open(editorCommand string, fileName, helperText string, removeSharps bool) (text string, err error) {
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

	text = string(content)

	if removeSharps {
		re := regexp.MustCompile("(?m)[\r\n]+^#.*$")
		text = re.ReplaceAllString(text, "")
	}

	text = strings.Trim(text, "\n")

	return
}
