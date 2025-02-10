package filter

/*
 *
 * Loading of lists of f result filtering
 *
 */

import (
	"fastburn/internal/utils"
	"fmt"
	"os"
	"regexp"

	csv "github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

type criteriaEntry struct {
	SHA256      string `csv:"sha256"`
	SHA1        string `csv:"sha1"`
	MD5         string `csv:"md5"`
	FileRE      string `csv:"file_re"`
	Description string `csv:"description"`
	Regexp      *regexp.Regexp
}

type Filter struct {
	whitelistCriteria []criteriaEntry
	whitelistFilename string

	blacklistCriteria []criteriaEntry
	blacklistFilename string
}

func loadCSVCriterias(filename string, criterias *[]criteriaEntry) error {
	f, err := os.Open(filename)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to load filters from '%s': %s", filename, err.Error()))
		return err
	}
	defer f.Close()

	err = csv.UnmarshalFile(f, criterias)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to load filters from '%s': %s", filename, err.Error()))
		return err
	}
	log.Debug(fmt.Sprintf("%d values loaded from %s", len(*criterias), filename))

	for _, c := range *criterias {
		log.Tracef("compiling Regexp [%s]", c.FileRE)
		c.Regexp, err = regexp.Compile(c.FileRE)
		if err != nil {
			return fmt.Errorf("invalid regexp in flag list '%s': [%s]", filename, c.FileRE)
		}

	}

	return nil
}

func (f *Filter) LoadWhitelistCSV(filename string) error {
	log.Debug(fmt.Sprintf("Loading Postprocessing filter from '%s'", filename))
	f.whitelistFilename = filename

	return loadCSVCriterias(filename, &f.whitelistCriteria)
}

func (f *Filter) LoadBlacklistCSV(filename string) error {
	log.Debug(fmt.Sprintf("Loading Blacklist filter from '%s'", filename))
	f.blacklistFilename = filename

	return loadCSVCriterias(filename, &f.blacklistCriteria)
}

func (f *Filter) LoadLists(whiteFilename string, blackFilename string) error {

	if whiteFilename != "" {
		utils.PrintAndLog(log.InfoLevel, "Whitelist file: %s", whiteFilename)
		csv_err := f.LoadWhitelistCSV(whiteFilename)
		if csv_err != nil {
			log.Errorf("Failed to load post processing flags from '%s': %s", whiteFilename, csv_err.Error())
			return csv_err
		}
	}

	if blackFilename != "" {
		utils.PrintAndLog(log.InfoLevel, "Blacklist file: %s", whiteFilename)
		csv_err := f.LoadBlacklistCSV(blackFilename)
		if csv_err != nil {
			log.Errorf("Failed to load post processing blacklist from '%s': %s", blackFilename, csv_err.Error())
			return csv_err
		}
	}
	return nil
}
