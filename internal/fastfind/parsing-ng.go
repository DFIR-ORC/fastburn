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
	for _, hive := range results.Registry.Hives {
		log.Tracef("Hivepath: %s", hive.HivePath)

		for _, match := range hive.RegfindMatches {
			description := match.Description
			for _, k := range match.Values {
				log.Tracef("Description:'%s' Key [%s] Value[%s] (lastmodif:%v type:%v size:%d)",
					description,
					k.Key, k.Value,
					k.LastmodifiedKey, k.Type, k.DataSize,
				)

				matches, nbm = recordRegMatchNG(
					fname,
					results.Computer, results.OS, results.Role, orcVersion,
					hive.HivePath, hive.VolumeID, hive.SnapshotID,
					description,
					k.Key, k.Value, k.Type, k.DataSize,
					0 /*subkeys count*/, 0, /*values count*/
					k.LastmodifiedKey,
					matches)

				c.NbMatches += nbm
				nreg += nbm
			}
			for _, k := range match.Keys {
				log.Tracef("Description:'%s' Key [%s] (lastmodif:%s subkeycount:%d valcount:%d)",
					description,
					k.Key,
					k.LastmodifiedKey, k.SubkeysCount, k.ValuesCount)

				matches, nbm = recordRegMatchNG(
					fname,
					results.Computer, results.OS, results.Role, orcVersion,
					hive.HivePath, hive.VolumeID, hive.SnapshotID,
					description,
					k.Key, "" /*Value*/, "" /*Type*/, 0, /*Size*/
					k.SubkeysCount, k.ValuesCount,
					k.LastmodifiedKey,
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
