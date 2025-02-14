package fastfind

import (
	"fmt"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var emocheckFilePattern = `^emocheck_results\/.*_emocheck\.json$`

// EmocheckBoolean - just a wrapper type for a boolean serialized as yes/no in JSON
type EmocheckBoolean bool

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func (bit EmocheckBoolean) String() string {
	if bit {
		return "true"
	}
	return "false"
}

func (bit *EmocheckBoolean) UnmarshalJSON(data []byte) error {
	asString := string(data)
	var err error

	log.Tracef("Converting to bool: string %s", asString)
	asString = trimQuotes(asString)
	if asString == "yes" {
		*bit = true
	} else if asString == "no" {
		*bit = false
	} else {
		res, err := strconv.ParseBool(asString)
		if err != nil {
			return fmt.Errorf("boolean unmarshal error: invalid input %s", asString)
		}
		if res {
			*bit = true
		}
	}
	log.Tracef("Converted bool: string '%s', bool %v", asString, *bit)
	return err
}

// EmocheckMatch - Structure to load the JSON result of an Emocheck execution
type EmocheckMatch struct {
	ScanTime        string          `json:"scan_time"`
	Hostname        string          `json:"hostname"`
	EmocheckVersion string          `json:"emocheck_version"`
	IsInfected      EmocheckBoolean `json:"is_infected"`
}

// IsEmocheckResult - returns whether a filename looks like a legit result from Emocheck
func IsEmocheckResult(fname string) bool {
	matched, err := regexp.MatchString(emocheckFilePattern, fname)
	if err != nil {
		log.Fatal("Invalid FileRE in match criteria: " + emocheckFilePattern)
		return false
	}
	if matched {
		return true
	}
	return false
}
