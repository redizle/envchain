package envfile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type schemaFile struct {
	Rules []struct {
		Key      string `yaml:"key"`
		Required bool   `yaml:"required"`
		Pattern  string `yaml:"pattern"`
	} `yaml:"rules"`
}

// LoadSchema reads a YAML schema file and returns a Schema.
// The YAML format is:
//
//	rules:
//	  - key: PORT
//	    required: true
//	    pattern: '^\d+$'
func LoadSchema(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Schema{}, fmt.Errorf("read schema file: %w", err)
	}

	var sf schemaFile
	if err := yaml.Unmarshal(data, &sf); err != nil {
		return Schema{}, fmt.Errorf("parse schema file: %w", err)
	}

	schema := Schema{Rules: make([]SchemaRule, 0, len(sf.Rules))}
	for _, r := range sf.Rules {
		if r.Key == "" {
			return Schema{}, fmt.Errorf("schema rule is missing 'key' field")
		}
		schema.Rules = append(schema.Rules, SchemaRule{
			Key:      r.Key,
			Required: r.Required,
			Pattern:  r.Pattern,
		})
	}

	return schema, nil
}
