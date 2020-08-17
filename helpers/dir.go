package helpers

import (
	"errors"
	"os"
)

// CreateDestDir check and create directory if not exist
func CreateDestDir(name string) error {
	destInfo, err := os.Stat(name)
	if os.IsNotExist(err) {
		errCreateDir := os.MkdirAll(name, 0755)
		if errCreateDir != nil {
			return err
		}
		return nil
	} else if !destInfo.IsDir() {
		return errors.New("Can not create dir, destination is file but should be a folder")
	} else {
		return err
	}
}
