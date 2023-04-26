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

/*
Questions:

- Write a program that will print out the total number of lines in the file.

- Notice that the 8th column contains a person’s name. Write a program that loads in this data and creates an array with all name strings. Print out the 432nd and 43243rd names.

- Notice that the 5th column contains a form of date. Count how many donations occurred in each month and print out the results.

- Notice that the 8th column contains a person’s name. Create an array with each first name. Identify the most common first name in the data and how many times it occurs.

*/

/*
For take4, is it possible to decrease the processing time by increasing
the number of workers we have by utilizing multiple go routines?
*/
func Take_4(path string) {

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

	var numberOfWorkers = 3
	// we'll send our lines over to a processor go routine in batches
	// of size 64K, with a small buffer.
	var rows = make(chan []Row, numberOfWorkers)
	var wg sync.WaitGroup

	wg.Add(1)

	// reader.
	var blockSize = 64 * 1024
	lineBuffer := make([]Row, 0, blockSize)
	go func() {
		for fileScanner.Scan() {
			var line = fileScanner.Text()
			if len(line) > 0 {

				// append the row to our buffer.
				lineBuffer = append(lineBuffer, Row{lineNumber: lineCounter, text: line})

				// and count the lines.
				lineCounter++

				// if our buffer size has gotten to our block size
				if len(lineBuffer) == blockSize {
					wg.Add(1)
					// send it over to our processor.
					rows <- lineBuffer

					// and reinitialize the buffer.
					lineBuffer = make([]Row, 0, blockSize)
				}
			}

		}

		// when done scanning, need to check buffer. it is unlkely it has
		// stopped right after the last 64K block.
		if len(lineBuffer) > 0 {
			wg.Add(1)
			rows <- lineBuffer
		}

		// now the reader is done.
		wg.Done()
	}()

	// some shared resources between go routines that will be
	// utilized to store calculations/data requested by the prompt.
	var firstNameFrequencies = make(map[string]int)
	var mostCommonNameFrequency = 0
	var mostCommonName = ""

	var frequenciesByMonth = make(map[string]int)
	var namesAtIndices = make(map[int]string)
	var namesFound = false

	// and so need to protect some of this stuff with a mutex.
	mutex := &sync.Mutex{}

	// need go routine workers to process it.
	for i := 0; i < numberOfWorkers; i++ {
		go func() {
			for {
				select {
				case bufferedRow, ok := <-rows:
					if ok {
						for _, r := range bufferedRow {
							var record, err = ProcessRow(r.text)
							if err == nil {
								//processed += 1
								// for the processed row, we need to calculate the first name
								// frequency (increment it) and do a max check.
								mutex.Lock()
								firstNameFrequencies[record.FirstName] += 1
								if firstNameFrequencies[record.FirstName] > mostCommonNameFrequency {
									mostCommonNameFrequency = firstNameFrequencies[record.FirstName]
									mostCommonName = record.FirstName
								}

								// now the months
								frequenciesByMonth[record.Month] += 1
								mutex.Unlock()

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

						}
						wg.Done()
					}

				}
			}
		}()
	}

	// wait for scanning and processing to complete.
	wg.Wait()

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

	// now it wants us to report the donation frequency. I'm just going to do it
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
