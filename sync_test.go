package main

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	strava "github.com/strava/go.strava"
	"github.com/stretchr/testify/assert"
)

func TestGetExistingActivityLookup(t *testing.T) {
	t.Run("non existent file", func(t *testing.T) {
		*outputFile = "non_existent.json"

		lookup, all, err := getExistingActivity()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(lookup))
		assert.Equal(t, 0, len(all))
	})

	t.Run("empty file", func(t *testing.T) {
		*outputFile = "empty.json"

		activity := []*strava.ActivityDetailed{}
		raw, err := json.Marshal(activity)
		assert.NoError(t, err)

		cleanup, err := provisionTestFile(*outputFile, raw)
		assert.NoError(t, err)
		defer cleanup()

		lookup, all, err := getExistingActivity()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(lookup))
		assert.Equal(t, 0, len(all))
	})

	t.Run("non empty files", func(t *testing.T) {
		*outputFile = "non_empty.json"

		activities := []*strava.ActivityDetailed{
			{ActivitySummary: strava.ActivitySummary{Id: 1}},
			nil,
			{ActivitySummary: strava.ActivitySummary{Id: 2}},
			nil,
		}
		raw, err := json.Marshal(activities)
		assert.NoError(t, err)

		cleanup, err := provisionTestFile(*outputFile, raw)
		assert.NoError(t, err)
		defer cleanup()

		lookup, all, err := getExistingActivity()
		assert.NoError(t, err)

		assert.Equal(t, 2, len(lookup))
		assert.Equal(t, 2, len(all))
		assert.True(t, lookup[1])
		assert.True(t, lookup[2])
		assert.False(t, lookup[3])
	})
}

func provisionTestFile(filename string, data []byte) (func(), error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return nil, err
	}

	return func() {
		if err := os.Remove(filename); err != nil {
			log.Println("Error removing file", filename, err)
		}
	}, nil
}
