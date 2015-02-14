// Site image version compare.
// Main file.

package main

import (
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// Path to the configuration file.
const config_file = "config.json"

// File name string format of the screenshots.
// Args:
// 	1 - id of the screenhot, usually name of the page,
//  2 - width of the capruted frame,
//  3 - incremental number to separate versions.
var shot_name_format string = "shot_%s_%d_%d.png"

// Root directory of the script.
var root_dir string

// Type of the page definition in the JSON config obejct.
type PageDef struct {
	Url           string `json:"url"`
	PhantomURLArg string `json:"phantom_url_arg"`
}

// Type for the configuration JSON file.
type Config struct {
	Width             []int              `json:"width"`
	ShotsDir          string             `json:"shots_dir"`
	ShotsDirPublicURL string             `json:"shots_dir_public_url"`
	Pages             map[string]PageDef `json:"pages"`
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
	_, filename, _, _ := runtime.Caller(0)
	root_dir = path.Dir(filename)

	file_content, err := os.Open(getPath(config_file))
	handleError(err, "Error occured during config file read")
	defer file_content.Close()

	decoder := json.NewDecoder(file_content)
	decoder.Decode(&configuration)
}

// Execute actions.
func runApp() {
	var wg sync.WaitGroup

	for id, page_def := range config.Pages {
		for _, width := range config.Width {
			wg.Add(1)
			go func(id string, page_def PageDef, width int) {
				generateShotAndDiff(id, page_def, width)
				defer wg.Done()
			}(id, page_def, width)
		}
	}

	wg.Wait()
}

// Handling one page scenario:
//  - creating a screenshot,
//  - correct image size issues,
//  - create diff.
func generateShotAndDiff(id string, page_def PageDef, width int) {
	fmt.Println(">> " + id + " | Process: " + page_def.Url)

	old_id := lastGenerationID(id)
	new_id := old_id + 1
	screenshot_name := getPath(fmt.Sprintf(config.ShotsDir+shot_name_format, id, width, new_id))

	cmd_capture := exec.Command("phantomjs", getPath("capture.js"), page_def.Url, screenshot_name, strconv.Itoa(width), page_def.PhantomURLArg)
	err := cmd_capture.Run()
	handleError(err, "Capture cannot run")

	// There is an old version.
	if old_id > 0 {
		generateDiff(id, old_id, new_id, width)
	} else {
		fmt.Println(">> " + id + " | No previous version")
	}
}

// Generate an image diff of two images.
func generateDiff(id string, old_num uint64, new_num uint64, width int) {
	file_name_old := getPath(fmt.Sprintf(config.ShotsDir+shot_name_format, id, width, old_num))
	file_name_new := getPath(fmt.Sprintf(config.ShotsDir+shot_name_format, id, width, new_num))
	file_name_diff := getPath(fmt.Sprintf(config.ShotsDir+"diff_"+shot_name_format, id, width, new_num))

	var err_fix error
	file_name_old, file_name_new, err_fix = fixImageHight(file_name_old, file_name_new, width)
	handleError(err_fix, "Cannot resize")

	cmd_diff := exec.Command("compare", "-metric", "PSNR", file_name_old, file_name_new, file_name_diff)
	// On some Unix based systems the exit status from compare is 1.
	// Even after -debug "All" -verbose the cause was unknown - and at the same time the diff was generated.
	// Avoiding error check until it's clear why is it happening.
	output, _ := cmd_diff.CombinedOutput()
	fmt.Println(">> " + id + " | Measured difference: " + strings.Trim(string(output), "\n\r\t "))
	fmt.Println(">> " + id + " | Created new diff: " + config.ShotsDirPublicURL + fmt.Sprintf("diff_"+shot_name_format, id, width, new_num))
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
	defer reader.Close()

	image, err_decode := png.Decode(reader)
	if err_decode != nil {
		return 0, err_decode
	}

	return image.Bounds().Dy(), nil
}

// Get the last generated incremental id of the same type of screenshot.
// Returns 0 if it doesn't exist yet.
func lastGenerationID(id string) uint64 {
	file, err := os.Open(getPath(config.ShotsDir))
	handleError(err, "Cannot open shots dir")
	defer file.Close()

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

// Extend path to absolute.
func getPath(file_path string) string {
	return path.Join(root_dir, file_path)
}

// Simple error helper.
func handleError(err error, message string) {
	if err != nil {
		log.Fatalln(message, err)
	}
}
