package main

/*

fastburn command line tool to parse FastFind results

**/

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"fastburn/internal/fastfind"
	"fastburn/internal/filter"
	"fastburn/internal/utils"

	_ "fastburn/cmd/fastburn/rsrc"

	log "github.com/sirupsen/logrus"
)

func Version() string {
	var (
		time     string
		revision string
		modified bool
	)

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.time":
				time = s.Value
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	if modified {
		return fmt.Sprintf("%s-%s-dirty", time, revision)
	}
	return fmt.Sprintf("%s-%s", time, revision)
}

func PrintUsage() {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	os.Exit(2)
}

func main() {
	var whiteFilename string
	var blackFilename string
	var outputFlag string
	var computersFlag string
	var statsFlag string
	var timelineFlag string
	var keyFlag string

	flag.Usage = PrintUsage

	infoFlag := flag.Bool("info", false, "Enable info verbosity level")
	debugFlag := flag.Bool("debug", false, "Enable debug verbosity level")
	traceFlag := flag.Bool("trace", false, "Enable trace verbosity level")
	versionFlag := flag.Bool("version", false, "Show version and exit")

	flag.StringVar(&whiteFilename, "whitelist", "", "Specify a CSV file containing flags to highligth in the results")
	flag.StringVar(&blackFilename, "blacklist", "", "Specify a CSV file containing flags to suppress from the results")
	flag.StringVar(&outputFlag, "output", "", "Specify output filename")
	flag.StringVar(&computersFlag, "computers", "", "Specify computers listing filename")
	flag.StringVar(&statsFlag, "stats", "", "Specify statistics filename")
	flag.StringVar(&timelineFlag, "timeline", "", "Specify a filename for timeline output")
	flag.StringVar(&keyFlag, "key", "", "Private key for decryption of encrypted archive filename")

	htmlFlag := flag.Bool("html", false, "Generate an HTML output")

	flag.Parse()

	args := flag.Args()

	if *versionFlag {
		version := Version()
		fmt.Printf("Fastburnt - version: %s\n", version)
		os.Exit(0)
	}

	if len(args) == 0 {
		PrintUsage()
		os.Exit(0)
	}

	utils.SetLogLevel(*infoFlag, *debugFlag, *traceFlag)

	var err error

	var postfilter filter.Filter
	err = postfilter.LoadLists(whiteFilename, blackFilename)
	if err != nil {
		log.Fatalf("Failed to load qualification lists: %v", err)
	}

	// look for archives and parse them to memory
	files, matches, computers, stats, err := parseFiles(keyFlag, args)
	if err != nil {
		log.Errorf("Failed to process files: %v", err)
		os.Exit(1)
	}

	// lets process collected data
	timeline, err := analyseData(matches, stats, &postfilter)
	if err != nil {
		log.Errorf("Failed to process files: %v", err)
		os.Exit(1)
	}

	// generating exports filenames
	prefix := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "_")
	csv_matches_fname := prefix + "-fastburn_matches.csv"
	csv_computers_fname := prefix + "-fastburn_computers.csv"
	csv_stats_fname := prefix + "-fastburn_stats.csv"
	timeline_fname := prefix + "-fastburn_timeline.csv"

	if timelineFlag != "" {
		timeline_fname = timelineFlag
	}

	if outputFlag != "" {
		csv_matches_fname = outputFlag
	}

	if computersFlag != "" {
		csv_computers_fname = computersFlag
	}

	if statsFlag != "" {
		csv_stats_fname = statsFlag
	}

	if *htmlFlag {
		html_fname := prefix + "-fastburn_matches.html"
		err = fastfind.ExportToHTML(html_fname, matches, stats, computers, timeline)
		if err != nil {
			log.Errorf("Failed to export results: %v", err)
			os.Exit(1)
		}
	}

	err = saveResults(csv_matches_fname, csv_computers_fname, csv_stats_fname, timeline_fname,
		files, matches, computers, stats, timeline, &postfilter)
	if err != nil {
		log.Errorf("Failed to export results: %v", err)
		os.Exit(1)
	}

}

//eof
