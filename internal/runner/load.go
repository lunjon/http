package runner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

func Load(filepath string) (runner *Runner, err error) {
	var spec *Spec
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

	err = spec.Validate()
	if err != nil {
		return
	}

	spec.SetHeaders()

	runner = NewRunner(spec)
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
