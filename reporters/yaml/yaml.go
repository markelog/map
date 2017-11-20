package yaml

import (
	"github.com/markelog/map/spider"

	"github.com/ghodss/yaml"
)

// Execute yaml reporter
func Execute(data *spider.Result) (output string, err error) {
	result, err := yaml.Marshal(data)

	return string(result), err
}
