// +build windows

package tail

import (
	"path"
	"github.com/hpcloud/tail/winfile"
	"os"
	"io/ioutil"
)

func OpenFile(name string) (file *os.File, isLink bool, err error) {
	tempFile, err := ioutil.TempFile(path.Dir(name), "")
	if err != nil {
		return nil, false, err
	}
	tempFileName := tempFile.Name()
	err = tempFile.Close()
	if err != nil {
		return nil, false, err
	}
	os.Remove(tempFileName)
	err = os.Link(name, tempFileName)
	if err != nil {
		return nil, false, err
	}
	file, err = winfile.OpenFile(tempFileName, os.O_RDONLY, 0)
	return file, true, err
}
