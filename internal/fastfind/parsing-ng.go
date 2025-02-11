package fastfind

import (
	"encoding/xml"

	log "github.com/sirupsen/logrus"
)

func processFFResultsNG(fname string, orcVersion string, resultData []byte, matches []*FastFindMatch, computers []*FastFindComputer) ([]*FastFindMatch, []*FastFindComputer, error) {
	// parsing a FastFind XML result
	var results FastFindResultNg
	xml.Unmarshal(resultData, &results)

	c := createFastFindComputer(fname, results.Computer, results.OS, results.Role, orcVersion)

	var nbm uint
	matches, nbm = recordFilesystemMatches(
		fname,
		results.Computer, results.OS, results.Role, orcVersion,
		results.Filesystem.FSMatches,
		matches)
	c.NbMatches += nbm

	for _, hive := range results.Registry.Hive {
		for _, match := range hive.RegMatches {
			for _, value := range match.Values {
				matches, nbm = recordRegMatch(
					fname,
					results.Computer, results.OS, results.Role, orcVersion,
					hive.HivePath, hive.VolumeID, hive.SnapshotID,
					match.Description,
					value, matches)
				c.NbMatches += nbm
			}
		}
	}

	computers = append(computers, &c)
	log.Debugf("Processing result: File '%s', Hostname %s matches: %v ",
		fname, c.Hostname, c.NbMatches)

	return matches, computers, nil
}
