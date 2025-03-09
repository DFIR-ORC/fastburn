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
		log.Errorf("Call to 7z tool to process '%s') failed with error: %v", fname, err)
		return matches, computers, err
	}

	/////////// List all files inside archive
	var mainLogFile string
	for _, e := range archive.Entries {
		log.Tracef("Archive content name: %s, size: %d ", e.Path, e.Size)
		if e.Path == "FastFind.log" {
			mainLogFile = e.Path
		}
	}
	if mainLogFile == "" {
		return matches, computers, fmt.Errorf("no FastFind.log in archive '%s'", fname)
	}

	////////// Process LogFile to determine ORC Version
	mainLogContent, err := readFileContentFromArchiveLZMADec(archive, mainLogFile)
	if err != nil {
		log.Errorf("Failed to decompress log '%s' from archive '%s' with error: %v", mainLogFile, fname, err)
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
	resultData, err := readFileContentFromArchiveLZMADec(archive, resultsFname)
	if err != nil {
		log.Errorf("Failed to decompress results '%s' from archive '%s' with error: %v", resultsFname, fname, err)
		return matches, computers, err
	}

	// parsing
	if isModern {
		return processFFResultsNG(fname, orcVersion, resultData, matches, computers)
	} else {
		return processFFResultsLegacy(fname, orcVersion, resultData, matches, computers)
	}
}
