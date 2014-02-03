package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

var config_file string = "config.json"
var shots_dir string = "./shots/"

type Config struct {
	Shots_dir string
	Pages     map[string]string
}

func main() {
	fmt.Println("******************************")
	fmt.Println("* Site version image compare *")
	fmt.Println("******************************")

	config := readConfiguration()
	runApp(config)
}

func readConfiguration() *Config {
	file_content, err := os.Open(config_file)
	handleError(err, "Error occured during config file read")

	decoder := json.NewDecoder(file_content)
	configuration := &Config{}
	decoder.Decode(&configuration)

	shots_dir = configuration.Shots_dir

	return configuration
}

func runApp(configuration *Config) {
	for id, url := range configuration.Pages {
		fmt.Println(">> Process " + id + ": " + url)
		generateShotAndDiff(id, url)
	}
}

func generateShotAndDiff(id string, url string) {
	file_name_format := "shot_" + id + "_%d.png"
	old_id := lastGenerationID(id)
	new_id := old_id + 1
	screenshot_name := fmt.Sprintf(file_name_format, new_id)

	cmd_capture := exec.Command("phantomjs", "capture.js", url, screenshot_name)
	err := cmd_capture.Run()
	handleError(err, "Capture cannot run")

	// There is an old version.
	if old_id > 0 {
		file_name_old := fmt.Sprintf(shots_dir+file_name_format, old_id)
		file_name_new := fmt.Sprintf(shots_dir+file_name_format, new_id)
		file_name_diff := fmt.Sprintf(shots_dir+"diff_"+file_name_format, new_id)

		cmd_diff := exec.Command("compare", file_name_old, file_name_new, file_name_diff)
		err := cmd_diff.Run()
		handleError(err, "Cannot create diff")

		fmt.Println(">> Created new diff: " + file_name_diff)
	}
}

func lastGenerationID(id string) uint64 {
	file, err := os.Open(shots_dir)
	handleError(err, "Cannot open shots dir")

	fi, err := file.Readdir(0)
	handleError(err, "Cannot scan shots dir")

	var max_id uint64 = 0
	rx, _ := regexp.Compile("^shot_" + id + "_(\\d+)\\.png")

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

func handleError(err error, message string) {
	if err != nil {
		log.Fatalln(message, err)
	}
}
