package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	strava "github.com/strava/go.strava"
)

type activityDetailsFunc func() (*strava.ActivityDetailed, error)

func detailWorkers(
	id int, jobs <-chan activityDetailsFunc, results chan<- *strava.ActivityDetailed,
) {
	for do := range jobs {
		activity, err := do()
		if err != nil {
			results <- nil
			log.Printf("[%d] Error: %v", id, err)
			continue
		}

		results <- activity
		log.Printf(
			"[%d] fetched details for '%s (%s)'",
			id, activity.Name, activity.StartDateLocal.Format("2006-01-02"),
		)
	}
}

func syncActivities(auth *strava.AuthorizationResponse) {
	// create a new strava client using the authorized token
	client := strava.NewClient(auth.AccessToken)
	req := strava.NewAthletesService(client).ListActivities(auth.Athlete.Id)
	activities := []*strava.ActivitySummary{}
	page := 1

	existingActivities, allActivities, err := getExistingActivity()
	if err != nil {
		log.Fatal(err)
	}

	for {
		currPage, err := req.PerPage(200).Page(page).Do()
		if err != nil {
			log.Fatal(err)
		}

		for _, activity := range currPage {
			if _, exists := existingActivities[activity.Id]; !exists {
				activities = append(activities, activity)
			}
		}

		log.Printf("Page: %d, Activities: %d", page, len(activities))
		if len(currPage) == 0 {
			break
		}

		page++
	}

	f, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	newActivities := getActivityDetails(client, activities)
	if len(newActivities) != 0 {
		detailBytes, err := json.Marshal(append(allActivities, newActivities...))
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.Write(detailBytes)
		if err != nil {
			log.Fatal(err)
		}
	}

	server.Shutdown(context.Background())
}

func getActivityDetails(
	client *strava.Client, activities []*strava.ActivitySummary,
) []*strava.ActivityDetailed {
	log.Println(len(activities), " activities")
	jobsChan := make(chan activityDetailsFunc, len(activities))
	resultsChan := make(chan *strava.ActivityDetailed, len(activities))
	allDetails := []*strava.ActivityDetailed{}

	// start workers
	for i := 0; i < *workers; i++ {
		i := i
		go detailWorkers(i, jobsChan, resultsChan)
	}

	for _, activity := range activities {
		activity := activity
		jobsChan <- func() (*strava.ActivityDetailed, error) {

			activity, err := strava.NewActivitiesService(client).Get(activity.Id).Do()
			if err != nil {
				return nil, err
			}

			return activity, nil
		}
	}

	for i := 0; i < len(activities); i++ {
		allDetails = append(allDetails, <-resultsChan)
	}
	close(jobsChan)
	close(resultsChan)

	log.Println("All activities fetched!")
	return allDetails
}

func getExistingActivity() (map[int64]bool, []*strava.ActivityDetailed, error) {
	f, err := os.Open(*outputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}

		return nil, nil, err
	}
	defer f.Close()

	existingActivities := []*strava.ActivityDetailed{}
	raw, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}

	if len(raw) == 0 {
		return nil, nil, nil
	}

	if err = json.Unmarshal(raw, &existingActivities); err != nil {
		return nil, nil, err
	}

	lookup := make(map[int64]bool, len(existingActivities))
	cleanList := []*strava.ActivityDetailed{}
	for _, activity := range existingActivities {
		if activity == nil {
			continue
		}

		lookup[activity.Id] = true
		cleanList = append(cleanList, activity)
	}

	return lookup, cleanList, nil
}
