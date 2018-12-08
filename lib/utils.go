package lib

import (
	"encoding/json"
	native "log"
	"os"
	"time"
)

// Global loggers for this application
var (
	Logger = native.New(os.Stdout, "", 0)
	Errors = native.New(os.Stderr, "", 0)
)

// PrintJSON print JSON marshaled value
func PrintJSON(records interface{}) {
	marshaled, _ := json.MarshalIndent(records, "", "  ") // nolint
	Logger.Println(string(marshaled))
}

// TimeFormat format a time pointer
func TimeFormat(t *time.Time) string {
	if t == nil {
		return ""
	}
	casted := *t
	return casted.Format(time.RFC3339)
}
