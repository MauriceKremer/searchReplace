package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// Config represents the JSON configuration
type Config struct {
	Folder       string            `json:"folder"`
	FileFilter   string            `json:"fileFilter"`
	Search       string            `json:"search"`
	Replacements map[string]string `json:"replacements"`
}

func main() {
	// Define and parse the command-line flag for selecting the environment
	env := flag.String("environment", "dev", "Environment to select the replacement string (e.g., dev, staging, production)")
	flag.Parse()

	// Read and parse the JSON configuration file
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Validate the environment and select the replacement pattern
	replacePattern, exists := config.Replacements[*env]
	if !exists {
		log.Fatalf("Invalid environment: %s", *env)
	}

	// Find files based on the folder and fileFilter from the config
	files, err := filepath.Glob(filepath.Join(config.Folder, config.FileFilter))
	if err != nil {
		log.Fatalf("Error finding files: %v", err)
	}

	// Compile the regex search pattern
	re, err := regexp.Compile(config.Search)
	if err != nil {
		log.Fatalf("Error compiling regex: %v", err)
	}

	// Iterate over each file and perform the search and replace
	for _, file := range files {
		fmt.Println("==============================================================================")
		fmt.Printf("Processing file: %s\n", file)

		// Read file content
		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf(" - error reading file %s: %v", file, err)
			continue
		}

		// Perform the search and replace
		updatedContent := re.ReplaceAllStringFunc(string(content), func(match string) string {
			replaced := re.ReplaceAllString(match, replacePattern)
			fmt.Printf(" - replaced '%s' with '%s' \n", match, replaced)
			return replaced
		})

		// Write the updated content back to the file
		err = os.WriteFile(file, []byte(updatedContent), os.ModePerm)
		if err != nil {
			log.Printf(" - error writing file %s: %v", file, err)
			continue
		}

		fmt.Println(" - successfully processed.")
	}
	fmt.Println("==============================================================================")

	fmt.Println("Search and replace completed.")
}
