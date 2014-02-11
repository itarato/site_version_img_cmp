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
	"strings"
)

// Path to the configuration file.
var config_file string = "config.json"

// File name string format of the screenshots.
// Args:
// 	1 - id of the screenhot, usually name of the page,
//  2 - width of the capruted frame,
//  3 - incremental number to separate versions.
var shot_name_format string = "shot_%s_%d_%d.png"

// Type for the configuration JSON file.
type Config struct {
	Width                []int
	Shots_dir            string
	Shots_dir_public_url string
	Pages                map[string]string
}

// Global configuration.
var config *Config = &Config{}

// Entry point.
func main() {
	fmt.Println("******************************")
	fmt.Println("* Site version image compare *")
	fmt.Println("******************************")

	readConfiguration(config)
	runApp()
}

// Read configuration file and sets default variables.
func readConfiguration(configuration *Config) {
	file_content, err := os.Open(config_file)
	handleError(err, "Error occured during config file read")

	decoder := json.NewDecoder(file_content)
	decoder.Decode(&configuration)
}

// Execute actions.
func runApp() {
	task_count := len(config.Pages) * len(config.Width)
	c := make(chan bool, task_count)
	for id, url := range config.Pages {
		for _, width := range config.Width {
			go generateShotAndDiff(id, url, width, c)
		}
	}

	for ; task_count > 0; (func() { <-c; task_count-- })() {
	}
}

// Handling one page scenario:
//  - creating a screenshot,
//  - correct image size issues,
//  - create diff.
func generateShotAndDiff(id string, url string, width int, c chan bool) {
	fmt.Println(">> " + id + " | Process: " + url)

	old_id := lastGenerationID(id)
	new_id := old_id + 1
	screenshot_name := fmt.Sprintf(shot_name_format, id, width, new_id)

	cmd_capture := exec.Command("phantomjs", "capture.js", url, screenshot_name, strconv.Itoa(width))
	err := cmd_capture.Run()
	handleError(err, "Capture cannot run")

	// There is an old version.
	if old_id > 0 {
		generateDiff(id, old_id, new_id, width)
	} else {
		fmt.Println(">> " + id + " | No previous version")
	}

	c <- true
}

// Generate an image diff of two images.
func generateDiff(id string, old_num uint64, new_num uint64, width int) {
	file_name_old := fmt.Sprintf(config.Shots_dir+shot_name_format, id, width, old_num)
	file_name_new := fmt.Sprintf(config.Shots_dir+shot_name_format, id, width, new_num)
	file_name_diff := fmt.Sprintf(config.Shots_dir+"diff_"+shot_name_format, id, width, new_num)

	var err_fix error
	file_name_old, file_name_new, err_fix = fixImageHight(file_name_old, file_name_new, width)
	handleError(err_fix, "Cannot resize")

	cmd_diff := exec.Command("compare", "-metric", "PSNR", file_name_old, file_name_new, file_name_diff)
	// On some Unix based systems the exit status from compare is 1.
	// Even after -debug "All" -verbose the cause was unknown - and at the same time the diff was generated.
	// Avoiding error check until it's clear why is it happening.
	output, _ := cmd_diff.CombinedOutput()
	fmt.Println(">> " + id + " | Measured difference: " + strings.Trim(string(output), "\n\r\t "))
	fmt.Println(">> " + id + " | Created new diff: " + config.Shots_dir_public_url + fmt.Sprintf("diff_"+shot_name_format, id, width, new_num))
}

// Check image sizes and synchronize them.
// Returns the new file names - as they might change during the resize.
func fixImageHight(file_a string, file_b string, width int) (string, string, error) {
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
		resizeImage(file_a, file_a_new, height_new, width)
	} else if height_new < height_old {
		file_b_new = file_b + ".fixed.png"
		resizeImage(file_b, file_b_new, height_old, width)
	}

	return file_a_new, file_b_new, nil
}

// Execute resize on an image.
func resizeImage(name string, output string, height int, width int) {
	cmd := exec.Command("convert", name, "-extent", fmt.Sprintf("%dx%d", width, height), output)
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
	file, err := os.Open(config.Shots_dir)
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
