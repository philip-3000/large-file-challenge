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
For take 1, let's try deviating from the prompt a bit. We will still report
what it asks for, however, let's try not storing all the names/data, and rather
let's just look for what the prompt asks, and try to calculate on the fly.
*/
func Take_1(path string) {
	// start the clock.
	start := time.Now()

	readFile, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// storage for some calculations and data the prompt asks for.
	var lineCounter = 0
	var firstNameFrequencies = make(map[string]int)
	var mostCommonNameFrequency = 0
	var mostCommonName = ""
	var frequenciesByMonth = make(map[string]int)
	var namesAtIndices = make(map[int]string)
	var namesFound = false

	for fileScanner.Scan() {
		var line = fileScanner.Text()
		if len(line) > 0 {

			// for each line, process the row.
			var record, err = ProcessRow(line)
			if err == nil {
				// calculate this as we go rather than store all of those records.
				// first, the most common name calculations/record keeping.
				firstNameFrequencies[record.FirstName] += 1
				if firstNameFrequencies[record.FirstName] > mostCommonNameFrequency {
					mostCommonNameFrequency = firstNameFrequencies[record.FirstName]
					mostCommonName = record.FirstName
				}

				// now the months
				frequenciesByMonth[record.Month] += 1

				// we also have to take care of particular 'name at index' questions.
				// for this, we can just store them off.
				if !namesFound {
					if lineCounter == 0 {
						namesAtIndices[lineCounter] = fmt.Sprintf("%v %v", record.FirstName, record.LastName)
					} else if lineCounter == 432 {
						namesAtIndices[lineCounter] = fmt.Sprintf("%v %v", record.FirstName, record.LastName)
					} else if lineCounter == 43243 {
						namesAtIndices[lineCounter] = fmt.Sprintf("%v %v", record.FirstName, record.LastName)
					}
				}
			}
			lineCounter++
		}

	}

	// ok, we read through the entire file now, and counted the lines.
	// we also collected the names and donation date data.

	// report total # of lines
	fmt.Printf("Total Line Count: %v\n", lineCounter)

	if len(namesAtIndices) == 0 {
		fmt.Println("No records available.")
		return
	}

	// report names at indices
	var reporingIndices = []int{0, 432, 43243}
	for _, nameAt := range reporingIndices {
		val, ok := namesAtIndices[nameAt]
		if ok {
			fmt.Printf("Name at index %v: %v\n", nameAt, val)
		} else {
			fmt.Printf("Not enough data for position %v.\n", nameAt)
		}
	}

	fmt.Printf("***Time Elapsed so far: %v***\n", time.Since(start))

	// Now it wants us to report the donation frequency. I'm just going to do it
	// for the months.
	fmt.Println("Donation frequencies by month:")

	// there's not much in here...going to sort these so that the output is consistent
	keys := make([]string, 0, len(frequenciesByMonth))
	for k := range frequenciesByMonth {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("\t(Month, Frequency) = (%v, %v)\n", k, frequenciesByMonth[k])
	}

	fmt.Printf("Most Common First Name: %v. Frequency: %v\n", mostCommonName, mostCommonNameFrequency)
	fmt.Printf("***Time Elapsed so far: %v***\n", time.Since(start))

}
