package fastfind

import (
	"encoding/xml"

	log "github.com/sirupsen/logrus"
)

func processFFResultsLegacy(fname string, orcVersion string, resultData []byte, emotetInfected bool, matches []*FastFindMatch, computers []*FastFindComputer) ([]*FastFindMatch, []*FastFindComputer, error) {
	// parsing a FastFind XML result
	var results FastFindResultLegacy
	xml.Unmarshal(resultData, &results)

	c := createFastFindComputer(fname, results.Computer, results.OS, results.Role, orcVersion, emotetInfected)

	var nbm uint
	matches, nbm = recordFilesystemMatches(
		fname,
		results.Computer, results.OS, results.Role, orcVersion,
		results.FSMatches,
		matches)
	c.NbMatches += nbm

	for _, match := range results.Registry {
		matches, nbm = recordRegMatch(
			fname,
			results.Computer, results.OS, results.Role, orcVersion,
			match.HivePath, match.VolumeID, match.SnapshotID,
			match.RegMatch.Description,
			match.RegMatch.Value, matches)
		c.NbMatches += nbm
	}

	computers = append(computers, &c)
	log.Debugf("Processing result: File '%s', Hostname %s Emotet: %v matches: %v ",
		fname, c.Hostname, c.EmotetInfected, c.NbMatches)

	return matches, computers, nil
}
