package fastfind

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// ExportMatchesToCSV - Export a list of matches to a CSV file appending filter columens from the provided functions
func ExportMatchesToCSV(filename string, matches *FastFindMatchesList, isInfected FilterFunc) error {
	log.Debug("Exporting CSV matches to " + filename)
	fout, err := os.Create(filename)
	if err != nil {
		log.Trace("Export to " + filename + " :" + err.Error())
		return err
	}
	defer fout.Close()
	w := csv.NewWriter(fout)

	err = w.Write([]string{
		"Ignore", "Computer", "ComputerRole", "ComputerOS", "ORCVersion", "MatchType",
		"Software", "Infection", "Reason",
		"Filename", "AltName", "RegKey", "RegType", "RegValue",
		"FileSize", "MD5", "SHA1", "SHA256",
		"FileCreation", "FileLastModification", "FileLastEntryChange", "FileLastAccess",
		"FilenameCreation", "FilenameLastModification", "FilenameLastEntryChange", "FilenameLastAccess",
		"AltFilenameCreation", "AltFilenameLastModification", "AltFilenameLastEntryChange", "AltFilenameLastAccess",
		"VolumeID", "SnapshotID",
		"ArchiveName"})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to write to CSV file '%s': %v", filename, err))
	}

	// processing results
	for _, m := range *matches {
		log.Trace("Match " + m.Computer + " " + m.Fullname)
		isBackdoored, backDescr := isInfected(m)
		presMsg := ""
		backMsg := ""

		if isBackdoored {
			backMsg = backDescr
		}

		err = w.Write([]string{
			strconv.FormatBool(m.Ignore),
			m.Computer, m.ComputerRole, m.ComputerOS, m.ORCVersion, m.Kind.String(),
			presMsg, backMsg, m.Reason,
			m.Fullname, m.AltFilename, m.RegKey, m.RegType, m.RegValue,
			fmt.Sprintf("%d", m.Size), m.MD5, m.SHA1, m.SHA256,
			m.Creation, m.LastModification, m.LastEntryChange, m.LastAccess,
			m.FilenameCreation, m.FilenameLastModification, m.FilenameLastEntryChange, m.FilenameLastAccess,
			m.AltFilenameCreation, m.AltFilenameLastModification, m.AltFilenameLastEntryChange, m.AltFilenameLastAccess,
			m.VolumeID, m.SnapshotID,
			m.ArchiveName,
		})
		if err != nil {
			log.Error(fmt.Sprintf("Failed to match write to CSV file '%s': %v", filename, err))
		}
	}
	w.Flush()

	log.Trace("Matches CSV Export to " + filename + " done")
	return nil
}

// ExportMatchesToCSV - Export a list of matches to a CSV file appending filter columens from the provided functions
func ExportComputersToCSV(filename string, computers *FastFindComputersList) error {
	log.Debug("Exporting CSV computers to " + filename)
	fout, err := os.Create(filename)
	if err != nil {
		log.Trace("Export to " + filename + " :" + err.Error())
		return err
	}
	defer fout.Close()
	w := csv.NewWriter(fout)

	err = w.Write([]string{
		"Computer", "ComputerRole", "ComputerOS", "ORCVersion",
		"Infection", "NbMatches",
		"ArchiveName"})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to write to CSV file '%s': %v", filename, err))
	}

	// processing results
	for _, c := range *computers {
		log.Trace("Computer " + c.Hostname)
		back_msg := ""
		if c.EmotetInfected {
			back_msg = "Emotet detected"
		}
		err = w.Write([]string{
			c.Hostname, c.Role, c.OS, c.ORCVersion,
			back_msg, fmt.Sprintf("%v", c.NbMatches),
			c.ArchiveName})
		if err != nil {
			log.Error(fmt.Sprintf("Failed to computer write to CSV file '%s': %v", filename, err))
		}
	}
	w.Flush()

	log.Trace("Computers CSV Export to " + filename + " done")
	return nil
}

func ExportStatsToCSV(filename string, stats *FastFindMatchesStats) error {
	log.Debug("Exporting CSV stats to " + filename)

	fout, err := os.Create(filename)
	if err != nil {
		log.Trace("Export to " + filename + " :" + err.Error())
		return err
	}
	defer fout.Close()

	stats.ToCSV(fout)
	log.Trace("Stats CSV Export to " + filename + " done")

	return nil
}

func ExportTimelineToCSV(filename string, timeline *Timeline) error {
	log.Debug("Exporting timeline to " + filename)

	fout, err := os.Create(filename)
	if err != nil {
		log.Trace("Export to " + filename + " :" + err.Error())
		return err
	}
	defer fout.Close()

	timeline.ToCSV(fout)
	log.Trace("Timeline Export to " + filename + " done")

	return nil
}
