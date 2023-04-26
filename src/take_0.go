package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

/*
Questions:

- Write a program that will print out the total number of lines in the file.

- Notice that the 8th column contains a person’s name. Write a program that loads in this data and creates an array with all name strings. Print out the 432nd and 43243rd names.

- Notice that the 5th column contains a form of date. Count how many donations occurred in each month and print out the results.

- Notice that the 8th column contains a person’s name. Create an array with each first name. Identify the most common first name in the data and how many times it occurs.

*/

/*
For take 0, let's try following the prompt pretty closely.  We'll store data into arrays
rather than calculating on the fly.
*/
func Take_0(path string) {
	// start the clock.
	start := time.Now()

	readFile, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// as it requests, we need to store off certain data from each line in
	// an array.
	var lineCounter = 0
	var donationRecords = make([]DonationDataRecord, 0)

	// so let's do that.
	for fileScanner.Scan() {
		var line = fileScanner.Text()
		if len(line) > 0 {
			lineCounter++
			// for each line, process the row.
			var record, err = ProcessRow(line)
			if err == nil {
				// and it append the record to a list to be
				// processed later.
				donationRecords = append(donationRecords, record)
			}
		}

	}

	// ok, we read through the entire file now, and counted the lines.
	// we also collected the names and donation date data.

	// report total # of lines
	fmt.Printf("Total Line Count: %v\n", lineCounter)

	// report names at indices
	var nameAt = 0
	if len(donationRecords) == 0 {
		fmt.Println("No records available.")
		return
	}

	fmt.Printf("Name at index %v: %v %v\n", nameAt, donationRecords[nameAt].FirstName, donationRecords[nameAt].LastName)

	nameAt = 432
	if len(donationRecords) >= nameAt-1 {
		fmt.Printf("Name at index %v: %v %v\n", nameAt, donationRecords[nameAt].FirstName, donationRecords[nameAt].LastName)
	} else {
		fmt.Printf("Not enough data for position %v.\n", nameAt)
	}

	nameAt = 43243
	if len(donationRecords) >= nameAt-1 {
		fmt.Printf("Name at index %v: %v %v\n", nameAt, donationRecords[nameAt].FirstName, donationRecords[nameAt].LastName)
	} else {
		fmt.Printf("Not enough data for position %v.\n", nameAt)
	}

	fmt.Printf("***Time Elapsed so far: %v***\n", time.Since(start))

	// Now it wants us to report the donation frequency. I'm just going to do it
	// for the months.
	var frequenciesByMonth = make(map[string]int)

	// notice here we have to run through all that data again and process it.
	for _, record := range donationRecords {
		frequenciesByMonth[record.Month] += 1
	}

	// there's not much in here...going to sort these so that the output is consistent
	keys := make([]string, 0, len(frequenciesByMonth))
	for k := range frequenciesByMonth {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("Donation frequencies by month:")
	for _, k := range keys {
		fmt.Printf("\t(Month, Frequency) = (%v, %v)\n", k, frequenciesByMonth[k])
	}

	// Identify the most common first name in the data and how many times it occurs.
	// ok, so, again we have to go through the records here.
	var firstNameFrequencies = make(map[string]int)
	var mostCommonNameFrequency = 0
	var mostCommonName = ""
	for _, record := range donationRecords {
		firstNameFrequencies[record.FirstName] += 1

		currentVal := firstNameFrequencies[record.FirstName]
		if currentVal > mostCommonNameFrequency {
			mostCommonNameFrequency = currentVal
			mostCommonName = record.FirstName
		}
	}

	fmt.Printf("Most Common First Name: %v. Frequency: %v\n", mostCommonName, mostCommonNameFrequency)
	fmt.Printf("***Time Elapsed so far: %v***\n", time.Since(start))

}
