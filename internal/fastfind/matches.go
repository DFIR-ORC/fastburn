package fastfind

import (
	"encoding/json"
	"fmt"
)

type MatchType int

const (
	FilesystemMatchType = iota
	RegistryMatchType
	I30MatchType
)

func (d MatchType) String() string {
	return [...]string{"Filesystem", "Registry", "I30"}[d]
}

// FastFindMatch - a data structure to store the flattened result of a FastFind_result.xml row
type FastFindMatch struct {
	Kind                        MatchType
	Reason                      string
	Computer                    string
	ComputerOS                  string
	ComputerRole                string
	ORCVersion                  string
	VolumeID                    string
	SnapshotID                  string
	Fullname                    string
	AltFilename                 string
	ParentFrn                   string
	AltParentFrn                string
	RegKey                      string
	RegType                     string
	RegValue                    string
	Creation                    string
	LastModification            string
	LastEntryChange             string
	LastAccess                  string
	FilenameCreation            string
	FilenameLastModification    string
	FilenameLastEntryChange     string
	FilenameLastAccess          string
	AltFilenameCreation         string
	AltFilenameLastModification string
	AltFilenameLastEntryChange  string
	AltFilenameLastAccess       string
	Size                        uint64
	MD5                         string
	SHA1                        string
	SHA256                      string
	ArchiveName                 string
	Ignore                      bool
}

type FastFindMatchesList []*FastFindMatch

// FastFindComputer - a data structure to synthetise the informations related to a particular computer seen in FastFind_result.xml
type FastFindComputer struct {
	Hostname    string
	ArchiveName string
	OS          string
	Role        string
	ORCVersion  string
	NbMatches   uint
}

type FastFindComputersList []*FastFindComputer

func (m FastFindMatch) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ("Failed serialisation of Match: " + err.Error())
	}
	return string(b)
}

func (m *FastFindMatch) Print() {
	fmt.Printf(" - %s\n", m.Fullname)
	fmt.Printf("\t- Match: %s\n", m.Reason)
	fmt.Printf("\t- Computer: %s\n", m.Computer)
	fmt.Printf("\t- OS: %s\n", m.ComputerOS)
	fmt.Printf("\t- Role: %s\n", m.ComputerRole)
	fmt.Printf("\t- ORCVersion: %s\n", m.ORCVersion)
	fmt.Printf("\t- File creation: %s\n", m.Creation)
	fmt.Printf("\t- File last modification: %s\n", m.LastModification)
	fmt.Printf("\t- File last entry change: %s\n", m.LastEntryChange)
	fmt.Printf("\t- Filename creation: %s\n", m.FilenameCreation)
	fmt.Printf("\t- Filename last modification: %s\n", m.FilenameLastModification)
	fmt.Printf("\t- Filename last entry change: %s\n", m.FilenameLastEntryChange)
	fmt.Printf("\t- File Size: %v\n", m.Size)
	fmt.Printf("\t- MD5: %s\n", m.MD5)
	fmt.Printf("\t- SHA1: %s\n", m.SHA1)
	fmt.Printf("\t- SHA256: %s\n", m.SHA256)
	//TODO gerer le regtype
}

type FilterFunc func(*FastFindMatch) (bool, string)

// URI - returns a string describing a synthetic description of a match path
func (m *FastFindMatch) URI() string {

	var uri string

	switch m.Kind {
	case FilesystemMatchType:
		uri = fmt.Sprintf("file://%s:%s", m.Computer, m.Fullname)
	case RegistryMatchType:
		uri = fmt.Sprintf("reg://%s:%s", m.Computer, m.Fullname)
		if m.RegKey != "" {
			uri = uri + ":" + m.RegKey
			if m.RegValue != "" {
				uri = uri + ":" + m.RegValue + "<" + m.RegType + ">"
			}
		}
	default:
		uri = m.Fullname
	}

	return uri
}
