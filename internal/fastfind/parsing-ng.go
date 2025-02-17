package fastfind

import (
	"encoding/xml"

	log "github.com/sirupsen/logrus"
)

func processFFResultsNG(fname string, orcVersion string, resultData []byte, matches []*FastFindMatch, computers []*FastFindComputer) ([]*FastFindMatch, []*FastFindComputer, error) {
	log.Tracef("Processing FF results modern format (v:%s)", orcVersion)

	// parsing a FastFind XML result
	var results FastFindResultNg
	xml.Unmarshal(resultData, &results)

	c := createFastFindComputer(fname, results.Computer, results.OS, results.Role, orcVersion)

	log.Trace("Processing filesystem matches")
	var nbm uint
	matches, nbm = recordFilesystemMatches(
		fname,
		results.Computer, results.OS, results.Role, orcVersion,
		results.Filesystem.FSMatches,
		matches)
	c.NbMatches += nbm
	log.Tracef("%d filesystem matches found", nbm)

	log.Trace("Processing registry matches")
	var nreg uint
	for _, hive := range results.Registry.Hive {
		log.Tracef("Hivepath: %s", hive.HivePath)

		for _, match := range hive.RegMatches {
			description := match.Description
			for _, k := range match.Values {
				log.Tracef("Description:'%s' Key [%s] (lastmodif:%v subkeycount:%v valcount:%v)",
					description,
					k.Key,
					k.LastmodifiedKey, k.SubkeysCount, k.ValuesCount)

				matches, nbm = recordRegMatchNG(
					fname,
					results.Computer, results.OS, results.Role, orcVersion,
					hive.HivePath, hive.VolumeID, hive.SnapshotID,
					description,
					k.Key, k.Value, k.Type, k.Size,
					k.LastmodifiedKey, k.SubkeysCount, k.ValuesCount,
					matches)

				c.NbMatches += nbm
				nreg += nbm
			}
		}
	}
	log.Tracef("%d registries matches found", nreg)

	computers = append(computers, &c)
	log.Debugf("Processing result: File '%s', Hostname %s matches: %v ",
		fname, c.Hostname, c.NbMatches)

	return matches, computers, nil
}
