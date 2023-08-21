package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func ExampleHandler_CreateShortURL() {
	req, _ := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://localhost:8080/"), strings.NewReader("example.org"))
	req.Header.Add("Content-Type", "text/plain")

	client := &http.Client{}

	resp, _ := client.Do(req)

	defer resp.Body.Close()
}

func ExampleHandler_GetByShort() {

	req, _ := http.NewRequest(http.MethodGet,
		fmt.Sprintf("http://localhost:8080/short_url.com"), nil)
	req.Header.Add("Content-Type", "text/plain")

	client := &http.Client{}

	resp, _ := client.Do(req)

	defer resp.Body.Close()
}

func ExampleHandler_CreateShortURLByJSON() {

	data, _ := json.Marshal(map[string]string{
		"url": "https://ya.ru",
	},
	)

	req, _ := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://localhost:8080/api/shorten/"), bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, _ := client.Do(req)

	defer resp.Body.Close()
}

func ExampleHandler_GetUserURLs() {

	req, _ := http.NewRequest(http.MethodGet,
		fmt.Sprintf("http://localhost:8080/api/shorten/user/urls/"), nil)
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "authToken", Value: "user_token_value"})

	client := &http.Client{}

	resp, _ := client.Do(req)

	defer resp.Body.Close()
}

func ExampleHandler_DeleteUserURLsBatch() {

	data, _ := json.Marshal(
		[]string{"https://short_url", "https://short_url_2"},
	)

	req, _ := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("http://localhost:8080/api/shorten/user/urls/"), bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "authToken", Value: "user_token_value"})

	client := &http.Client{}

	resp, _ := client.Do(req)

	defer resp.Body.Close()
}

func ExampleHandler_BatchCreateURLByJSON() {

	data, _ := json.Marshal(
		[]map[string]string{
			{
				"url":            "https://ya.ru",
				"correlation_id": "vfwt4312",
			},
			{
				"url":            "https://yand.ru",
				"correlation_id": "fwef13",
			},
		},
	)

	req, _ := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://localhost:8080/api/shorten/batch"), bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, _ := client.Do(req)

	defer resp.Body.Close()
}
