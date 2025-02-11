package fastfind

import (
	"embed"
	"os"

	thtml "html/template"

	log "github.com/sirupsen/logrus"
)

//go:embed  templates/*.html
var templateFS embed.FS

type htmlData struct {
	Title     string
	Matches   *FastFindMatchesList
	Stats     *FastFindMatchesStats
	Computers *FastFindComputersList
	Timeline  *FastFindTimeline
}

// ExportMatchesToCSV - Export a list of matches to a HTML file appending filter columens from the provided functions
func ExportToHTML(filename string, matches *FastFindMatchesList, stats *FastFindMatchesStats, computers *FastFindComputersList, timeline *FastFindTimeline) error {
	log.Debug("Exporting HTML to " + filename)
	fout, err := os.Create(filename)
	if err != nil {
		log.Trace("Export to " + filename + " :" + err.Error())
		return err
	}
	defer fout.Close()

	t, err := thtml.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		log.Trace("HTML export failed to compile text template " + err.Error())
		return err
	}

	err = t.ExecuteTemplate(
		fout, "frame.html",
		htmlData{
			Title: "FastFind results",
			Stats: stats, Matches: matches, Computers: computers, Timeline: timeline})
	if err != nil {
		log.Trace("HTML export : failed to fill HTML template " + err.Error())
		return err
	}

	log.Trace("HTML Export to " + filename + " done")
	return nil
}
