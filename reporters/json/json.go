package json

import (
	"encoding/json"

	"github.com/markelog/map/spider"
)

// Execute json reporter
func Execute(data *spider.Result) (string, error) {
	result, err := json.Marshal(data)

	return string(result), err
}
