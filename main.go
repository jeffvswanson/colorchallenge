package main

import (
	"bufio"
	"fmt"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/jeffvswanson/colorchallenge/pkg/exporttocsv"

	log "github.com/jeffvswanson/colorchallenge/pkg/errorlogging"
)

type colorCode struct {
	Red, Green, Blue int
}

type logInfo struct {
	Level, Message string
	ErrorMessage   error
}

func init() {
	log.FormatLog()
}

func main() {

	// Setup
	status := logInfo{
		Level:   "Info",
		Message: "Beginning setup.",
	}
	log.WriteToLog(status.Level, status.Message, nil)

	// CSV file setup
	status = logInfo{
		Level:   "Info",
		Message: csvSetup("ColorChallengeOutput"),
	}
	log.WriteToLog(status.Level, status.Message, nil)

	imgColorPrevalence, message := extractURLs("input.txt")
	status = logInfo{
		Level:   "Info",
		Message: message,
	}
	log.WriteToLog(status.Level, status.Message, nil)

	var keyCount int
	for key := range imgColorPrevalence {
		keyCount++
		resp, err := http.Get(key)
		if log.ErrorCheck("Warn", fmt.Sprintf("http.Get(%v) failure", key), err) {
			continue
		}
		defer resp.Body.Close()
		imgData, err := jpeg.DecodeConfig(resp.Body)
		if log.ErrorCheck("Warn", fmt.Sprintf("%v image decode error", key), err) {
			continue
		}
		fmt.Printf("URL %d: %v\n\tH: %d, W: %d\n", keyCount, key, imgData.Height, imgData.Width)
		// xDim := imgData.Width
		// yDim := imgData.Height
	}
}

// Start with the end in mind.

// Result written to CSV file in the form of url, top_color1, top_color2, top_color3. O(n) Key = url, value = string of top 3 hexadecimal values

// Convert RGB color scheme (0 - 255, 0 - 255, 0 - 255) (256 bytes or 2^8) to hexadecimal format (#000000 - #FFFFFF) O(1) due to only needing
// to deal with 3 colors.

// Utilize quicksort to sort colors into ascending/descending order and slice off top 3. O(n lg n)

// 1st approach to get colors from image
// Scan image pixel by pixel and increment a counter relating to each color found in the image. Best way in a map. Key = color code,
// value = number of times color found.

// Navigate to the image

// Keep a counter of what image we're on

// Allocate a map to hold the url and it's index of colors. As the algorithm progresses the slice holding the colors will be converted to

// Load in the input.txt file line by line and send off a gofunc for as many lines are as possible while staying within memory and CPU constraints.

// Constant to set max memory used. 512 MB. I'm guessing their using a Docker container or something similar.

// Setup function to initialize log file and csv file to write to.

/*
Ideas:

1. Make sections of the code supporting packages. For example, not all the code needs to be in one main file. The CSV handler could be a package and called into main.

2. Nothing says I'm explicitly limited, just that I may be limited.

3. Given list 1000 urls to an image, to simulate a larger number keep looping around the list. Will this cause a denial of service?

4. Take a wide sample, say 1000 pixels apart. If the pixels are the same value assume all pixels have the same value in between. If not, cut
the sample in half to find where the pixels would be the same.

5. Benchmark how running different numbers of goroutines would affect performance.
	Should the goroutine start after the URLs get extracted or part of the
	extraction process?
	a. Launch a goroutine for each URL
	b. Launch a goroutine for every 10 URLs
	c. Launch a goroutine for every 100 URLs
	d. Launch a goroutine for every 1000 URLs

6. Once program runs dockerize it.
*/

func csvSetup(filename string) string {

	filename = exporttocsv.CreateCSV(filename)
	headerRecord := []string{"URL", "top_color1", "top_color2", "top_color3"}
	exporttocsv.Export(filename, headerRecord)

	return "CSV setup complete."
}

func extractURLs(filename string) (map[string]map[colorCode]int, string) {
	// Think on this, should I do a batch extraction or have a go func deal
	// with each individual URL with a pointer/counter to reference the last
	// URL extracted?

	f, err := os.Open(filename)
	log.ErrorCheck("Fatal", "URL extraction failed during setup.", err)
	defer f.Close()

	// Continue to think on data structure
	// URL is the key
	// URL represents an image with RGB color codes
	// Color codes are a key
	// The number of times a color code appears is the value
	imgColorPrevalence := make(map[string]map[colorCode]int)

	scanner := bufio.NewScanner(f)

	// Default behavior is to scan line-by-line
	for scanner.Scan() {
		imgColorPrevalence[scanner.Text()] = make(map[colorCode]int)
	}

	return imgColorPrevalence, "URLs extracted."
}
