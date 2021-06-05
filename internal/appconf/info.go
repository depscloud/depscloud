package appconf

import (
	"encoding/json"
)

// Current returns the running applications version.
func Current() *V {
	return &V{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}
}

// V wraps version metadata about the current application.
type V struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func (i V) String() string {
	data, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(data)
}
