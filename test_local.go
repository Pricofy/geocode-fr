package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pricofy/geocode-fr/internal/application"
	"github.com/pricofy/geocode-fr/internal/domain"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_local.go <operation> [args...]")
		fmt.Println("Example: go run test_local.go geocode-by-postal 75001")
		os.Exit(1)
	}

	// Initialize service
	service := application.NewPostalCodeService()

	operation := os.Args[1]

	switch operation {
	case "geocode-by-postal":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run test_local.go geocode-by-postal <postal_code>")
			os.Exit(1)
		}
		postalCode := os.Args[2]
		event := domain.LambdaEvent{
			Body: fmt.Sprintf(`{"operation":"geocode-by-postal","postalCode":"%s"}`, postalCode),
		}
		result, err := service.GeocodeByPostal(event)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))

	case "geocode-by-municipality":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run test_local.go geocode-by-municipality <municipality>")
			os.Exit(1)
		}
		municipality := os.Args[2]
		event := domain.LambdaEvent{
			Body: fmt.Sprintf(`{"operation":"geocode-by-municipality","municipality":"%s"}`, municipality),
		}
		result, err := service.GeocodeByPostal(event) // This will try postal first, then municipality
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))

	case "reverse-geocode":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run test_local.go reverse-geocode <lat> <lon>")
			os.Exit(1)
		}
		lat := os.Args[2]
		lon := os.Args[3]
		event := domain.LambdaEvent{
			Body: fmt.Sprintf(`{"operation":"reverse-geocode","lat":%s,"lon":%s}`, lat, lon),
		}
		result, err := service.ReverseGeocode(event)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))

	case "validate-postal":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run test_local.go validate-postal <postal_code>")
			os.Exit(1)
		}
		postalCode := os.Args[2]
		event := domain.LambdaEvent{
			Body: fmt.Sprintf(`{"operation":"validate-postal","postalCode":"%s"}`, postalCode),
		}
		result, err := service.ValidatePostal(event)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))

	default:
		fmt.Printf("Unknown operation: %s\n", operation)
		fmt.Println("Supported operations: geocode-by-postal, geocode-by-municipality, reverse-geocode, validate-postal")
		os.Exit(1)
	}
}

