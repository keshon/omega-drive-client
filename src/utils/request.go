package utils

import (
	"app/src/conf"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
	Request wrapper simplifies use of http request by wrapping most common operations

	Innokentiy Sokolov
	https://github.com/keshon

	2022-03-29
*/

func Request(authBasic, method, url string, payload, unmarshal interface{}) error {
	// Marshal payload to bytes
	var body io.ReadWriter

	if payload == nil {
		payload = strings.NewReader("")
	}

	if payload != nil {
		buf, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(buf)
	}

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	if authBasic != "" {
		req.Header.Set("Authorization", "Basic "+authBasic)
	}

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Print response
	if conf.IsDev {
		fmt.Println(string(respBody))
	}

	// Unmarshal response
	if unmarshal != nil {
		err = json.Unmarshal([]byte(respBody), &unmarshal)
		//err = json.NewDecoder(resp.Body).Decode(&unmarshal)
		if err != nil {
			return err
		}
	}

	return nil
}
