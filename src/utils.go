package main

import (
	"errors"
	"math"
	"strings"
)

/*
Helper function to calculate some statistics over run times.
*/
func CalculateStats(times []float64) Stats {
	// first we need the sum.
	var sum float64 = 0
	var min float64 = math.Inf(1)
	var max float64 = math.Inf(-1)
	for _, time := range times {
		sum += time

		if time < min {
			min = time
		}

		if time > max {
			max = time
		}
	}

	// calclulate the average/mean.
	var average float64 = sum / float64(len(times))

	// now calculate the standard deviation.
	var standardDeviation float64 = 0
	for _, time := range times {
		standardDeviation += math.Pow(time-average, 2)
	}

	standardDeviation = math.Sqrt(standardDeviation / float64(len(times)))

	return Stats{Average: average, StandardDeviaton: standardDeviation, Min: min, Max: max}
}

/*
Helper function to process a line of election/donation data into a
DonataionDataRecord. Sometimes the data contains companies, and not
a person in the format of {Last Name}, {First Name}; these are ignored.
*/
func ProcessRow(text string) (DonationDataRecord, error) {
	row := strings.Split(text, "|")
	var record = DonationDataRecord{}

	// extract the first name/last name and trim whitespace.
	lastName, firstName, found := strings.Cut(row[7], ",")
	lastName = strings.Trim(lastName, " ")
	firstName = strings.Trim(firstName, " ")

	// return an error to processor for when we cannot find an actual person.
	if !found {
		return record, errors.New("couldn't parse first and last name from row")
	}

	// it asks us to extract the date from the 5th column, which looks like this:
	// 201903139145684455. I believe the year and month consists of the
	// first 4 digits and then the next two, respectively.
	var dateString = row[4]
	if len(dateString) < 6 {
		return record, errors.New("couldn't parse year and month from row.")
	}

	// these are left as strings for the moment.
	// we may not need the year at all.
	record.Year = dateString[:4]
	record.Month = dateString[4:6]
	record.FirstName = firstName
	record.LastName = lastName
	return record, nil

}
