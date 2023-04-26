package main

// Struct for storing relevant information from a line of election data.
type DonationDataRecord struct {
	FirstName string
	LastName  string
	Month     string
	Year      string
}

type Version struct {
	Exec func(path string)
	Info string
}

type Stats struct {
	Average          float64
	StandardDeviaton float64
	Min              float64
	Max              float64
}
