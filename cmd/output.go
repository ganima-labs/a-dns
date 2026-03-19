package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func outputResult(data interface{}) error {
	switch outputFormat {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		defer encoder.Close()
		return encoder.Encode(data)
	case "table":
		return printTable(data)
	default:
		return fmt.Errorf("format non supporté: %s", outputFormat)
	}
}

func printTable(data interface{}) error {
	switch v := data.(type) {
	case []interface{}:
		fmt.Println("DATA:", v)
	default:
		fmt.Printf("%+v\n", v)
	}
	return nil
}
