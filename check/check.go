package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pivotal-cf-experimental/cron-resource/models"
)

var daysOfTheWeek = [...]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

func getVersions(request models.CheckRequest) (models.CheckResponse, error) {

	loc, err := time.LoadLocation(request.Source.Location)
	if err != nil {
		return nil, err
	}

	dayToFire := request.Source.DayToFire
	hourToFire := request.Source.HourToFire

	err = validateWhenToFire(dayToFire, hourToFire)
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)
	previouslyFiredAt := request.Version.Time

	if !timeWithinRange(now, dayToFire, hourToFire) {
		// TODO: Think about this
		// if previouslyFiredAt.IsZero() {
		// 	return append([]models.Version{}, models.Version{
		// 		Time: now.Add(-72 * time.Hour),
		// 	}), nil
		// }
		return []models.Version{}, nil
	}

	// We are within the firing range
	if previouslyFiredAt.IsZero() {
		return append([]models.Version{}, models.Version{
			Time: now,
		}), nil
	}

	if !previouslyFiredAt.After(now.Add(-2 * time.Hour)) {
		return append([]models.Version{}, models.Version{
			Time: now,
		}), nil
	}

	return []models.Version{}, nil
}

func validateWhenToFire(dayToFire string, hourToFire int) error {

	if !isADayOfTheWeek(dayToFire) {
		return fmt.Errorf(`"day_to_fire" should be one of the following: %q`, strings.Join(daysOfTheWeek[:], ","))
	}

	if hourToFire < 0 || hourToFire > 23 {
		return errors.New(`"hour_to_fire" should be in the 0-23 range`)
	}

	return nil
}

func isADayOfTheWeek(day string) bool {
	sanitizedInput := strings.TrimSpace(strings.ToLower(day))

	for _, dayOfTheWeek := range daysOfTheWeek {
		if sanitizedInput == strings.ToLower(dayOfTheWeek) {
			return true
		}
	}
	return false
}

func timeWithinRange(now time.Time, dayToFire string, hourToFire int) bool {

	if strings.ToLower(now.Weekday().String()) == strings.ToLower(dayToFire) && now.Hour() == hourToFire {
		return true
	}

	return false
}

func main() {
	var request models.CheckRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error decoding payload: "+err.Error())
		os.Exit(1)
	}

	versions, err := getVersions(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid resource config: "+err.Error())
		os.Exit(1)
	}

	json.NewEncoder(os.Stdout).Encode(versions)
}
