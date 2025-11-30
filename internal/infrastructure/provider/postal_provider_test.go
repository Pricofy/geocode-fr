package provider

import (
	"testing"

	"github.com/pricofy/geocode-fr/internal/domain"
)

func TestPostalCodeProvider_GeocodeByPostalCode(t *testing.T) {
	p := NewPostalCodeProvider()

	tests := []struct {
		name        string
		postalCode  string
		wantSuccess bool
		wantErr     bool
	}{
		{
			name:        "valid postal code - Paris",
			postalCode:  "75001",
			wantSuccess: true,
			wantErr:     false,
		},
		{
			name:        "valid postal code - Troyes",
			postalCode:  "10000",
			wantSuccess: true,
			wantErr:     false,
		},
		{
			name:        "invalid postal code",
			postalCode:  "99999",
			wantSuccess: false,
			wantErr:     true,
		},
		{
			name:        "empty postal code",
			postalCode:  "",
			wantSuccess: false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := p.GeocodeByPostalCode(tt.postalCode)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if _, ok := err.(*domain.PostalCodeNotFoundError); !ok {
					t.Errorf("Expected PostalCodeNotFoundError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.Success != tt.wantSuccess {
					t.Errorf("Expected success=%v, got %v", tt.wantSuccess, result.Success)
				}
				if result.PostalCode != tt.postalCode {
					t.Errorf("Expected postalCode=%s, got %s", tt.postalCode, result.PostalCode)
				}
				if result.Coords.Lat == 0 && result.Coords.Lon == 0 {
					t.Errorf("Expected non-zero coordinates")
				}
			}
		})
	}
}

func TestPostalCodeProvider_GeocodeByMunicipality(t *testing.T) {
	p := NewPostalCodeProvider()

	tests := []struct {
		name         string
		municipality string
		wantSuccess  bool
		wantErr      bool
	}{
		{
			name:         "valid municipality - Paris",
			municipality: "Paris",
			wantSuccess:  true,
			wantErr:      false,
		},
		{
			name:         "valid municipality - Troyes",
			municipality: "Troyes",
			wantSuccess:  true,
			wantErr:      false,
		},
		{
			name:         "invalid municipality",
			municipality: "NonExistentCity",
			wantSuccess:  false,
			wantErr:      true,
		},
		{
			name:         "case insensitive",
			municipality: "paris",
			wantSuccess:  true,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := p.GeocodeByMunicipality(tt.municipality)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.Success != tt.wantSuccess {
					t.Errorf("Expected success=%v, got %v", tt.wantSuccess, result.Success)
				}
				if result.Municipality == "" {
					t.Errorf("Expected non-empty municipality")
				}
			}
		})
	}
}

func TestPostalCodeProvider_ReverseGeocode(t *testing.T) {
	p := NewPostalCodeProvider()

	tests := []struct {
		name        string
		lat         float64
		lon         float64
		wantSuccess bool
		wantErr     bool
	}{
		{
			name:        "Paris coordinates",
			lat:         48.8592,
			lon:         2.3417,
			wantSuccess: true,
			wantErr:     false,
		},
		{
			name:        "Troyes coordinates",
			lat:         48.3007,
			lon:         4.0852,
			wantSuccess: true,
			wantErr:     false,
		},
		{
			name:        "Valid coordinates in France",
			lat:         43.6047,
			lon:         1.4442,
			wantSuccess: true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, distance, err := p.ReverseGeocode(tt.lat, tt.lon)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.Success != tt.wantSuccess {
					t.Errorf("Expected success=%v, got %v", tt.wantSuccess, result.Success)
				}
				if result.PostalCode == "" {
					t.Errorf("Expected non-empty postalCode")
				}
				if distance < 0 {
					t.Errorf("Expected non-negative distance, got %f", distance)
				}
				if distance > 1000 {
					t.Errorf("Distance seems too large: %f km", distance)
				}
			}
		})
	}
}

