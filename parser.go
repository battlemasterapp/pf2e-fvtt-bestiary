package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var localizationDict map[string]interface{}

func loadLocalization(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &localizationDict)
}

func getLocalizedText(id string) (string, error) {
	parts := strings.Split(id, ".")
	var current interface{} = localizationDict
	for _, part := range parts {
		if dict, ok := current.(map[string]interface{}); ok {
			current = dict[part]
		} else {
			return "", fmt.Errorf("key %s not found", id)
		}
	}
	// insert special characters
	localizedText := fmt.Sprintf("%v", current)
	localizedText = strings.ReplaceAll(localizedText, "\n", "\\n")
	localizedText = strings.ReplaceAll(localizedText, "\t", "\\t")
	return localizedText, nil
}

// Main parser
func parseText(input string) string {
	uuidRegex := regexp.MustCompile(`@UUID\[(.+?)\](?:{(.+?)})?`)
	localizeRegex := regexp.MustCompile(`@Localize\[(.+?)\](?:{(.+?)})?`)
	checkRegex := regexp.MustCompile(`@Check\[(.*?)\]`)
	templateRegex := regexp.MustCompile(`@Template\[(.*?)\](?:{(.+?)})?`)
	damageRegex := regexp.MustCompile(`@Damage\[(\d+d\d+(?:\[.*?\])?)\]`)

	// UUID parser
	result := uuidRegex.ReplaceAllStringFunc(input, func(match string) string {
		matches := uuidRegex.FindStringSubmatch(match)
		id, name := matches[1], matches[2]
		if name != "" {
			return name
		}
		// If there's no name, use the last id part
		parts := strings.Split(id, ".")
		return parts[len(parts)-1]
	})

	// Localize parser
	result = localizeRegex.ReplaceAllStringFunc(result, func(match string) string {
		matches := localizeRegex.FindStringSubmatch(match)
		id, name := matches[1], matches[2]
		localizedText, err := getLocalizedText(id)
		if err != nil {
			return err.Error()
		}
		if name != "" {
			return name
		}
		return localizedText
	})

	// Check parser
	result = checkRegex.ReplaceAllStringFunc(result, func(match string) string {
		matches := checkRegex.FindStringSubmatch(match)
		params := matches[1]
		options := strings.Split(params, "|")
		var dc, ability, traits, name, basic string
		for _, option := range options {
			if strings.HasPrefix(option, "dc:") {
				dc = strings.TrimPrefix(option, "dc:")
			} else if option == "basic" {
				basic = "Basic "
			} else if option == "name" {
				name = "Name "
			} else {
				ability = option
			}
		}
		result := fmt.Sprintf("%sDC %s %s", name, dc, basic+ability)
		if traits != "" {
			result += fmt.Sprintf(" (%s)", traits)
		}
		return result
	})

	// Template parser
	result = templateRegex.ReplaceAllStringFunc(result, func(match string) string {
		matches := templateRegex.FindStringSubmatch(match)
		params := matches[1]
		name := matches[2]
		var templateType string
		var distance string
		for _, param := range strings.Split(params, "|") {
			if strings.HasPrefix(param, "distance:") {
				distance = strings.TrimPrefix(param, "distance:")
			} else {
				templateType = param
			}
		}
		if name != "" {
			return name
		}
		return fmt.Sprintf("%s-foot %s", distance, templateType)
	})

	// Damage parser
	result = damageRegex.ReplaceAllStringFunc(result, func(match string) string {
		matches := damageRegex.FindStringSubmatch(match)
		ndmType := matches[1]
		parts := strings.Split(ndmType, "[")
		dice := parts[0]
		damageType := ""
		if len(parts) > 1 {
			damageType = strings.Trim(parts[1], "]")
			if damageType == "untyped" {
				return dice
			}
		}
		return fmt.Sprintf("%s %s", dice, damageType)
	})

	return result
}

func processFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	result := parseText(string(data))

	// Overwrite file with parsed text
	err = ioutil.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		return err
	}
	return nil
}

func processDirectory(dirPath string) error {
	return filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".json") {
			fmt.Println("Processing file:", path)
			err = processFile(path)
			if err != nil {
				fmt.Println("Error processing file:", err)
			}
		}
		return nil
	})
}

func main() {
	err := loadLocalization("./pf2e/static/lang/en.json")
	if err != nil {
		fmt.Println("Error loading file:", err)
		return
	}

	dirPath := "./bestiaries"
	err = processDirectory(dirPath)
	if err != nil {
		fmt.Println("Error processing directory:", err)
	}
}
