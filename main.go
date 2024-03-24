package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var postURL = "https://faxcaster.fly.dev/"

const (
	ButtonActionPost         = "post"
	ButtonActionPostRedirect = "post_redirect"
	ButtonActionLink         = "link"
	ButtonActionMint         = "mint"
	ButtonActionTx           = "tx"
)

type Button struct {
	Label   string
	Action  string
	Target  string
	PostURL string
}

type Frame struct {
	FormatText string
	Buttons    []Button
	Input      string
}

var frames = []Frame{
	// p0: landing page
	{
		FormatText: "Faxcaster!\n\nClick the button below to get started.",
		Buttons: []Button{
			{
				Label:  "Let us fax!",
				Action: ButtonActionPost,
			},
		},
	},
}

func main() {
	var bind string

	// for Fly.io
	port := os.Getenv("PORT")
	if port == "" {
		bind = "localhost:9080"
		postURL = "http://localhost:9080/"
	} else {
		bind = ":" + port
	}

	// landing page handler
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// generate the HTML for the landing page
		htmlContent, err := generateFrameHTML(postURL, frames[0])
		if err != nil {
			log.Println("failed to generate HTML:", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		// header(s) and output
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, htmlContent)
	})

	// POST handler
	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		// dump out the HTTP request, including headers
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		log.Println("Request received:\n" + string(dump))

		// parse body into DataRepresentation
		var data DataRepresentation
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Println("failed to decode JSON:", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		// decode message bytes
		message, err := decodeMessageBytes(data.TrustedData.MessageBytes)
		if err != nil {
			log.Println("failed to decode message bytes:", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		// marshal to JSON and log the message
		messageJSON, err := json.MarshalIndent(message, "", "  ")
		if err != nil {
			log.Println("failed to marshal message to JSON:", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		log.Printf("Message: %s", string(messageJSON))

		// generate the HTML for the required page
		htmlContent, err := generateFrameHTML(postURL, frames[0])
		if err != nil {
			log.Println("failed to generate HTML:", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		// header(s) and output
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, htmlContent)
	})

	log.Println("listening on", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}