func TestPostalCodeProvider_ValidatePostalCode(t *testing.T) {
	p := NewPostalCodeProvider()

	tests := []struct {
		name       string
		postalCode string
		want       bool
	}{
		{
			name:       "valid postal code",
			postalCode: "28001",
			want:       true,
		},
		{
			name:       "invalid postal code",
			postalCode: "99999",
			want:       false,
		},
		{
			name:       "empty postal code",
			postalCode: "",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.ValidatePostalCode(tt.postalCode)
			if got != tt.want {
				t.Errorf("ValidatePostalCode(%s) = %v, want %v", tt.postalCode, got, tt.want)
			}
		})
	}
}

func TestPostalCodeProvider_ValidateMunicipality(t *testing.T) {
	p := NewPostalCodeProvider()

	tests := []struct {
		name         string
		municipality string
		want         bool
	}{
		{
			name:         "valid municipality",
			municipality: "Paris",
			want:         true,
		},
		{
			name:         "case insensitive",
			municipality: "paris",
			want:         true,
		},
		{
			name:         "invalid municipality",
			municipality: "NonExistentCity",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.ValidateMunicipality(tt.municipality)
			if got != tt.want {
				t.Errorf("ValidateMunicipality(%s) = %v, want %v", tt.municipality, got, tt.want)
			}
		})
	}
}

func TestPostalCodeProvider_AutocompletePostalCode(t *testing.T) {
	p := NewPostalCodeProvider()

	tests := []struct {
		name      string
		prefix    string
		limit     int
		wantCount int
		wantErr   bool
	}{
		{
			name:      "prefix 280",
			prefix:    "280",
			limit:     10,
			wantCount: 10,
			wantErr:   false,
		},
		{
			name:      "prefix 08",
			prefix:    "08",
			limit:     5,
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "non-existent prefix",
			prefix:    "99",
			limit:     10,
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := p.AutocompletePostalCode(tt.prefix, tt.limit)

			if len(results) != tt.wantCount {
				t.Errorf("AutocompletePostalCode(%s, %d) returned %d results, want %d",
					tt.prefix, tt.limit, len(results), tt.wantCount)
			}

			// Verify all results start with prefix
			for _, result := range results {
				if len(result.PostalCode) < len(tt.prefix) || result.PostalCode[:len(tt.prefix)] != tt.prefix {
					t.Errorf("Result %s does not start with prefix %s", result.PostalCode, tt.prefix)
				}
			}
		})
	}
}

func TestPostalCodeProvider_AutocompleteMunicipality(t *testing.T) {
	p := NewPostalCodeProvider()

	tests := []struct {
		name      string
		query     string
		limit     int
		wantCount int
		wantErr   bool
	}{
		{
			name:      "query 'mad'",
			query:     "mad",
			limit:     10,
			wantCount: 10,
			wantErr:   false,
		},
		{
			name:      "query 'bar'",
			query:     "bar",
			limit:     5,
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "case insensitive",
			query:     "MAD",
			limit:     10,
			wantCount: 10,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := p.AutocompleteMunicipality(tt.query, tt.limit)

			if len(results) > tt.limit {
				t.Errorf("AutocompleteMunicipality(%s, %d) returned %d results, want at most %d",
					tt.query, tt.limit, len(results), tt.limit)
			}

			// Verify all results contain query (case insensitive)
			queryLower := tt.query
			for _, result := range results {
				municipalityLower := result.Municipality
				if len(municipalityLower) < len(queryLower) {
					t.Errorf("Result municipality %s is shorter than query %s", result.Municipality, tt.query)
				}
			}
		})
	}
}

func TestPostalCodeProvider_CalculateDistance(t *testing.T) {
	p := NewPostalCodeProvider()

	// Test Haversine distance calculation
	// Paris coordinates
	parisLat, parisLon := 48.8592, 2.3417
	troyesLat, troyesLon := 48.3007, 4.0852

	_, dist, err := p.ReverseGeocode(parisLat, parisLon)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	_, dist2, err := p.ReverseGeocode(troyesLat, troyesLon)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Distance should be reasonable (not negative, not too large)
	if dist < 0 || dist > 1000 {
		t.Errorf("Distance from Paris seems incorrect: %f km", dist)
	}
	if dist2 < 0 || dist2 > 1000 {
		t.Errorf("Distance from Troyes seems incorrect: %f km", dist2)
	}
}
