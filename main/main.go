package main

import (
	"fmt"
	"log"

	"github.com/kaidev1024/posm"
)

func main() {
	// SanitizeAddress test cases
	fmt.Println("=== Testing SanitizeAddress ===")
	sanitizeCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "collapse multiple spaces",
			input:    "123   Main   St",
			expected: "123 Main St",
		},
		{
			name:     "remove space before comma",
			input:    "San Francisco , CA",
			expected: "San Francisco, CA",
		},
		{
			name:     "collapse and remove space before comma",
			input:    "  615   John   Muir   Dr  ,   San   Francisco  ",
			expected: "615 John Muir Dr, San Francisco",
		},
		{
			name:     "tabs and newlines",
			input:    "Line1\n\tLine2  ,\tLine3",
			expected: "Line1 Line2, Line3",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}
	failed := 0
	for _, tc := range sanitizeCases {
		got := posm.SanitizeAddress(tc.input)
		if got != tc.expected {
			failed++
			fmt.Printf("[FAIL] %s\n  input:    %q\n  expected: %q\n  got:      %q\n", tc.name, tc.input, tc.expected, got)
		} else {
			fmt.Printf("[PASS] %s\n", tc.name)
		}
	}
	if failed > 0 {
		fmt.Printf("SanitizeAddress: %d failed\n\n", failed)
	} else {
		fmt.Println("SanitizeAddress: all passed\n")
	}

	// Initialize with your LocationIQ access token
	accessToken := "your_locationiq_access_token_here"
	posm.Init(accessToken)

	fmt.Println("=== Testing POSM APIs ===\n")

	// Test 1: GetStreetByText
	// fmt.Println("Test 1: GetStreetByText")
	// street, err := posm.GetStreetByText("Main St, San Francisco")
	// if err != nil {
	// 	log.Printf("Error: %v\n", err)
	// } else {
	// 	fmt.Printf("Street: %+v\n\n", street)
	// }

	// // Test 2: GetCityByText
	// fmt.Println("Test 2: GetCityByText")
	// city, err := posm.GetCityByText("San Francisco")
	// if err != nil {
	// 	log.Printf("Error: %v\n", err)
	// } else {
	// 	fmt.Printf("City: %+v\n\n", city)
	// }

	// // Test 3: GetPointsBySearch
	// fmt.Println("Test 3: GetPointsBySearch")
	// points, err := posm.GetPointsBySearch("615 John Muir Dr, San Francisco")
	// if err != nil {
	// 	log.Printf("Error: %v\n", err)
	// } else {
	// 	fmt.Printf("Found %d points:\n", len(points))
	// 	for i, p := range points {
	// 		fmt.Printf("  %d: %+v\n", i+1, p)
	// 	}
	// 	fmt.Println()
	// }

	// Test 4: GetCitiesBySearch
	fmt.Println("Test 4: GetCitiesBySearch")
	cities, err := posm.GetCitiesBySearch("xxxxxxxx")
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d cities:\n", len(cities))
		for i, c := range cities {
			fmt.Printf("  %d: %+v\n", i+1, c)
		}
		fmt.Println()
	}

	// Test 5: GetCitiesByAutocomplete
	// fmt.Println("Test 5: GetCitiesByAutocomplete")
	// autocompleteCities, err := posm.GetCitiesByAutocomplete("San")
	// if err != nil {
	// 	log.Printf("Error: %v\n", err)
	// } else {
	// 	fmt.Printf("Found %d autocomplete cities:\n", len(autocompleteCities))
	// 	for i, c := range autocompleteCities {
	// 		fmt.Printf("  %d: %+v\n", i+1, c)
	// 	}
	// 	fmt.Println()
	// }

	// Test 6: GetPointByTID (requires a valid TID from previous search)
	// fmt.Println("Test 6: GetPointByTID (example)")
	// if len(points) > 0 && points[0].OsmID != 0 {
	// 	tid := fmt.Sprintf("N%d", points[0].OsmID)
	// 	point, err := posm.GetPointByTID(tid)
	// 	if err != nil {
	// 		log.Printf("Error: %v\n", err)
	// 	} else {
	// 		fmt.Printf("Point by TID: %+v\n\n", point)
	// 	}
	// } else {
	// 	fmt.Println("Skipping - need valid OsmID from previous search\n")
	// }

	// Test 7: GetCityByTID (requires a valid TID from previous search)
	fmt.Println("Test 7: GetCityByTID (example)")
	if len(cities) > 0 && cities[0].OsmID != 0 {
		tid := fmt.Sprintf("R%d", cities[0].OsmID)
		cityByTID, err := posm.GetCityByTID(tid)
		if err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("City by TID: %+v\n\n", cityByTID)
		}
	} else {
		fmt.Println("Skipping - need valid OsmID from previous search\n")
	}

	fmt.Println("=== All tests completed ===")
}
