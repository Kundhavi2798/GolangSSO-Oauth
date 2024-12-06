package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	if state != RandomString {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	data, err := getUserData(code)
	if err != nil {
		http.Error(w, "Error fetching user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Data: %s", data)
}

func getUserData(code string) ([]byte, error) {
	token, err := ssogolang.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken
	response, err := http.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer response.Body.Close()

	// Check if the API response is OK
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("non-OK response from Google API: %s", body)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}
