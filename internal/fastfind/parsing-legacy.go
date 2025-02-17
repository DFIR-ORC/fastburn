package fastfind

import (
	"encoding/xml"

	log "github.com/sirupsen/logrus"
)

func processFFResultsLegacy(fname string, orcVersion string, resultData []byte, matches []*FastFindMatch, computers []*FastFindComputer) ([]*FastFindMatch, []*FastFindComputer, error) {

	log.Tracef("Processing FF results legacy format (v:%s)", orcVersion)

	// parsing a FastFind XML result
	var results FastFindResultLegacy
	xml.Unmarshal(resultData, &results)

	c := createFastFindComputer(fname, results.Computer, results.OS, results.Role, orcVersion)

	log.Trace("Processing filesystem matches")
	var nbm uint
	matches, nbm = recordFilesystemMatches(
		fname,
		results.Computer, results.OS, results.Role, orcVersion,
		results.FSMatches,
		matches)
	c.NbMatches += nbm
	log.Tracef("%d filesystem matches found", nbm)

	log.Trace("Processing registry matches")
	var nreg uint
	for _, match := range results.Registry {
		matches, nbm = recordRegMatchLegacy(
			fname,
			results.Computer, results.OS, results.Role, orcVersion,
			match.HivePath, match.VolumeID, match.SnapshotID,
			match.RegMatch.Description,
			match.RegMatch.Value, matches)
		c.NbMatches += nbm
		nreg += nbm
	}
	log.Tracef("%d registries matches found", nreg)

	computers = append(computers, &c)
	log.Debugf("Processing result: File '%s', Hostname %s matches: %v ",
		fname, c.Hostname, c.NbMatches)

	return matches, computers, nil
}
