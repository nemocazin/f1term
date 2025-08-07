package fetchData

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FetchByYears get all the data for every season from now to 1950
// If there is no data for the season, we increment a counter
// If for 3 consecutive years, we have no data, it means we finish fetching
func FetchByYears(done chan struct{}) error {
	currentYear := time.Now().Year()
	consecutiveEmpty := 0

	for year := currentYear; year >= 1950; year-- {
		url := fmt.Sprintf("https://api.openf1.org/v1/meetings?year=%d", year)

		// Connect to OpenF1 API
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("cannot Access API")
		}

		// Read data for the year
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("could not read for %d", year)
		}

		// Tranform it into readble data
		var result []any
		err = json.Unmarshal(body, &result)
		// Check if data is empty
		if err != nil || len(result) == 0 {
			consecutiveEmpty++
			if consecutiveEmpty >= 3 {
				break
			}
			continue
		}

		// Reset counter
		consecutiveEmpty = 0
	}

	done <- struct{}{}
	return nil
}
