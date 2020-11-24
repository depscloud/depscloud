package v

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func (i Info) String() string {
	data, err := json.Marshal(i)
	if err != nil {
		logrus.Fatal(err)
	}
	return string(data)
}
