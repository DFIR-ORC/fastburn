package fastfind

import (
	"fmt"
	"io"
	"sort"

	net "github.com/THREATINT/go-net"
	log "github.com/sirupsen/logrus"
)

type FastFindMatchesStats struct {
	Fields map[string]map[string]uint
}

func getFieldsList() []string {
	flds := []string{
		"Total",
		"Computers-Role",
		"Matches-ComputerRole",
		"Computers-OS",
		"Matches-ComputerOS",
		"Computers-Domain",
		"Matches-Domain",
		"Matches-Reason",
		"Matches-Fullname",
		"Matches-MD5",
		"Matches-SHA1",
		"Matches-FileSize",
		"Matches-CreationDay",
		"Matches-ModificationDay",
		"Matches-CreationMonth",
		"Matches-ModificationMonth",
	}
	return flds
}

// truncateISOToDay truncate a string containing an ISO date to the day: YYYY-MM-DD
func truncateISOToDay(isoDate string) string {
	if len(isoDate) > 9 {
		return isoDate[:10]
	} else {
		log.Warn("Expected an ISO 8660 timestamp and got '" + isoDate + "'")
		return isoDate
	}
}

// truncateISOToMonth truncate a string containing an ISO date to the month: YYYY-MM
func truncateISOToMonth(isoDate string) string {
	if len(isoDate) > 6 {
		return isoDate[:7]
	} else {
		log.Warn("Expected an ISO 8660 timestamp and got '" + isoDate + "'")
		return isoDate
	}
}

// update the stats with a computer information
func (stats *FastFindMatchesStats) UpdateComputers(comp *FastFindComputer) {
	domain := net.DomainFromFqdn(comp.Hostname)

	stats.Fields["Total"]["Computers"]++

	stats.Fields["Computers-OS"][comp.OS]++
	stats.Fields["Computers-Role"][comp.Role]++
	stats.Fields["Computers-Domain"][domain]++
}

// update the stats with a match information
func (stats *FastFindMatchesStats) UpdateMatches(match *FastFindMatch) {

	strSize := fmt.Sprintf("%09d", match.Size)
	domain := net.DomainFromFqdn(match.Computer)

	stats.Fields["Total"]["Matches"]++

	stats.Fields["Matches-Domain"][domain]++
	stats.Fields["Matches-Reason"][match.Reason]++
	stats.Fields["Matches-ComputerOS"][match.ComputerOS]++
	stats.Fields["Matches-ComputerRole"][match.ComputerRole]++
	stats.Fields["Matches-FileSize"][strSize]++
	stats.Fields["Matches-Fullname"][match.Fullname]++
	stats.Fields["Matches-MD5"][match.MD5]++
	stats.Fields["Matches-SHA1"][match.SHA1]++
	stats.Fields["Matches-ModificationDay"][truncateISOToDay(match.LastModification)]++
	stats.Fields["Matches-ModificationMonth"][truncateISOToMonth(match.LastModification)]++
	stats.Fields["Matches-CreationDay"][truncateISOToDay(match.Creation)]++
	stats.Fields["Matches-CreationMonth"][truncateISOToMonth(match.Creation)]++
}

func (stats *FastFindMatchesStats) ToCSV(w io.Writer) error {

	fields := getFieldsList()
	separator := "----------------\n"

	io.WriteString(w, "Field,Value,Count\n")
	for _, key := range fields {
		err := stats.FieldToCSV(key, w)
		if err != nil {
			log.Errorf("failed to compute statistics for key [%s]", key)
			return err
		}
		io.WriteString(w, separator)
	}
	return nil
}

func (stats *FastFindMatchesStats) FieldToCSV(field string, w io.Writer) error {

	values := stats.Fields[field]

	fields := make([]string, 0, len(values))
	for key := range values {
		fields = append(fields, key)
	}
	sort.Strings(fields)

	for _, key := range fields {
		value := values[key]
		line := fmt.Sprintf("\"%v\",\"%v\",%v\n", field, key, value)
		io.WriteString(w, line)
	}
	return nil
}

func CreateStats() *FastFindMatchesStats {
	stats := FastFindMatchesStats{}
	stats.Fields = make(map[string]map[string]uint)

	for _, fld := range getFieldsList() {
		stats.Fields[fld] = make(map[string]uint)
	}
	return &stats
}
