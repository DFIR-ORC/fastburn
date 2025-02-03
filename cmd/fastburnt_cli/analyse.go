package main

/*
fastburnt_cli command line tool to parse FastFind results

BUILD

  export GOPATH="$GOPATH:$(pwd)"

  go get  "github.com/kjk/lzmadec"
  go get "github.com/sirupsen/logrus"

  go build fastburnt/cmd/fastburnt_cli

RUN

Attention: requires 7z binary in the path named as 7z (or 7z.exe on Windows)

  ./fastburnt_cli <7z archive1 ... n>

**/

import (
	"fmt"

	fastfound "dfir-orc/fastburnt/internal/fastfind"
	"dfir-orc/fastburnt/internal/filter"
	"dfir-orc/fastburnt/internal/utils"

	log "github.com/sirupsen/logrus"
)

// parseFiles - process the command line to list files and parse them to the in-memory data structures
func parseFiles(args []string) ([]string, *fastfound.FastFindMatchesList, *fastfound.FastFindComputersList, *fastfound.FastFindMatchesStats, error) {

	matches := make(fastfound.FastFindMatchesList, 0)
	computers := make(fastfound.FastFindComputersList, 0)

	stats := fastfound.CreateStats()
	log.Debugf("Processing file list: %v", args)
	files, err := utils.ExpandArchiveFilePaths(utils.Uniq(args))
	if err != nil {
		log.Errorf("Failed to expand paths: %v", err)
		return nil, nil, nil, nil, err
	}
	log.Info(fmt.Sprintf("%d files to process", len(files)))
	log.Debug(fmt.Sprintf("Processing %d files", len(files)))

	for _, fname := range files {
		matches, computers, err = fastfound.ProcessFile(fname, matches, computers)
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
func analyseData(matches *fastfound.FastFindMatchesList, stats *fastfound.FastFindMatchesStats, postfilter *filter.Filter) (*fastfound.Timeline, error) {
	log.Debug(fmt.Sprintf("Post-processing %v results", len(*matches)))
	var blacklistCount uint64 = 0
	var whitelistCount uint64 = 0
	timeline := fastfound.InitTL()

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
	matches *fastfound.FastFindMatchesList, computers *fastfound.FastFindComputersList, stats *fastfound.FastFindMatchesStats, timeline *fastfound.Timeline,
	postfilter *filter.Filter) error {

	// CSV export
	log.Debug("Exporting results")

	err := fastfound.ExportComputersToCSV(csv_computers_fname, computers)
	if err != nil {
		log.Warning("Export to '" + csv_computers_fname + "' failed: " + err.Error())
	} else {
		log.Info("Computers exported to '" + csv_computers_fname + "'")
	}

	if len(*matches) > 0 {

		// is infected ?
		f := func(m *fastfound.FastFindMatch) (bool, string) {
			return postfilter.IsWhitelisted(m)
		}
		// exporting
		err = fastfound.ExportMatchesToCSV(csv_matches_fname, matches, f)
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
	err = fastfound.ExportStatsToCSV(csv_stats_fname, stats)
	if err != nil {
		log.Warning("Stat report export to '" + csv_stats_fname + "' failed: " + err.Error())
	} else {
		log.Info("Stat report exported to '" + csv_stats_fname + "'")
	}

	// exporting timeline
	if timeline_fname != "" {
		fastfound.ExportTimelineToCSV(timeline_fname, timeline)
		if err != nil {
			log.Warning("Timeline export to '" + timeline_fname + "' failed: " + err.Error())
		} else {
			log.Info("Timeline exported to '" + timeline_fname + "'")
		}
	}
	return nil
}
