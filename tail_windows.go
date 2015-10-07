// +build windows

package tail

import (
	"path"
	"github.com/hpcloud/tail/winfile"
	"os"
	"io/ioutil"
	"path/filepath"
)

func OpenFile(name string, needLink bool) (file *os.File, err error) {
	linkFileName := name
	if needLink {
		linkFileName, err = createLink(name)
		if err != nil {
			return nil, err
		}
	}
	file, err = winfile.OpenFile(linkFileName, os.O_RDONLY, 0)
	return file, err
}

func createLink(originalName string) (linkFileName string, err error) {
	originalPath, err := filepath.Abs(originalName)
	if err != nil {
		return "", err
	}
	tempDirPath, err := filepath.Abs(os.TempDir())
	if err != nil {
		return "", err
	}
	if filepath.VolumeName(originalPath) != filepath.VolumeName(tempDirPath) {
		tempDirPath = path.Dir(originalPath)
	}
	linkFileName, err = getUniqueTempFileName(tempDirPath)
	if err != nil {
		return "", err
	}
	err = os.Link(originalName, linkFileName)
	if err != nil {
		return "", err
	}
	return linkFileName, nil
}

func getUniqueTempFileName(tempDir string) (name string, err error) {
	tempFile, err := ioutil.TempFile(tempDir, "~")
	if err != nil {
		return "", err
	}
	tempFileName := tempFile.Name()
	err = tempFile.Close()
	if err != nil {
		return "", err
	}
	os.Remove(tempFileName)
	return tempFileName, nil
}