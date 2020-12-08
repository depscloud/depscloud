package v

import (
	"encoding/json"
)

// Info wraps version metadata about the current application.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func (i Info) String() string {
	data, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(data)
}
