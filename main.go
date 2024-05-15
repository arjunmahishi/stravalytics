package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	strava "github.com/strava/go.strava"
)

var (
	authenticator *strava.OAuthAuthenticator
	server        *http.Server

	port       = flag.Int("port", 3000, "Port for the local server")
	workers    = flag.Int("concurrency", runtime.NumCPU(), "Number of concurrent workers")
	outputFile = flag.String("output", "activities.json", "Output file for activities")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU()) // setup to use all the cores

	clientID, err := strconv.Atoi(os.Getenv("STRAVA_CLIENT_ID"))
	if err != nil {
		log.Fatal("Invalid STRAVA_CLIENT_ID")
	}

	strava.ClientId = clientID
	strava.ClientSecret = os.Getenv("STRAVA_CLIENT_SECRET")

	// define a strava.OAuthAuthenticator to hold state.
	// The callback url is used to generate an AuthorizationURL.
	// The RequestClientGenerator can be used to generate an http.RequestClient.
	// This is usually when running on the Google App Engine platform.
	authenticator = &strava.OAuthAuthenticator{
		CallbackURL:            fmt.Sprintf("http://localhost:%d/exchange_token", *port),
		RequestClientGenerator: nil,
	}

	path, err := authenticator.CallbackPath()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc(path, authenticator.HandlerFunc(oAuthSuccess, oAuthFailure))

	// login to strava
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, authenticator.AuthorizationURL("state1", "activity:read", false), http.StatusFound)
	})

	// start the server
	fmt.Printf("Visit %s to authenticate\n", fmt.Sprintf("http://localhost:%d/login", *port))
	server = &http.Server{Addr: fmt.Sprintf(":%d", *port)}
	log.Fatal(server.ListenAndServe())
}

func oAuthSuccess(auth *strava.AuthorizationResponse, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Successfully Authenticated\n\n")
	fmt.Fprintf(w, "State: %s\n\n", auth.State)

	fmt.Fprintf(w, "The Authenticated Athlete (you):\n")
	content, _ := json.MarshalIndent(auth.Athlete, "", " ")
	fmt.Fprint(w, string(content), "\n\n")

	go syncActivities(auth)
}

func oAuthFailure(err error, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Authorization Failure:\n")

	// some standard error checking
	if err == strava.OAuthAuthorizationDeniedErr {
		fmt.Fprint(w, "The user clicked the 'Do not Authorize' button on the previous page.\n")
		fmt.Fprint(w, "This is the main error your application should handle.")
	} else if err == strava.OAuthInvalidCredentialsErr {
		fmt.Fprint(w, "You provided an incorrect client_id or client_secret.\nDid you remember to set them at the begininng of this file?")
	} else if err == strava.OAuthInvalidCodeErr {
		fmt.Fprint(w, "The temporary token was not recognized, this shouldn't happen normally")
	} else if err == strava.OAuthServerErr {
		fmt.Fprint(w, "There was some sort of server error, try again to see if the problem continues")
	} else {
		fmt.Fprint(w, err)
	}
}
