#!/bin/bash

echo "🧪 Performance Testing for pricofy-geocode-fr"
echo "============================================"

BINARY="./pricofy-geocode-fr"
LOG_FILE="performance.log"

# Clean log file
> "$LOG_FILE"

echo "Testing geocoding operations..."
echo ""

# Test 1: Geocode by postal code (Paris)
echo "1. Geocode by postal code (Paris 75001):"
start=$(date +%s%N)
result=$($BINARY '{"operation":"geocode-by-postal","postalCode":"75001"}' 2>>"$LOG_FILE")
end=$(date +%s%N)
duration=$(( (end - start) / 1000000 )) # Convert to milliseconds
echo "   Result: $(echo $result | jq -r '.municipality') ($(echo $result | jq -r '.province'))"
echo "   Latency: ${duration}ms"
echo ""

# Test 2: Geocode by municipality (Paris)
echo "2. Geocode by municipality (Paris):"
start=$(date +%s%N)
result=$($BINARY '{"operation":"geocode-by-municipality","municipality":"Paris"}' 2>>"$LOG_FILE")
end=$(date +%s%N)
duration=$(( (end - start) / 1000000 ))
echo "   Result: $(echo $result | jq -r '.postalCode') ($(echo $result | jq -r '.province'))"
echo "   Latency: ${duration}ms"
echo ""

# Test 3: Reverse geocode (Paris coordinates)
echo "3. Reverse geocode (Paris coordinates):"
start=$(date +%s%N)
result=$($BINARY '{"operation":"reverse-geocode","lat":48.8592,"lon":2.3417}' 2>>"$LOG_FILE")
end=$(date +%s%N)
duration=$(( (end - start) / 1000000 ))
echo "   Result: $(echo $result | jq -r '.city') $(echo $result | jq -r '.postalCode') ($(echo $result | jq -r '.province'))"
echo "   Distance: $(echo $result | jq -r '.distance') km"
echo "   Latency: ${duration}ms"
echo ""

# Test 4: Validate postal code
echo "4. Validate postal code (75001):"
start=$(date +%s%N)
result=$($BINARY '{"operation":"validate-postal","postalCode":"75001"}' 2>>"$LOG_FILE")
end=$(date +%s%N)
duration=$(( (end - start) / 1000000 ))
echo "   Valid: $(echo $result | jq -r '.valid')"
echo "   Latency: ${duration}ms"
echo ""

# Test 5: Validate municipality
echo "5. Validate municipality (Paris):"
start=$(date +%s%N)
result=$($BINARY '{"operation":"validate-municipality","municipality":"Paris"}' 2>>"$LOG_FILE")
end=$(date +%s%N)
duration=$(( (end - start) / 1000000 ))
echo "   Valid: $(echo $result | jq -r '.valid')"
echo "   Latency: ${duration}ms"
echo ""

# Test 6: Autocomplete postal code
echo "6. Autocomplete postal code (75):"
start=$(date +%s%N)
result=$($BINARY '{"operation":"autocomplete-postal","prefix":"75","limit":5}' 2>>"$LOG_FILE")
end=$(date +%s%N)
duration=$(( (end - start) / 1000000 ))
count=$(echo $result | jq '. | length')
echo "   Results: $count matches"
echo "   Latency: ${duration}ms"
echo ""

# Test 7: Autocomplete municipality
echo "7. Autocomplete municipality (Par):"
start=$(date +%s%N)
result=$($BINARY '{"operation":"autocomplete-municipality","query":"Par","limit":3}' 2>>"$LOG_FILE")
end=$(date +%s%N)
duration=$(( (end - start) / 1000000 ))
count=$(echo $result | jq '. | length')
echo "   Results: $count matches"
echo "   Latency: ${duration}ms"
echo ""

# Memory usage (rough estimate)
echo "8. Memory usage estimate:"
echo "   JSON file size: $(ls -lh internal/infrastructure/provider/postal-codes-fr.json | awk '{print $5}')"
echo "   Unique postal codes: $(grep -c '"[0-9]' internal/infrastructure/provider/postal-codes-fr.json)"
echo "   Estimated in-memory usage: ~3-4MB"
echo ""

echo "✅ Performance testing complete!"
echo "📊 Check $LOG_FILE for detailed logs"