package main

/*
Supporting functions for fastburn command line tool

**/

import (
	"fmt"

	fbn "fastburn/internal/fastfind"
	"fastburn/internal/filter"
	"fastburn/internal/utils"

	log "github.com/sirupsen/logrus"
)

// parseFiles - process the command line to list files and parse them to the in-memory data structures
func parseFiles(args []string) ([]string, *fbn.FastFindMatchesList, *fbn.FastFindComputersList, *fbn.FastFindMatchesStats, error) {

	matches := make(fbn.FastFindMatchesList, 0)
	computers := make(fbn.FastFindComputersList, 0)

	stats := fbn.CreateStats()
	log.Debugf("Processing file list: %v", args)
	files, err := utils.ExpandArchiveFilePaths(utils.Uniq(args))
	if err != nil {
		log.Errorf("Failed to expand paths: %v", err)
		return nil, nil, nil, nil, err
	}
	log.Info(fmt.Sprintf("%d files to process", len(files)))
	log.Debug(fmt.Sprintf("Processing %d files", len(files)))

	for _, fname := range files {
		matches, computers, err = fbn.ProcessFile(fname, matches, computers)
		if err != nil {
			log.Warning("Failed to process '" + fname + "': " + err.Error())
		} else {
			c := computers[len(computers)-1]
			stats.UpdateComputers(c)
			if c.EmotetInfected {
				log.Warning("File '" + fname + "', Hostname " + c.Hostname + " matches: " + fmt.Sprintf("%v", c.NbMatches) + ": Emotet infected")
			} else {
				log.Info("File '" + fname + "', Hostname " + c.Hostname + " matches: " + fmt.Sprintf("%v", c.NbMatches))
			}
		}

	} //eo foreach filename
	return files, &matches, &computers, stats, nil
}

// analyseData - process the collected data in memory
func analyseData(matches *fbn.FastFindMatchesList, stats *fbn.FastFindMatchesStats, postfilter *filter.Filter) (*fbn.Timeline, error) {
	log.Debug(fmt.Sprintf("Post-processing %v results", len(*matches)))
	var blacklistCount uint64 = 0
	var whitelistCount uint64 = 0
	timeline := fbn.InitTL()

	// processing results: looking for sure matches
	for _, m := range *matches {

		// updating registers
		stats.UpdateMatches(m)
		timeline.Register(m)

		// processing remarquable IOCs
		isWhilelisted, back_descr := postfilter.IsWhitelisted(m)
		pres_msg := ""
		back_msg := ""

		if isWhilelisted {
			back_msg = back_descr
			if pres_msg != "" {
				pres_msg = ", " + pres_msg
			}
		}
		if isWhilelisted {
			log.Info("Detection on " + m.Computer + " [" + m.Kind.String() + "] : " + m.URI() + " -" + pres_msg + back_msg + " - Reason:<" + m.Reason + "> - Archive:" + m.ArchiveName)
			whitelistCount = whitelistCount + 1
		}

		// marking blackisted entries
		blacklisted, _ := postfilter.IsBlacklisted(m)
		if blacklisted {
			m.Ignore = true
			blacklistCount = blacklistCount + 1
		}
	}
	totalCount := len(*matches)
	log.Debugf("%d/%d whitelisted entries %d/%d blacklisted entries",
		whitelistCount, totalCount, blacklistCount, totalCount)

	return timeline, nil
}

// saveResults - export results to CSV files
func saveResults(csv_matches_fname string, csv_computers_fname string, csv_stats_fname string, timeline_fname string,
	files []string,
	matches *fbn.FastFindMatchesList, computers *fbn.FastFindComputersList, stats *fbn.FastFindMatchesStats, timeline *fbn.Timeline,
	postfilter *filter.Filter) error {

	// CSV export
	log.Debug("Exporting results")

	err := fbn.ExportComputersToCSV(csv_computers_fname, computers)
	if err != nil {
		log.Warning("Export to '" + csv_computers_fname + "' failed: " + err.Error())
	} else {
		log.Info("Computers exported to '" + csv_computers_fname + "'")
	}

	if len(*matches) > 0 {

		// is infected ?
		f := func(m *fbn.FastFindMatch) (bool, string) {
			return postfilter.IsWhitelisted(m)
		}
		// exporting
		err = fbn.ExportMatchesToCSV(csv_matches_fname, matches, f)
		if err != nil {
			log.Warning("Export to '" + csv_matches_fname + "' failed: " + err.Error())
		} else {
			log.Info("Matches exported to '" + csv_matches_fname + "'")
		}
	} else {
		log.Infof("No match found in %d computers", len(*computers))
	}
	log.Infof("%d archives processed for %d computers, %d matches found", len(files), len(*computers), len(*matches))

	// exporting stats
	err = fbn.ExportStatsToCSV(csv_stats_fname, stats)
	if err != nil {
		log.Warning("Stat report export to '" + csv_stats_fname + "' failed: " + err.Error())
	} else {
		log.Info("Stat report exported to '" + csv_stats_fname + "'")
	}

	// exporting timeline
	if timeline_fname != "" {
		fbn.ExportTimelineToCSV(timeline_fname, timeline)
		if err != nil {
			log.Warning("Timeline export to '" + timeline_fname + "' failed: " + err.Error())
		} else {
			log.Info("Timeline exported to '" + timeline_fname + "'")
		}
	}
	return nil
}
