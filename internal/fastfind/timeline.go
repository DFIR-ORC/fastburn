package fastfind

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"fastburn/internal/utils"

	log "github.com/sirupsen/logrus"
)

const PathSeparator = "\\"

type TLEvent struct {
	Timestamp string
	Match     *FastFindMatch
	SI_M      bool
	SI_A      bool
	SI_C      bool
	SI_B      bool
	FN_M      bool
	FN_A      bool
	FN_C      bool
	FN_B      bool
}

var csvHeaders []string = []string{
	"Timestamp",
	"SI_MACB",
	"FN_MACB",
	"ComputerName",
	"File",
	"ParentName",
	"FullName",
	"Extension",
	"SizeInBytes",
	"CreationDate",
	"LastModificationDate",
	"LastAccessDate",
	"LastAttrChangeDate",
	"FileNameCreationDate",
	"FileNameLastModificationDate",
	"FileNameLastAccessDate",
	"FileNameLastAttrModificationDate",
	"MD5",
	"SHA1",
	"SHA256",
	"Reason",
	"ArchiveName",
}

const csvQuote string = "\""
const csvSep string = ";"

type FastFindTimeline struct{ Events map[string][]TLEvent }

func writeToCSV(w io.Writer, parts []string) (n int, err error) {
	strLine := csvQuote +
		strings.Join(parts, csvQuote+csvSep+csvQuote) + csvQuote + "\n"
	return io.WriteString(w, strLine)
}

func (e *TLEvent) SetMACB(timestamp string) {
	e.Timestamp = timestamp

	if timestamp == e.Match.Creation {
		e.SI_M = true
		e.SI_A = true
		e.SI_C = true
		e.SI_B = true
	}
	if timestamp == e.Match.LastModification {
		e.SI_M = true
	}
	if timestamp == e.Match.LastEntryChange {
		e.SI_C = true
	}
	if timestamp == e.Match.LastAccess {
		e.SI_A = true
	}
	if timestamp == e.Match.FilenameCreation {
		e.FN_M = true
		e.FN_A = true
		e.FN_C = true
		e.FN_B = true
	}
	if timestamp == e.Match.FilenameLastModification {
		e.FN_M = true
	}
	if timestamp == e.Match.FilenameLastEntryChange {
		e.FN_C = true
	}
	if timestamp == e.Match.FilenameLastAccess {
		e.FN_A = true
	}
	// TODO: check EntryChange vs Modification
	// TODO: account for Filename et AltFilename timestamps
}

func InitTL() *FastFindTimeline {
	tl := FastFindTimeline{}
	tl.Events = make(map[string][]TLEvent)
	return &tl
}

// small device to simplify timestamps deduplication
type timestampMatchMerger struct {
	Timestamps map[string]int
	Match      *FastFindMatch
}

func (merger *timestampMatchMerger) addStamp(timestamp string, stampname string) {
	if timestamp != "" {
		(merger.Timestamps)[timestamp]++
	} else {
		log.Debug("Empty " + stampname + " time for match " + merger.Match.Fullname)
	}
}

func (tl *FastFindTimeline) Register(match *FastFindMatch) {

	// Building the list of timestamps
	merger := timestampMatchMerger{Match: match}
	merger.Timestamps = make(map[string]int)

	merger.addStamp(match.Creation, "Creation")
	merger.addStamp(match.LastModification, "LastModification")
	merger.addStamp(match.LastEntryChange, "LastEntryChange")
	merger.addStamp(match.LastAccess, "LastAccess")

	merger.addStamp(match.FilenameCreation, "FilenameCreation")
	merger.addStamp(match.FilenameLastModification, "FilenameLastModification")
	merger.addStamp(match.FilenameLastEntryChange, "FilenameLastEntryChange")
	merger.addStamp(match.FilenameLastAccess, "FilenameLastAccess")

	merger.addStamp(match.AltFilenameCreation, "AltFilenameCreation")
	merger.addStamp(match.AltFilenameLastModification, "AltFilenameLastModification")
	merger.addStamp(match.AltFilenameLastEntryChange, "AltFilenameLastEntryChange")
	merger.addStamp(match.AltFilenameLastAccess, "AltFilenameLastAccess")

	if len(merger.Timestamps) == 0 {
		log.Warn("Missing timestamps for match " + match.Fullname)
	}

	// Foreach different timestamp
	// create a new event
	for stamp := range merger.Timestamps {
		event := TLEvent{Match: match}
		event.SetMACB(stamp)
		if tl.Events[stamp] == nil {
			tl.Events[stamp] = []TLEvent{}
		}
		tl.Events[stamp] = append(tl.Events[stamp], event)
	}
}

func macbToString(M bool, A bool, C bool, B bool) string {
	macb := ""
	if M {
		macb += "M"
	} else {
		macb += "."
	}
	if A {
		macb += "A"
	} else {
		macb += "."
	}
	if C {
		macb += "C"
	} else {
		macb += "."
	}
	if B {
		macb += "B"
	} else {
		macb += "."
	}
	return macb
}

func (e *TLEvent) ToCSV(w io.Writer) error {

	si_macb := macbToString(e.SI_M, e.SI_A, e.SI_C, e.SI_B)
	fn_macb := macbToString(e.FN_M, e.FN_A, e.FN_C, e.FN_B)

	fakeName := filepath.ToSlash(e.Match.Fullname)

	ext := utils.WinExt(fakeName)
	dir := utils.WinDir(fakeName)
	basename := utils.WinBase(fakeName)

	strSze := fmt.Sprintf("%d", e.Match.Size)

	line := []string{
		e.Timestamp,
		si_macb,
		fn_macb,
		e.Match.Computer,
		basename,
		dir,
		e.Match.Fullname,
		ext,
		strSze,
		e.Match.Creation,
		e.Match.LastModification,
		e.Match.LastAccess,
		e.Match.LastEntryChange,
		e.Match.FilenameCreation,
		e.Match.FilenameLastModification,
		e.Match.FilenameLastAccess,
		e.Match.AltFilenameLastEntryChange,
		e.Match.MD5,
		e.Match.SHA1,
		e.Match.SHA256,
		e.Match.Reason,
		e.Match.ArchiveName,
	}
	writeToCSV(w, line)

	return nil
}

func (tl *FastFindTimeline) ToCSV(w io.Writer) {

	writeToCSV(w, csvHeaders)

	// sorting timestamps
	stamps := make([]string, 0, len(tl.Events))
	for stamp := range tl.Events {
		stamps = append(stamps, stamp)
	}
	sort.Strings(stamps)

	// foreach timestamp
	for _, stamp := range stamps {
		if stamp != "" {
			evts := tl.Events[stamp]
			// dump each event
			for i := range evts {
				evt := evts[i]
				evt.ToCSV(w)
			}
		}
	}
}
