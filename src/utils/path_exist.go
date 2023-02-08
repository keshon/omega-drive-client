package utils

import "os"

/*
	Path exist funcs allow to check if path ..well exists
	The implementation is taken from here: https://stackoverflow.com/a/10510783

	Innokentiy Sokolov
	https://github.com/keshon

	2022-03-25
*/

func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
