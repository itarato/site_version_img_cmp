// Site image version compare.
// Main file.

// @todo Make width an array an iterate through.
// @todo Provide links in case it's a CI - could be handy.

package main

import (
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

// Path to the configuration file.
var config_file string = "config.json"

// Default directory for the screenshots.
var shots_dir string = "./shots/"

// File name string format of the screenshots.
// Args:
// 	1 - id of the screenhot, usually name of the page,
//  2 - width of the capruted frame,
//  3 - incremental number to separate versions.
var shot_name_format string = "shot_%s_%d_%d.png"

// Default width of the screenshot.
var shot_width int = 960

// Type for the configuration JSON file.
type Config struct {
	Width     int
	Shots_dir string
	Pages     map[string]string
}

// Entry point.
func main() {
	fmt.Println("******************************")
	fmt.Println("* Site version image compare *")
	fmt.Println("******************************")

	config := readConfiguration()
	runApp(config)
}

// Read configuration file and sets default variables.
func readConfiguration() *Config {
	file_content, err := os.Open(config_file)
	handleError(err, "Error occured during config file read")

	decoder := json.NewDecoder(file_content)
	configuration := &Config{}
	decoder.Decode(&configuration)

	shots_dir = configuration.Shots_dir
	shot_width = configuration.Width

	return configuration
}

// Execute actions.
func runApp(configuration *Config) {
	for id, url := range configuration.Pages {
		fmt.Println(">> Process " + id + ": " + url)
		generateShotAndDiff(id, url)
	}
}

// Handling one page scenario:
//  - creating a screenshot,
//  - correct image size issues,
//  - create diff.
func generateShotAndDiff(id string, url string) {
	old_id := lastGenerationID(id)
	new_id := old_id + 1
	screenshot_name := fmt.Sprintf(shot_name_format, id, shot_width, new_id)

	cmd_capture := exec.Command("phantomjs", "capture.js", url, screenshot_name)
	err := cmd_capture.Run()
	handleError(err, "Capture cannot run")

	// There is an old version.
	if old_id > 0 {
		generateDiff(id, old_id, new_id)
	} else {
		fmt.Println(">> No previous version of " + id)
	}
}

// Generate an image diff of two images.
func generateDiff(id string, old_id uint64, new_id uint64) {
	file_name_old := fmt.Sprintf(shots_dir+shot_name_format, id, shot_width, old_id)
	file_name_new := fmt.Sprintf(shots_dir+shot_name_format, id, shot_width, new_id)
	file_name_diff := fmt.Sprintf(shots_dir+"diff_"+shot_name_format, id, shot_width, new_id)

	var err_fix error
	file_name_old, file_name_new, err_fix = fixImageHight(file_name_old, file_name_new)
	handleError(err_fix, "Cannot resize")

	cmd_diff := exec.Command("compare", file_name_old, file_name_new, file_name_diff)
	err_run := cmd_diff.Run()
	handleError(err_run, "Cannot create diff")

	fmt.Println(">> Created new diff: " + file_name_diff)
}

// Check image sizes and synchronize them.
// Returns the new file names - as they might change during the resize.
func fixImageHight(file_a string, file_b string) (string, string, error) {
	height_old, err_a := getImageHeight(file_a)
	if err_a != nil {
		return "", "", err_a
	}

	height_new, err_b := getImageHeight(file_b)
	if err_b != nil {
		return "", "", err_b
	}

	file_a_new := file_a
	file_b_new := file_b
	if height_new > height_old {
		file_a_new = file_a + ".fixed.png"
		resizeImage(file_a, file_a_new, height_new)
	} else if height_new < height_old {
		file_b_new = file_b + ".fixed.png"
		resizeImage(file_b, file_b_new, height_old)
	}

	return file_a_new, file_b_new, nil
}

// Execute resize on an image.
func resizeImage(name string, output string, height int) {
	cmd := exec.Command("convert", name, "-extent", fmt.Sprintf("%dx%d", shot_width, height), output)
	err := cmd.Run()
	handleError(err, "Cannot resize image: "+name)
	fmt.Println(">> Corrected image size: " + output)
}

// Get the height of an image.
func getImageHeight(path string) (int, error) {
	reader, err_open := os.Open(path)
	if err_open != nil {
		return 0, err_open
	}

	image, err_decode := png.Decode(reader)
	if err_decode != nil {
		return 0, err_decode
	}

	return image.Bounds().Dy(), nil
}

// Get the last generated incremental id of the same type of screenshot.
// Returns 0 if it doesn't exist yet.
func lastGenerationID(id string) uint64 {
	file, err := os.Open(shots_dir)
	handleError(err, "Cannot open shots dir")

	fi, err := file.Readdir(0)
	handleError(err, "Cannot scan shots dir")

	var max_id uint64 = 0
	// @todo add the current size there, not just the pattern
	rx, _ := regexp.Compile("^shot_" + id + "_\\d+_(\\d+)\\.png$")

	for _, file_info := range fi {
		file_name := file_info.Name()
		if rx.MatchString(file_name) {
			id := rx.ReplaceAllString(file_name, "$1")
			current_id, err := strconv.ParseUint(id, 10, 32)
			handleError(err, "Cannot convert id to uint")
			if current_id > max_id {
				max_id = current_id
			}
		}
	}

	return max_id
}

// Simple error helper.
func handleError(err error, message string) {
	if err != nil {
		log.Fatalln(message, err)
	}
}
