package utils

import (
	"io/ioutil"
	"log"
)

/*
	Read file funcs reads file content to bytes

	Innokentiy Sokolov
	https://github.com/keshon

	2022-03-25
*/

func ReadFile(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		log.Println(err)
	}
	return b
}
