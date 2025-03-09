package fastfind

import (
	"fmt"
	"strings"

	unarr "github.com/gen2brain/go-unarr"
	log "github.com/sirupsen/logrus"
)

func readFileContentFromArchiveUnarr(archive *unarr.Archive, fname string) ([]byte, error) {

	err := archive.EntryFor(fname)
	if err != nil {
		return nil, err
	}

	data, err := archive.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ProcessFile - parse a FastFind Archive into arrays of results
func ProcessFileUnarr(fname string, matches []*FastFindMatch, computers []*FastFindComputer) ([]*FastFindMatch, []*FastFindComputer, error) {

	// early return to skip un-interesting files
	if !strings.HasSuffix(fname, ".7z") {
		log.Debug("Skipping " + fname)
		return matches, computers, fmt.Errorf("%s is not a FastFind archive", fname)
	}

	log.Debug("Processing " + fname)
	archive, err := unarr.NewArchive(fname)
	if err != nil {
		log.Errorf("Failed to open '%s' with error: %v", fname, err)
		return matches, computers, err
	}
	defer archive.Close()

	/////////// List all files inside archive
	var mainLogFile string
	files, err := archive.List()
	if err != nil {
		log.Errorf("Failed to list '%s' content with error : %v", fname, err)
		return matches, computers, err
	}
	for _, e := range files {
		log.Tracef("Archive content name: %s", e)
		if e == "FastFind.log" {
			mainLogFile = e
		}
	}
	if mainLogFile == "" {
		return matches, computers, fmt.Errorf("no FastFind.log in archive '%s'", fname)
	}

	////////// Process LogFile to determine ORC Version
	mainLogContent, err := readFileContentFromArchiveUnarr(archive, mainLogFile)
	if err != nil {
		log.Errorf("Failed to decompress log '%s' from archive '%s' with: %v", mainLogFile, fname, err)
		return matches, computers, err
	}
	orcVersion, isModern, err := scanORCVersion(string(mainLogContent))
	if err != nil {
		log.Errorf("Failed to scan ORC Version from archive '%s' with error : %v", fname, err)
		return matches, computers, err
	}
	log.Debugf("Orc version: %s (modern:%v)", orcVersion, isModern)

	/////////// Process FastFind_Result.xml

	// decompress to in-memory buffer for result file
	resultData, err := readFileContentFromArchiveUnarr(archive, resultsFname)
	if err != nil {
		log.Errorf("Failed to decompress results '%s' from archive '%s' with: %v", resultsFname, fname, err)
		return matches, computers, err
	}

	// parsing
	if isModern {
		return processFFResultsNG(fname, orcVersion, resultData, matches, computers)
	} else {
		return processFFResultsLegacy(fname, orcVersion, resultData, matches, computers)
	}
}
