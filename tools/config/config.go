package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TemplateConfig holds metadata for a single mail template type.
type TemplateConfig struct {
	Name   string         `json:"name"`
	Desc   string         `json:"desc"`
	Path   []string       `json:"path"`
	Params map[string]any `json:"params"`
}

// TemplatesConfig is the root structure of templates_config.json.
type TemplatesConfig struct {
	Templates map[string]TemplateConfig `json:"templates"`
}

// Load reads and parses a templates_config.json file.
func Load(path string) (*TemplatesConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file %s: %w", path, err)
	}

	var cfg TemplatesConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config file %s: %w", path, err)
	}

	if len(cfg.Templates) == 0 {
		return nil, fmt.Errorf("config file %s contains no template definitions", path)
	}

	return &cfg, nil
}

// PathStr returns the file path of a template definition joined with OS separator.
func (t *TemplateConfig) PathStr() string {
	return filepath.Join(t.Path...)
}

// Registry builds a lookup map suitable for serializing into rendered.js
// as window.__REGISTRY__: maps template ID → {name, desc, path}.
func Registry(cfg *TemplatesConfig) map[string]map[string]string {
	reg := make(map[string]map[string]string, len(cfg.Templates))
	for id, t := range cfg.Templates {
		reg[id] = map[string]string{
			"name": t.Name,
			"desc": t.Desc,
			"path": strings.Join(t.Path, "/"),
		}
	}
	return reg
}

// FlattenParams converts nested template params into a flat dot-notation map
// suitable for the preview panel display in rendered.js.
func FlattenParams(params map[string]any) map[string]string {
	flat := make(map[string]string)
	flatten("", params, flat)
	return flat
}

func flatten(prefix string, val any, out map[string]string) {
	switch v := val.(type) {
	case map[string]any:
		for k, inner := range v {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}
			flatten(key, inner, out)
		}
	case []any:
		out[prefix] = fmt.Sprintf("[%d items]", len(v))
	default:
		out[prefix] = fmt.Sprintf("%v", v)
	}
}
