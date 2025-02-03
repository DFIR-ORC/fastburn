package fastfind

import (
	"fmt"
	"io"
	"strings"

	"github.com/kjk/lzmadec"
	log "github.com/sirupsen/logrus"
)

// readFileContentFromArchive - decompress to in-memory buffer for result file
func readFileContentFromArchiveLZMADec(archive *lzmadec.Archive, fname string) ([]byte, error) {

	r, err := archive.GetFileReader(fname)
	if err != nil {
		log.Error("readFileContentFromArchive('" + fname + "') failed with :" + err.Error())
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}

// ProcessFile - parse a FastFind Archive into arrays of results
func ProcessFileLZMADec(fname string, matches []*FastFindMatch, computers []*FastFindComputer) ([]*FastFindMatch, []*FastFindComputer, error) {

	// early return to skip un-interesting files
	if !strings.HasSuffix(fname, ".7z") {
		log.Debug("Skipping " + fname)
		return matches, computers, fmt.Errorf("%s is not a FastFind archive", fname)
	}

	log.Debug("Processing " + fname)
	var archive *lzmadec.Archive
	var err error
	archive, err = lzmadec.NewArchive(fname)
	if err != nil {
		log.Error("Call to 7z tool to process '" + fname + "') failed with :" + err.Error())
		return matches, computers, err
	}

	/////////// List all files inside archive
	var emocheckFile string
	var mainLogFile string
	for _, e := range archive.Entries {
		isEmocheck := IsEmocheckResult(e.Path)
		log.Trace(fmt.Sprintf(
			"Archive content name: %s, size: %d is_emocheck:%v",
			e.Path, e.Size, isEmocheck))
		if isEmocheck {
			emocheckFile = e.Path
		}
		if e.Path == "FastFind.log" {
			mainLogFile = e.Path
		}
	}

	/////////// Process emocheck result
	emotetInfected := false
	if emocheckFile != "" {
		log.Trace(fmt.Sprintf("Emocheck result found: %s", emocheckFile))
		emocheck_data, err := readFileContentFromArchiveLZMADec(archive, emocheckFile)
		if err != nil {
			log.Errorf(
				"Failed to decompress '%s' from archive '%s' with failed with :%v",
				emocheckFile, fname, err)
			return matches, computers, err
		}
		emotetInfected = processEmocheck(fname, emocheck_data)
	}

	////////// Process LogFile to determine ORC Version
	mainLogContent, err := readFileContentFromArchiveLZMADec(archive, mainLogFile)
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
	log.Info("Orc version: " + orcVersion)

	/////////// Process FastFind_Result.xml

	// decompress to in-memory buffer for result file
	resultData, err := readFileContentFromArchiveLZMADec(archive, resultsFname)
	if err != nil {
		log.Error(fmt.Sprintf(
			"Failed to decompress '%s' from archive '%s' with failed with :%s",
			emocheckFile, fname, err.Error()))
		return matches, computers, err
	}

	// parsing
	if isModern {
		return processFFResultsNG(fname, orcVersion, resultData, emotetInfected, matches, computers)
	} else {
		return processFFResultsLegacy(fname, orcVersion, resultData, emotetInfected, matches, computers)
	}
}
