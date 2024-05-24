package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	strava "github.com/strava/go.strava"
)

var (
	authenticator *strava.OAuthAuthenticator
	server        *http.Server

	port     = flag.Int("port", 8000, "Port for the local server")
	workers  = flag.Int("max-concurrency", runtime.NumCPU(), "Number of concurrent workers")
	dataFile = flag.String("data", "activities.json", "File to store activity data")
	dbHost   = flag.String("db-host", "localhost", "ClickHouse host")
	dbPort   = flag.Int("db-port", 9000, "ClickHouse port")
	sync     = flag.Bool("sync", false, "Run only the sync process")
	load     = flag.Bool("load", false, "Only load the already synced data into the db")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	syncLoadXNOR := (*sync && *load) || (!*sync && !*load)

	if *sync || syncLoadXNOR {
		if err := runSync(); err != nil {
			log.Fatal(err)
		}
	}

	if *load || syncLoadXNOR {
		db, err := newDB()
		if err != nil {
			log.Fatal(err)
		}

		loadData(db)
	}
}

func loadData(db driver.Conn) {
	f, err := os.Open(*dataFile)
	if err != nil {
		log.Fatal(err)
	}

	raw, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	activities := []*strava.ActivityDetailed{}
	if err = json.Unmarshal(raw, &activities); err != nil {
		log.Fatal(err)
	}

	allQueries := []string{}
	for _, activity := range activities {
		if activity == nil {
			continue
		}

		allQueries = append(allQueries, insertActivityQueries(activity)...)
	}

	fmt.Print("inserting data into the db")
	bulkInsertData(db, allQueries)
	fmt.Print("...DONE\n")
}

func runSync() error {
	clientID, err := strconv.Atoi(os.Getenv("STRAVA_CLIENT_ID"))
	if err != nil {
		return err
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
		return err
	}

	http.HandleFunc(path, authenticator.HandlerFunc(oAuthSuccess, oAuthFailure))

	// login to strava
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, authenticator.AuthorizationURL("state1", "activity:read", false), http.StatusFound)
	})

	// start the server
	fmt.Printf("Visit %s to authenticate\n", fmt.Sprintf("http://localhost:%d/login", *port))
	server = &http.Server{Addr: fmt.Sprintf(":%d", *port)}
	log.Println(server.ListenAndServe())
	return nil
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
