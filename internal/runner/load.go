package runner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

func LoadSpec(filepath string) (spec *Spec, err error) {
	switch {
	case strings.HasSuffix(filepath, ".json"):
		spec, err = loadJSONSpec(filepath)
	case strings.HasSuffix(filepath, ".json"):
		spec, err = loadJSONSpec(filepath)
	case strings.HasSuffix(filepath, ".yaml"), strings.HasSuffix(filepath, ".yml"):
		spec, err = loadYAMLSpec(filepath)
	default:
		err = fmt.Errorf("invalid file ending: requires JSON or YAML file")
		return
	}

	if err != nil {
		return
	}

	return
}

func loadJSONSpec(filepath string) (*Spec, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var spec Spec
	err = json.Unmarshal(bytes, &spec)
	return &spec, err
}

func loadYAMLSpec(filepath string) (*Spec, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var spec Spec
	err = yaml.Unmarshal(bytes, &spec)
	return &spec, err
}
