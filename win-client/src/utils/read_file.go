package utils

import (
	"io/ioutil"
	"log"
)

func ReadFile(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		log.Println("[utils/read_file.go][ReadFile] Error reading file. Details:")
		log.Println(err)
	}
	return b
}
