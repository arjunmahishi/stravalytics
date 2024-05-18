package main

import (
	"strings"
	"testing"

	strava "github.com/strava/go.strava"
	"github.com/stretchr/testify/assert"
)

func TestInsertQueries(t *testing.T) {
	queries := insertActivityQueries(strava.ActivityDetailed{
		ActivitySummary: strava.ActivitySummary{
			Id:   1,
			Name: "Test",
		},
		BestEfforts: []*strava.BestEffort{
			{
				EffortSummary: strava.EffortSummary{
					Id:       1,
					Name:     "Test Best Effort",
					Distance: 1000,
				},
			},
		},
	})

	t.Logf(strings.Join(queries, "\n\n====================\n\n"))

	db, err := newDB()
	assert.NoError(t, err)

	for _, query := range queries {

		_, err := db.Exec(query)
		assert.NoError(t, err)
	}
}
