package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
)

func getHandler(input []byte, authToken string, instance string, endpoint string) []byte {

	connString := fmt.Sprintf("%s%s?api_token=%s", instance, endpoint, authToken)
	req, _ := http.NewRequest(
		"GET",
		connString,
		bytes.NewBuffer(input),
	)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(req) // send an HTTP using `req` object
	if err != nil {                        // check for response error
		log.Fatal("Error:O", err)
	}
	defer req.Body.Close()
	if res.StatusCode != 200 {
		log.Fatal("Error:O", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	return body
}

func tokenGenerator(apiKey string) string {
	if apiKey == "" {
		return ""
	}
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	h := sha256.New()
	h.Write([]byte(apiKey + timestamp))
	token := fmt.Sprintf("%x%s", h.Sum(nil), timestamp)
	return token
}

func deleteHandler(input []byte, authToken string, instance string, endpoint string) int {

	connString := fmt.Sprintf("%s%s?api_token=%s", instance, endpoint, authToken)
	req, _ := http.NewRequest(
		"DELETE",
		connString,
		bytes.NewBuffer(input),
	)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(req) // send an HTTP using `req` object
	if err != nil {                        // check for response error
		log.Fatal("Error:O", err)
	}
	defer res.Body.Close()
	log.Printf("POST status: %d\n", res.StatusCode)
	return res.StatusCode
}

func postHandler(input []byte, authToken string, instance string, endpoint string) int {
	connString := fmt.Sprintf("%s%s?api_token=%s", instance, endpoint, authToken)
	req, _ := http.NewRequest(
		"POST",
		connString,
		bytes.NewBuffer(input),
	)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(req) // send an HTTP using `req` object
	if err != nil {                        // check for response error
		log.Fatal("Error:O", err)
	}
	defer req.Body.Close()
	if res.StatusCode != 200 && res.StatusCode != 204 {
		log.Fatal("Error:O", res.StatusCode)
	}
	return res.StatusCode
}

type CallbackEntry struct {
	URL     string   `json:"url"`
	Method  string   `json:"method"`
	Headers []string `json:"headers"`
}

func retCodeHandler(ret int) {

	switch ret {
	case 204:
		color.Green("Operation completed successfully")
	case 401:
		color.Yellow("Unauthorized - Invalid API key")
	case 403:
		color.Yellow("Invalid unlock mode value")
	case 404:
		color.Red("Element not found")
	case 405:
		color.Red("Device disconnected")
	case 409:
		color.Red("Internal storage of Tedee Bridge is busy")
	default:
		color.Cyan("Ret code without message -  %s", ret)

	}
}
