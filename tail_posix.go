// +build linux darwin freebsd

package tail

import (
	"os"
)

func OpenFile(name string, bool needLink) (file *os.File, err error) {
	file, err = os.Open(name)
	return file, err
}
