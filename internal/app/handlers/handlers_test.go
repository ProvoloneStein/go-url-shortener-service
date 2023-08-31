package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

func ExampleHandler_CreateShortURL() {
	req, err := http.NewRequest(http.MethodPost,
		"http://localhost:8080/", strings.NewReader("example.org"))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "text/plain")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func ExampleHandler_GetByShort() {

	req, err := http.NewRequest(http.MethodGet,
		"http://localhost:8080/short_url.com", nil)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "text/plain")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func ExampleHandler_CreateShortURLByJSON() {

	data, err := json.Marshal(map[string]string{
		"url": "https://ya.ru",
	},
	)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost,
		"http://localhost:8080/api/shorten/", bytes.NewReader(data))

	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func ExampleHandler_GetUserURLs() {

	req, err := http.NewRequest(http.MethodGet,
		"http://localhost:8080/api/shorten/user/urls/", nil)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "authToken", Value: "user_token_value"})

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func ExampleHandler_DeleteUserURLsBatch() {

	data, err := json.Marshal(
		[]string{"https://short_url", "https://short_url_2"},
	)

	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodDelete,
		"http://localhost:8080/api/shorten/user/urls/", bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "authToken", Value: "user_token_value"})

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func ExampleHandler_BatchCreateURLByJSON() {

	data, err := json.Marshal(
		[]map[string]string{
			{
				"original_url":   "https://ya.ru",
				"correlation_id": "vfwt4312",
			},
			{
				"original_url":   "https://yand.ru",
				"correlation_id": "fwef13",
			},
		},
	)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost,
		"http://localhost:8080/api/shorten/batch", bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}
