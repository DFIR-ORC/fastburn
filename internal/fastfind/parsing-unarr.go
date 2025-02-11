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
		log.Error("failed to open '" + fname + "' with :" + err.Error())
		return matches, computers, err
	}
	defer archive.Close()

	/////////// List all files inside archive
	var emocheckFile string
	var mainLogFile string
	files, err := archive.List()
	if err != nil {
		log.Error("failed to list '" + fname + "' content with :" + err.Error())
		return matches, computers, err
	}
	for _, e := range files {
		isemocheck := IsEmocheckResult(e)
		log.Tracef("Archive content name: %s, is_emocheck:%v", e, isemocheck)
		if isemocheck {
			emocheckFile = e
		}
		if e == "FastFind.log" {
			mainLogFile = e
		}
	}

	////////// Process LogFile to determine ORC Version
	mainLogContent, err := readFileContentFromArchiveUnarr(archive, mainLogFile)
	if err != nil {
		log.Errorf(
			"Failed to decompress '%s' from archive '%s' with failed with :%v",
			mainLogFile, fname, err)
		return matches, computers, err
	}
	orcVersion, isModern, err := scanORCVersion(string(mainLogContent))
	if err != nil {
		log.Errorf(
			"Failed to decompress '%s' from archive '%s' with failed with :%v",
			emocheckFile, fname, err)
		return matches, computers, err
	}
	log.Debug("Orc version: " + orcVersion)

	/////////// Process FastFind_Result.xml

	// decompress to in-memory buffer for result file
	resultData, err := readFileContentFromArchiveUnarr(archive, resultsFname)
	if err != nil {
		log.Error(fmt.Sprintf(
			"Failed to decompress '%s' from archive '%s' with failed with :%s",
			emocheckFile, fname, err.Error()))
		return matches, computers, err
	}

	// parsing
	if isModern {
		return processFFResultsNG(fname, orcVersion, resultData, matches, computers)
	} else {
		return processFFResultsLegacy(fname, orcVersion, resultData, matches, computers)
	}
}
