package main

import (
	"fmt"

	strava "github.com/strava/go.strava"
)

const (
	insertActivity = `INSERT into strava.activities (
id, name, distance, moving_time, elapsed_time, total_elevation_gain, type, start_date, start_latitude, start_longitude, end_latitude, end_longitude, average_speed, max_speed, average_cadence, kilojoules, average_heartrate, max_heartrate, has_kudoed)
VALUES (%v, '%v', %v, %v, %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', %v, %v, %v, %v, %v, %v, %v)`

	insertBestEffort = `INSERT into strava.best_efforts (
id, activity_id, name, distance, moving_time, elapsed_time, start_date, pr_rank)
VALUES (%v, %v, '%v', %v, %v, %v, '%v', %v)`
)

func insertActivityQueries(activity *strava.ActivityDetailed) []string {
	queries := []string{}
	startDate := activity.StartDate.Format("2006-01-02 15:04:05")

	queries = append(queries, fmt.Sprintf(
		insertActivity,
		activity.Id, activity.Name, activity.Distance, activity.MovingTime,
		activity.ElapsedTime, activity.TotalElevationGain, activity.Type,
		startDate, activity.StartLocation[0], activity.StartLocation[1],
		activity.EndLocation[0], activity.EndLocation[1], activity.AverageSpeed,
		activity.MaximunSpeed, activity.AverageCadence, activity.Kilojoules,
		activity.AverageHeartrate, activity.MaximumHeartrate, activity.HasKudoed,
	))

	for _, bestEffort := range activity.BestEfforts {
		startDate := bestEffort.StartDate.Format("2006-01-02 15:04:05")
		queries = append(queries, fmt.Sprintf(
			insertBestEffort,
			bestEffort.Id, activity.Id, bestEffort.Name, bestEffort.Distance, bestEffort.MovingTime, bestEffort.ElapsedTime, startDate, bestEffort.PRRank,
		))
	}

	return queries
}
