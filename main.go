package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <path> <search string1> [<search string2> ...]")
		return
	}

	path := os.Args[1]
	searchStrings := os.Args[2:]

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
				yamlFile, err := ioutil.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}
				var yamlData interface{}
				err = yaml.Unmarshal(yamlFile, &yamlData)
				if err == nil {
					findStringInYAML(yamlData, path, searchStrings)
				}
			} else {
				file, err := os.Open(path)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				scanner.Buffer(nil, 1048576*10) // Set scanner buffer to 10MB

				for scanner.Scan() {
					line := scanner.Text()
					for _, searchString := range searchStrings {
						if strings.Contains(line, searchString) {
							re := regexp.MustCompile(searchString)
							fmt.Printf("%s:%s\n", path, re.ReplaceAllStringFunc(line, func(match string) string {
								return fmt.Sprintf("\033[1;31m%s\033[0m", match)
								// The above line adds ANSI escape codes to highlight the matched string in red color
							}))
						}
					}
				}
				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func findStringInYAML(yamlData interface{}, path string, searchStrings []string) {
	switch node := yamlData.(type) {
	case map[string]interface{}:
		for _, value := range node {
			findStringInYAML(value, path, searchStrings)
		}
	case []interface{}:
		for _, value := range node {
			findStringInYAML(value, path, searchStrings)
		}
	case string:
		for _, searchString := range searchStrings {
			if strings.Contains(node, searchString) {
				re := regexp.MustCompile(searchString)
				fmt.Printf("%s: %s\n", path, re.ReplaceAllStringFunc(node, func(match string) string {
					return fmt.Sprintf("\033[1;31m%s\033[0m", match)
					// The above line adds ANSI escape codes to highlight the matched string in red color
				}))
			}
		}
	}
}
