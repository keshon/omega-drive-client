package settings

import (
	"app/src/conf"
	"app/src/utils"
	"io/ioutil"
	"log"
	"strings"
)

func ReadAccessKey() string {
	strData := strings.TrimSuffix(string(utils.ReadFile("key")), "\r\n")

	if len(strData) == 0 {
		log.Println("[settings/settings.go][ReadAccessKey] Access key is empty")
	}

	return strData
}

func WriteAccessKey(strData string) bool {
	log.Println([]byte(strData))

	// Simple validate
	strData = strings.TrimSuffix(string(strData), "\r\n")
	if len(strData) == 0 {
		log.Println("[settings/settings.go][WriteAccessKey] Access key is empty")
		return false
	}

	// Encrypting
	key := []byte(conf.EncryptKey) // 32 bytes
	plaintext := []byte(strData)
	//fmt.Printf("%s\n", plaintext)

	ciphertext, err := utils.Encrypt(key, plaintext)
	if err != nil {
		log.Println("[settings/settings.go][WriteAccessKey] Error encrypting access key. Details:")
		log.Println(err)
	}

	// write the whole body at once
	err = ioutil.WriteFile("key", []byte(ciphertext), 0644)
	if err != nil {
		log.Println("[settings/settings.go][WriteAccessKey] Error write access key to file key. Details:")
		log.Println(err)
	}

	return true
}
