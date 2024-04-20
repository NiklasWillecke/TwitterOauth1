package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

var secret string

func handleHome(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="login">Twitter Login</a></body></html>`
	fmt.Fprint(w, html)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {

	requestToken, requestSecret, _ := config.RequestToken()

	authorizationURL, _ := config.AuthorizationURL(requestToken)

	secret = requestSecret
	// handle err
	http.Redirect(w, r, authorizationURL.String(), http.StatusFound)

	fmt.Println(requestToken)
	fmt.Println(requestSecret)

	/*req, err := http.NewRequest("POST", "https://api.twitter.com/oauth/request_token", nil)

	if err != nil {
		log.Printf("Request Failed %s\n", err.Error())
	}
	req.Header.Set("oauth_consumer_key", os.Getenv("oauth_consumer_key"))
	req.Header.Set("oauth_callback", os.Getenv("callback_url"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Fehler beim Ausführen der Anfrage:", err)
		return
	}
	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		fmt.Println(error)
	}

	fmt.Println(string(body))*/

}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	requestToken, verifier, _ := oauth1.ParseAuthorizationCallback(r)

	fmt.Println("Request Token: ", requestToken)
	fmt.Println("verifier: ", verifier)

	accessToken, accessSecret, _ := config.AccessToken(requestToken, secret, verifier)
	// handle error
	token := oauth1.NewToken(accessToken, accessSecret)

	fmt.Println("Token: ", token)

	req, err := http.NewRequest("GET", "https://api.twitter.com/1.1/account/verify_credentials.json", nil)

	if err != nil {
		log.Printf("Request Failed %s\n", err.Error())
	}
	req.Header.Set("oauth_consumer_key", os.Getenv("oauth_consumer_key"))
	req.Header.Set("oauth_callback", os.Getenv("callback_url"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Fehler beim Ausführen der Anfrage:", err)
		return
	}
	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		fmt.Println(error)
	}

	fmt.Println(string(body))
}

var config oauth1.Config

func main() {
	godotenv.Load()

	config = oauth1.Config{
		ConsumerKey:    os.Getenv("oauth_consumer_key"),
		ConsumerSecret: os.Getenv("oauth_consumer_secret"),
		CallbackURL:    os.Getenv("callback_url"),
		Endpoint:       twitter.AuthorizeEndpoint,
	}

	router := chi.NewMux()

	router.Get("/", handleHome)
	router.Get("/login", handleLogin)
	router.Get("/oauth/twitter/callback", handleCallback)

	http.ListenAndServe(":3000", router)
}
