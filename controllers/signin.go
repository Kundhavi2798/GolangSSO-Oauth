package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var ssogolang *oauth2.Config
var RandomString = "random-string"

type Credentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

func init() {
	file, err := os.Open("evn.json")
	if err != nil {
		log.Fatalf("Error opening credentials file: %v", err)
	}
	defer file.Close()
	fmt.Println("file", file)
	var creds Credentials
	err = json.NewDecoder(file).Decode(&creds)
	if err != nil {
		log.Fatalf("Error parsing credentials file: %v", err)
	}

	ssogolang = &oauth2.Config{
		RedirectURL:  creds.RedirectURL,
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func Signin(w http.ResponseWriter, r *http.Request) {
	url := ssogolang.AuthCodeURL(RandomString)
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
