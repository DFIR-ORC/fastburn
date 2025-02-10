package fastfind

import (
	"bytes"
	"os"

	thtml "html/template"
	ttxt "text/template"

	log "github.com/sirupsen/logrus"
)

var tmplHTMLFrame = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>{{ .Title }}</title>
		<style>
		body {
			font-family: helvetica, sans-serif;
		}
    	table.matches {
      		border-collapse: collapse;
    	}
    	.matches tr:nth-child(even) {
      		background-color: #EAEAFA;
    	}
		.matches td{
			border: 1px solid black;
		}
		.matches th{
      		background-color: black;
			border: 1px solid #EAEAFA;
      	color: #EAEAFA;
			border: 1px black;
		}
		.matches {
			width: 100%;
			margin: auto;
		}
		</style>
	</head>
	<body>
		<h2>Statistics</h2>
		{{ .StatsTable }}
		<hr>
		<h2>Matches</h2>
		{{ .MatchesTable }}
	 </body>
 </html>
`

var tmplHTMLMatches = `
		<table class="matches">
			<thead>
				<tr>				
					<th>Computer</th>
					<th>ComputerRole</th>
					<th>ComputerOS</th>
					<th>ORCVersion</th>
					<th>MatchType</th>
					<th>Reason</th>
					<th>Filename
					<th>AltName</th>
					<th>RegKey</th>
					<th>RegType</th>
					<th>RegValue</th>
					<th>FileSize</th>
					<th>MD5</th>
					<th>SHA1</th>
					<th>SHA256</th>
					<th>FileCreation<th>
					<th>FileLastModification</th>
					<th>FileLastEntryChange</th>
					<th>FileLastAccess</th>
					<th>FilenameCreation</th>
					<th>FilenameLastModification</th>
					<th>FilenameLastEntryChange</th>
					<th>FilenameLastAccess</th>
					<th>AltFilenameCreation</th>
					<th>AltFilenameLastModification</th>
					<th>AltFilenameLastEntryChange</th>
					<th>AltFilenameLastAccess</th>
					<th>VolumeID</th>
					<th>SnapshotID</th>
					<th>ArchiveName</th>
				</tr>
			</thead>
			<tbody>
			{{ range .Matches }}
			 	{{ if .Ignore }}<tr class="ignore">	{{ else }}<tr class="match">{{end}}
					<td>{{ .Computer }}</td>
					<td>{{ .ComputerRole }}</td>
					<td>{{ .ComputerOS }}</td>
					<td>{{ .ORCVersion }}</td>
					<td>{{ printf "%s" .Kind }}</td>
					<td>{{ .Reason }}</td>
					<td>{{ .Fullname }}</td>
					<td>{{ .AltFilename }}</td>
					<td>{{ .RegKey }}</td>
					<td>{{ .RegType }}</td>
					<td>{{ .RegValue }}</td>
					<td>{{ .Size }}</td>
					<td>{{ .MD5 }}</td>
					<td>{{ .SHA1 }}</td>
					<td>{{ .SHA256 }}</td>
					<td>{{ .Creation }}<td>
					<td>{{ .LastModification }}</td>
					<td>{{ .LastEntryChange }}</td>
					<td>{{ .LastAccess }}</td>
					<td>{{ .FilenameCreation }}</td>
					<td>{{ .FilenameLastModification }}</td>
					<td>{{ .FilenameLastEntryChange }}</td>
					<td>{{ .FilenameLastAccess }}</td>
					<td>{{ .AltFilenameCreation }}</td>
					<td>{{ .AltFilenameLastModification }}</td>
					<td>{{ .AltFilenameLastEntryChange }}</td>
					<td>{{ .AltFilenameLastAccess }}</td>
					<td>{{ .VolumeID }}</td>
					<td>{{ .SnapshotID }}</td>
					<td>{{ .ArchiveName }}</td>
				</tr>
			{{ end }}
			</tbody>
		</table>
`

var tmplHTMLStats = `
	<table class="matches">
			<thead>
				<tr>				
					<th>Category</th>
					<th>Value</th>
					<th>Count</th>
				</tr>
			</thead>
			<tbody>
			{{ range $key,$value := .Stats.Fields }}
			 	{{ range $innerKey, $innerValue := $value }}
			 	<tr>
					<td>{{ $key }}</td>
					<td>{{ $innerKey }}</td>
					<td>{{ $innerValue }}</td>
				</tr>
				{{ end }}
			{{ end }}
			</tbody>
		</table>
`

type frameData struct {
	Title        string
	StatsTable   string
	MatchesTable string
}

type matchesData struct {
	Matches *FastFindMatchesList
}

type statsData struct {
	Stats *FastFindMatchesStats
}

// ExportMatchesToCSV - Export a list of matches to a HTML file appending filter columens from the provided functions
func ExportToHTML(filename string, matches *FastFindMatchesList, stats *FastFindMatchesStats) error {
	log.Debug("Exporting HTML to " + filename)
	fout, err := os.Create(filename)
	if err != nil {
		log.Trace("Export to " + filename + " :" + err.Error())
		return err
	}
	defer fout.Close()

	t, err := ttxt.New("frame").Parse(tmplHTMLFrame)
	if err != nil {
		log.Trace("HTML export failed to compile text template " + err.Error())
		return err
	}

	statsHTML, err := GenerateStatsHTML(stats)
	if err != nil {
		log.Trace("HTML export failed to generate stats HTML " + err.Error())
		return err
	}

	log.Trace(statsHTML)

	matchesHTML, err := GenerateMatchesHTML(matches)
	if err != nil {
		log.Trace("HTML export failed to generate matches HTML " + err.Error())
		return err
	}

	err = t.Execute(
		fout,
		frameData{Title: "FastFind results", StatsTable: statsHTML, MatchesTable: matchesHTML})
	if err != nil {
		log.Trace("HTML export : failed to fill HTML template " + err.Error())
		return err
	}

	log.Trace("HTML Export to " + filename + " done")
	return nil
}

// ExportMatchesToCSV - Export a list of matches to a HTML file appending filter columens from the provided functions
func GenerateMatchesHTML(matches *FastFindMatchesList) (string, error) {
	log.Debug("Exporting HTML matches")

	t, err := thtml.New("matches").Parse(tmplHTMLMatches)
	if err != nil {
		log.Trace("HTML export failed to compile HTML template for matches" + err.Error())
		return "", err
	}

	var rawHTML bytes.Buffer
	err = t.Execute(&rawHTML, matchesData{Matches: matches})
	if err != nil {
		log.Trace("HTML matches export : failed to fill HTML template " + err.Error())
		return "", err
	}

	log.Trace("HTML export of matches done")
	return rawHTML.String(), nil
}

/*
// ExportComputersToCSV - Export a list of matches to a CSV file appending filter columens from the provided functions

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
*/

func GenerateStatsHTML(stats *FastFindMatchesStats) (string, error) {
	log.Debug("Exporting HTML stats")

	t, err := thtml.New("stats").Parse(tmplHTMLStats)
	if err != nil {
		log.Trace("HTML export failed to compile HTML template " + err.Error())
		return "", err
	}

	var rawHTML bytes.Buffer
	err = t.Execute(&rawHTML, statsData{Stats: stats})
	if err != nil {
		log.Trace("HTML stats export : failed to fill HTML template " + err.Error())
		return "", err
	}

	log.Trace("Stats HTML Export  done")

	return rawHTML.String(), nil
}

/*
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
*/
