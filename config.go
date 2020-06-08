package main

import (
	"fmt"
	"io/ioutil"

	"github.com/shimmerglass/bar3x/ui"
	"gopkg.in/yaml.v2"
)

func getConfig(path string) (ui.Context, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	ctx := ui.Context{}
	err = yaml.Unmarshal(contents, &ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}

	return ctx, nil
}
