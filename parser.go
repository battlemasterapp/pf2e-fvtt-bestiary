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

func uuidParser(matches []string) string {
	_, name := matches[1], matches[2]
	if name != "" {
		return name
	}
	return ""
	// If there's no name, use the last id part
	// parts := strings.Split(id, ".")
	// return parts[len(parts)-1]
}

func localizeParser(matches []string) string {
	id, name := matches[1], matches[2]
	localizedText, err := getLocalizedText(id)
	if err != nil {
		return err.Error()
	}
	if name != "" {
		return name
	}
	return localizedText
}

func checkParser(matches []string) string {
	params := matches[1]
	options := strings.Split(params, "|")
	var dc, traits, name, basic string
	name = strings.Title(options[0])
	for _, option := range options[1:] {
		if strings.HasPrefix(option, "dc:") {
			dc = strings.TrimPrefix(option, "dc:")
		} else if strings.HasPrefix(option, "traits:") {
			traits = strings.TrimPrefix(option, "traits:")
		} else if option == "basic" {
			basic = "Basic "
		}
	}

	result := ""

	if dc != "" {
		result += fmt.Sprintf("DC %s ", dc)
	}
	if basic != "" {
		result += " Basic "
	}
	result += fmt.Sprintf(" %s", name)

	if traits != "" {
		result += fmt.Sprintf(" (%s)", traits)
	}

	// trim the space at the beginning
	return strings.Trim(result, " ")
}

func templateParser(matches []string) string {
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
}

func damageParser(matches []string) string {
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
		return uuidParser(matches)
	})

	// Localize parser
	result = localizeRegex.ReplaceAllStringFunc(result, func(match string) string {
		matches := localizeRegex.FindStringSubmatch(match)
		return localizeParser(matches)
	})

	// Check parser
	result = checkRegex.ReplaceAllStringFunc(result, func(match string) string {
		matches := checkRegex.FindStringSubmatch(match)
		return checkParser(matches)
	})

	// Template parser
	result = templateRegex.ReplaceAllStringFunc(result, func(match string) string {
		matches := templateRegex.FindStringSubmatch(match)
		return templateParser(matches)
	})

	// Damage parser
	result = damageRegex.ReplaceAllStringFunc(result, func(match string) string {
		matches := damageRegex.FindStringSubmatch(match)
		return damageParser(matches)
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
