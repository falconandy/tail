// +build linux darwin freebsd

package tail

import (
	"os"
)

func OpenFile(name string) (file *os.File, isLink bool, err error) {
	file, err = os.Open(name)
	return file, false, err
}
