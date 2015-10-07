// +build windows

package tail

import (
	"github.com/hpcloud/tail/winfile"
	"os"
	"io/ioutil"
	"path/filepath"
	"errors"
	"strings"
)

func OpenFile(name string, needLink bool, linkDirectory string) (file *os.File, err error) {
	linkFileName := name
	if needLink {
		linkFileName, err = createLink(name, linkDirectory)
		if err != nil {
			return nil, err
		}
	}
	file, err = winfile.OpenFile(linkFileName, os.O_RDONLY, 0)
	return file, err
}

func createLink(originalName string, linkDirectory string) (linkFileName string, err error) {
	linkDirectoryPath, err := getTempLinkDirectoryPath(originalName, linkDirectory)
	if err != nil {
		return "", err
	}
	linkFileName, err = getUniqueTempFileName(linkDirectoryPath)
	if err != nil {
		return "", err
	}
	err = os.Link(originalName, linkFileName)
	if err != nil {
		return "", err
	}
	return linkFileName, nil
}

func getTempLinkDirectoryPath(originalName string, linkDirectory string) (linkDirectoryPath string, err error) {
	originalPath, err := filepath.Abs(originalName)
	if err != nil {
		return "", err
	}
	if linkDirectory != "" {
		configLinkDirectoryPath, err := filepath.Abs(linkDirectory)
		if err != nil {
			return "", err
		}
		if strings.ToLower(filepath.VolumeName(configLinkDirectoryPath)) != strings.ToLower(filepath.VolumeName(originalPath)) {
			return "", errors.New("Volumes of files and directory for temp hard links must be identical")
		}
		return configLinkDirectoryPath, nil
	}
	tempDirectoryPath, err := filepath.Abs(os.TempDir())
	if err != nil {
		return "", err
	}
	if strings.ToLower(filepath.VolumeName(originalPath)) == strings.ToLower(filepath.VolumeName(tempDirectoryPath)) {
		return tempDirectoryPath, nil
	}
	return filepath.Dir(originalPath), nil
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
	err = os.Remove(tempFileName)
	if err != nil {
		return "", err
	}
	return tempFileName, nil
}