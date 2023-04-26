package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	versions := []Version{
		{Exec: Take_0, Info: "Processes the file sequentially and reports calculations and time stamps. Follows the prompt by creating arrays of names."},
		{Exec: Take_1, Info: "Processes the file sequentially and reports calculations and time stamps. Optimizes by not storing names."},
		{Exec: Take_2, Info: "Processes the file concurrently by reading in one go routine and processing in another over a buffered channel. Optimizes by not storing names."},
		{Exec: Take_3, Info: "Processes the file concurrently by reading in one go routine and sending over blocks of them to be processed over a buffered channel.Optimizes by not storing names."},
		{Exec: Take_4, Info: "Processes the file concurrently by reading in one go routine and sending over blocks of them to be processed by multiple go routines.Optimizes by not storing names."},
	}
	fileFlag := flag.String("input-file", "", "Specifies the data file to process.")
	listAllFlag := flag.Bool("list-versions", false, "Lists available versions to try.")
	versionFlag := flag.Int("version", -1, "The version of solution to the challenge to run.")
	calculateStatsModeFlag := flag.Bool("stats", false, "Calculate status mode. Runs the specified version n (default 5) number of times.")
	iterationsFlag := flag.Int("iterations", 5, "Specifies how many iterations to run when in stats mode. Defaults to 5.")

	flag.Parse()

	if *listAllFlag {
		fmt.Printf("Available versions:\n")
		for idx, version := range versions {
			fmt.Printf("\t%v: %v\n", idx, version.Info)
		}

		return
	}

	if *versionFlag == -1 {
		log.Fatal("The solution version must be specified.")
	}

	if *versionFlag < 0 || *versionFlag > (len(versions)-1) {
		log.Fatal("The solution version is invalid.")
	}

	if *fileFlag == "" {
		log.Fatal("Input path to file to process must be specified.")
	}

	// check to see if we are just running 1
	if !*calculateStatsModeFlag {
		// now run it.
		fmt.Printf("Running version %v...\n", *versionFlag)
		versions[*versionFlag].Exec(*fileFlag)
		return
	}

	if *iterationsFlag < 1 {
		log.Fatal("Number of iterations must be greater than 0.")
	}

	fmt.Printf("Running version '%v' %v times...\n\n", *versionFlag, *iterationsFlag)

	// ok, stats mode. need to do some calculations
	var runTimes = make([]float64, 0, *iterationsFlag)
	for i := 0; i < *iterationsFlag; i++ {
		start := time.Now()
		versions[*versionFlag].Exec(*fileFlag)
		duration := time.Since(start)

		runTimes = append(runTimes, duration.Seconds())
	}

	// calculate stats
	var stats = CalculateStats(runTimes)
	fmt.Printf("\nStatistics: \n")

	fmt.Printf("\tAverage Run Time: %v seconds\n", stats.Average)
	fmt.Printf("\tStandard Deviation: %v seconds\n", stats.StandardDeviaton)
	fmt.Printf("\tMax Run Time: %v seconds\n", stats.Max)
	fmt.Printf("\tMin Run Time: %v seconds\n", stats.Min)

}
