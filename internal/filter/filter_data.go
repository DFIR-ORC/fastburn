package filter

/*
 *
 * Loading of lists of f result filtering
 *
 */

import (
	"fastburn/internal/utils"
	"fmt"
	"io"
	"os"
	"regexp"

	"encoding/csv"

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

var rxNotHex = regexp.MustCompile("[^0-9A-Fa-f]")

func isHexLen(val string, nbbits int) bool {
	nbHexCars := (nbbits / 8) * 2
	if len(val) != nbHexCars {
		return false
	}
	if rxNotHex.MatchString(val) {
		return false
	}
	return true
}

func checkHashOrEmpty(val string, nbBits int, name string, rowValid *bool, strErr *string) {
	if val == "" {
		return
	}
	if isHexLen(val, nbBits) {
		return
	}
	*strErr += fmt.Sprintf("'%s'is not a valid %s,", val, name)
	*rowValid = false
}

func loadCSVCriterias(filename string, criterias *[]criteriaEntry) error {

	log.Debugf("Loading filters from %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Debugf("Failed to load filters from '%s': %s", filename, err.Error())
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rowNum := 1     // row number
	filtersNum := 0 // filter number
	nbRegexp := 0   // regexp number

	// skipping headers
	if _, err := reader.Read(); err != nil {
		log.Debugf("Failed while loading filters from '%s': failed to skip headers - %s", filename, err.Error())
		return err
	}

	for {
		// reading and splitting a row
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Debugf("Failed while loading filters from '%s': %s", filename, err.Error())
				return err
			}
		}

		// fields validation
		if row[0] == "" && row[1] == "" && row[2] == "" && row[3] == "" {
			log.Infof("Empty filter on line %d of file %s", rowNum, filename)
		} else {
			var rowValid bool
			crit := criteriaEntry{SHA256: row[0], SHA1: row[1], MD5: row[2], FileRE: row[3], Description: row[4]}
			rowValid = true
			strErr := ""
			checkHashOrEmpty(crit.SHA256, 256, "SHA256", &rowValid, &strErr)
			checkHashOrEmpty(crit.SHA1, 160, "SHA1", &rowValid, &strErr)
			checkHashOrEmpty(crit.MD5, 128, "MD5", &rowValid, &strErr)
			if crit.FileRE != "" {
				log.Tracef("compiling regexp [%s]", crit.FileRE)
				crit.Regexp, err = regexp.Compile(crit.FileRE)
				if err != nil {
					strErr += fmt.Sprintf("[%s] is not a valid regexp", crit.FileRE)
					rowValid = false
				} else {
					nbRegexp += 1
					log.Tracef("compiled regexp [%v]", crit.Regexp)
				}
			}
			if rowValid {
				filtersNum += 1
				*criterias = append(*criterias, crit)
			} else {
				err := fmt.Errorf("invalid filter on row %d: %s", rowNum, strErr)
				log.Debugf("Failed while loading filters from '%s' : %s", filename, err.Error())
				return err
			}
		}
		rowNum += 1
	}
	log.Debugf("%d filters loaded from '%s'", len(*criterias), filename)
	log.Tracef("%d filters regexps compiled from '%s'", nbRegexp, filename)

	return nil
}

func (f *Filter) LoadWhitelistCSV(filename string) error {
	log.Debugf("Loading Postprocessing filter from '%s'", filename)
	f.whitelistFilename = filename

	return loadCSVCriterias(filename, &f.whitelistCriteria)
}

func (f *Filter) LoadBlacklistCSV(filename string) error {
	log.Debugf("Loading Blacklist filter from '%s'", filename)
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
		utils.PrintAndLog(log.InfoLevel, "Blacklist file: %s", blackFilename)
		csv_err := f.LoadBlacklistCSV(blackFilename)
		if csv_err != nil {
			log.Errorf("Failed to load post processing blacklist from '%s': %s", blackFilename, csv_err.Error())
			return csv_err
		}
	}

	log.Tracef("Filters loading %d blacklist criteria, %d whitelist criteria", len(f.blacklistCriteria), len(f.whitelistCriteria))
	return nil
}
