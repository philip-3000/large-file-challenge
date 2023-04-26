package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

type Row struct {
	lineNumber int
	text       string
}

/*
Questions:

- Write a program that will print out the total number of lines in the file.

- Notice that the 8th column contains a person’s name. Write a program that loads in this data and creates an array with all name strings. Print out the 432nd and 43243rd names.

- Notice that the 5th column contains a form of date. Count how many donations occurred in each month and print out the results.

- Notice that the 8th column contains a person’s name. Create an array with each first name. Identify the most common first name in the data and how many times it occurs.

*/

/*
For take 2, rather than process the file sequentially, can we utilize a channel
and a go routine to split up the work? that is, we can scan the lines
on one go routine, but, send them over a channel to a 'worker' go routine
to process them.
*/
func Take_2(path string) {

	// start the clock.
	start := time.Now()

	readFile, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var lineCounter = 0

	// we'll send our lines over to a processor go routine.
	// we can try playing with the buffer size - maybe there's
	// an optimal buffer size.
	var rows = make(chan Row, 64*1024)
	var wg sync.WaitGroup

	wg.Add(1)

	// reader.
	go func() {
		for fileScanner.Scan() {
			var line = fileScanner.Text()
			if len(line) > 0 {
				wg.Add(1)
				// send the row over our channel
				rows <- Row{lineNumber: lineCounter, text: line}

				// and count the lines.
				lineCounter++
			}

		}
		wg.Done()

	}()

	var firstNameFrequencies = make(map[string]int)
	var mostCommonNameFrequency = 0
	var mostCommonName = ""

	var frequenciesByMonth = make(map[string]int)

	var namesAtIndices = make(map[int]string)

	var namesFound = false

	// and we need a go routine to process it.
	//processed := 0
	go func() {
		for {
			select {
			case r, ok := <-rows:
				if ok {
					var record, err = ProcessRow(r.text)
					// if err != nil {
					// 	fmt.Println("WARNING: " + err.Error())
					// }
					if err == nil {
						//processed += 1
						// for the processed row, we need to calculate the first name
						// frequency (increment it) and do a max check.
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
							if r.lineNumber == 0 {
								namesAtIndices[r.lineNumber] = fmt.Sprintf("%v %v", record.FirstName, record.LastName)
							} else if r.lineNumber == 432 {
								namesAtIndices[r.lineNumber] = fmt.Sprintf("%v %v", record.FirstName, record.LastName)
							} else if r.lineNumber == 43243 {
								namesAtIndices[r.lineNumber] = fmt.Sprintf("%v %v", record.FirstName, record.LastName)
								namesFound = true
							}
						}
					}
					wg.Done()
				}

			}
		}
	}()

	// wait for scanning and processing to complete.
	wg.Wait()

	// c1 - report total # of lines
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

	fmt.Println("Donation frequencies by month:")
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
