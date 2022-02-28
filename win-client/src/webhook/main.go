package webhook

import (
	"app/src/conf"
	"app/src/utils"
	"log"
	"net/http"
	"strings"
)

// Verify if access key exists in Airtable
func PostToWebhook() Response {
	log.Println("[webhook/main.go][PostToWebhook] Validating access key has started")

	// Prepare payload
	type Payload struct {
		AccessKey string
	}
	var payload Payload

	// Decrypt access key
	key, err := utils.Decrypt([]byte(conf.EncryptKey), utils.ReadFile("key"))

	//log.Println(string(key))

	if err != nil {
		log.Println("[webhook/main.go][PostToWebhook] Error decrypting access key. Details:")
		log.Println(err)
	}

	// Clean access key
	payload.AccessKey = strings.TrimSuffix(string(key), "\r\n")

	// Access key is missing in `key` file
	if len(payload.AccessKey) == 0 {
		log.Println("[webhook/main.go][PostToWebhook] Warning: access key is missing. Do nothing")
		return Response{}
	}

	// Which URL to use. N8n offers two of them
	var path string
	if conf.IsDev {
		path = conf.N8nDevURL // Development webhook
	} else {
		path = conf.N8nURL // Production webhook
	}

	// Request
	response := []Response{}
	utils.Request(conf.N8nAuthEncoded, http.MethodPost, path, payload, &response)

	// Access key not found
	if len(response) <= 0 {
		log.Println("[webhook/main.go][PostToWebhook] Warning: access key was not found in Airtable. Do nothing")
		return Response{}
	}

	log.Println("[webhook/main.go][PostToWebhook] Validating access key has finished.")

	return response[0]
}
