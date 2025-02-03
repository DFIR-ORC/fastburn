package fastfind

import (
	"encoding/xml"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type FastFindConfig struct {
	XMLName    xml.Name `xml:"fastfind"`
	Version    string   `xml:"version,attr"`
	Filesystem struct {
		Location struct {
			Shadows string `xml:"shadows,attr"`
		} `xml:"location"`
		Yara []struct {
			Block      string `xml:"block,attr"`
			Overlap    string `xml:"overlap,attr"`
			Timeout    string `xml:"timeout,attr"`
			Source     string `xml:"source,attr"`
			ScanMethod string `xml:"scan_method,attr"`
		} `xml:"yara"`
		NtfsFind []struct {
			SizeLe   string `xml:"size_le,attr"`
			Header   string `xml:"header,attr"`
			YaraRule string `xml:"yara_rule,attr"`
			Size     string `xml:"size,attr"`
			SHA256   string `xml:"sha256,attr"`
			Name     string `xml:"name,attr"`
		} `xml:"ntfs_find"`
	} `xml:"filesystem"`
}

func (cfg FastFindConfig) String() string {
	return "Shadows: " + cfg.Filesystem.Location.Shadows + ", Yaras:" + fmt.Sprintf("%v", cfg.Filesystem.Yara) + ", Ntfs:" + fmt.Sprintf("%v", cfg.Filesystem.NtfsFind)
}

// ReadFastFindConfig read a FastFind configuration and returns a FastFindConfig struct
func ReadFastfindConfig(fname string) (*FastFindConfig, error) {
	log.Debug("Reading FastFind config from " + fname)
	cfg_data, err := os.ReadFile(fname)
	if err != nil {
		log.Error("Failed to read configuration from file '" + fname + "' :" + err.Error())
		return nil, err
	}
	config := new(FastFindConfig)
	err = xml.Unmarshal(cfg_data, config)
	if err != nil {
		log.Error("Failed to parse configuration from file '" + fname + "' :" + err.Error())
		return nil, err
	}
	log.Debug("FastFind config from '" + fname + "': " + config.String())
	return config, nil
}
