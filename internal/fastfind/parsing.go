package fastfind

/*


 **/

import (
	"fastburn/internal/utils"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func recordRegMatch(
	fname string, computer string, os string, role string, orcVersion string,
	hivepath string, volumeid string, snapshotid string,
	description string, value FastFind_RegValue,
	matches []*FastFindMatch) ([]*FastFindMatch, uint) {

	var nbmatches uint = 0

	m := FastFindMatch{
		Kind:             RegistryMatchType,
		Fullname:         hivepath,
		Reason:           description,
		Computer:         computer,
		ComputerOS:       os,
		ComputerRole:     role,
		ORCVersion:       orcVersion,
		VolumeID:         volumeid,
		SnapshotID:       snapshotid,
		LastModification: value.LastModifiedKey,
		Size:             value.DataSize,
		RegKey:           value.Key,
		RegType:          value.Type,
		RegValue:         value.Value,
		ArchiveName:      fname}

	if m.RegKey != "" || m.RegValue != "" {
		matches = append(matches, &m)
		nbmatches++
	}
	return matches, nbmatches
}
func recordFilesystemMatches(
	fname string, computer string, os string, role string, orcVersion string,
	fsmatches []FastFind_FileMatch, matches []*FastFindMatch) ([]*FastFindMatch, uint) {
	var nbmatches uint = 0
	for _, match := range fsmatches {
		m := FastFindMatch{
			Kind:             FilesystemMatchType,
			Reason:           match.Description,
			Computer:         computer,
			ComputerOS:       os,
			ComputerRole:     role,
			ORCVersion:       orcVersion,
			VolumeID:         match.Record.VolumeID,
			SnapshotID:       match.Record.SnapshotID,
			Creation:         match.Record.StandardInformation.Creation,
			LastModification: match.Record.StandardInformation.LastModification,
			LastEntryChange:  match.Record.StandardInformation.LastEntryChange,
			LastAccess:       match.Record.StandardInformation.LastAccess,
			Size:             match.Record.Data.FileSize,
			MD5:              match.Record.Data.MD5,
			SHA1:             match.Record.Data.SHA1,
			SHA256:           match.Record.Data.SHA256,
			ArchiveName:      fname,
			Ignore:           false,
		}
		/* Multiple $FILE_NAME
		- two: fullname + dos8.3 are very common
		- more (with high cardinality): for hardlink mostly
		*/
		nbFilenames := len(match.Record.Filenames)
		if nbFilenames > 0 {
			fn := match.Record.Filenames[0]
			m.Fullname = fn.Fullname
			m.FilenameCreation = fn.Creation
			m.FilenameLastModification = fn.LastModification
			m.FilenameLastEntryChange = fn.LastEntryChange
			m.FilenameLastAccess = fn.LastAccess
			m.ParentFrn = fn.Parentfrn
			if len(match.Record.Filenames) > 1 {
				fn = match.Record.Filenames[1]
				m.AltFilename = fn.Fullname
				m.AltFilenameCreation = fn.Creation
				m.AltFilenameLastModification = fn.LastModification
				m.AltFilenameLastEntryChange = fn.LastEntryChange
				m.AltFilenameLastAccess = fn.LastAccess
				m.AltParentFrn = fn.Parentfrn
			}
			if len(match.Record.Filenames) > 2 {
				log.Warning(
					fmt.Sprintf("Filesystem entry has %d $FILE_NAME entries", nbFilenames))
			}
		}
		/* Is it actually an I30 match ?
		   It is assumed we wont have mixed FILENAME and I30 matches
		*/
		nbI30 := len(match.Record.I30s)
		if nbI30 > 0 {
			if nbI30 > 1 {
				log.Warning(
					fmt.Sprintf("Filesystem entry has %d $I30 entries, only considering the first one", nbI30))
			}
			i30 := match.Record.I30s[0]
			log.Trace("I30 match for " + i30.Fullname)
			m.Kind = I30MatchType
			m.Fullname = i30.Fullname
			m.FilenameCreation = i30.Creation
			m.FilenameLastModification = i30.LastModification
			m.FilenameLastEntryChange = i30.LastEntryChange
			m.FilenameLastAccess = i30.LastAccess
			m.ParentFrn = i30.Parentfrn
		}
		//adding entry
		matches = append(matches, &m)
		nbmatches += 1
	}

	return matches, nbmatches
}

func createFastFindComputer(fname string, computer string, os string, role string, orcversion string) FastFindComputer {
	return FastFindComputer{
		Hostname:    computer,
		OS:          os,
		Role:        role,
		ORCVersion:  orcversion,
		ArchiveName: fname,
		NbMatches:   0}
}

// Read a FastFind logfile and try to extract a version string
// returns: the version string, true for modern vesrions, false for legacy versions <10.2.2, the error in case of error
func scanORCVersion(logContent string) (string, bool, error) {
	if logContent == "" {
		return "", false, nil
	}
	versionMatches, err := utils.ReScanStrings(logContent, `^FastFind v(\d+\.\d+(\.\d+)?)`)
	if err != nil {
		log.Errorf("failed to parse ORC log content with error :%v", err)
		return "", false, err
	}
	longVersion := versionMatches[0]
	//TODO extract short version
	if utils.VersionOrdinal(longVersion) < utils.VersionOrdinal("10.2.2") {
		return longVersion, false, nil
	} else {
		return longVersion, true, nil
	}

}

// Wrapper
func ProcessFile(fname string, matches []*FastFindMatch, computers []*FastFindComputer) ([]*FastFindMatch, []*FastFindComputer, error) {
	return ProcessFileUnarr(fname, matches, computers)
}
