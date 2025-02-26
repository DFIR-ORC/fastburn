/*
 * FastFind results parsing
 * Based on <https://dfir-orc.github.io/FastFind.html>
 */
package fastfind

type FastFind_NTFSFilename struct {
	Fullname         string `xml:"fullname,attr"`
	Parentfrn        string `xml:"parentfrn,attr"`
	Creation         string `xml:"creation,attr"`
	LastModification string `xml:"lastmodification,attr"`
	LastAccess       string `xml:"lastaccess,attr"`
	LastEntryChange  string `xml:"lastentrychange,attr"`
}

type FastFind_StandardInformation struct {
	Creation         string `xml:"creation,attr"`
	LastModification string `xml:"lastmodification,attr"`
	LastAccess       string `xml:"lastaccess,attr"`
	LastEntryChange  string `xml:"lastentrychange,attr"`
	Attributes       string `xml:"attributes,attr"`
}

type FastFind_FileData struct {
	FileSize uint64 `xml:"filesize,attr"`
	MD5      string `xml:"MD5,attr"`
	SHA1     string `xml:"SHA1,attr"`
	SHA256   string `xml:"SHA256,attr"`
}

type FastFind_FileMatchRecord struct {
	Frn                 string                       `xml:"frn,attr"`
	VolumeID            string                       `xml:"volume_id,attr"`
	SnapshotID          string                       `xml:"snapshot_id,attr"`
	StandardInformation FastFind_StandardInformation `xml:"standardinformation"`
	Filenames           []FastFind_NTFSFilename      `xml:"filename"`
	Data                FastFind_FileData            `xml:"data"`
	I30s                []FastFind_NTFSFilename      `xml:"i30"`
}

type FastFind_FileMatch struct {
	Description string                   `xml:"description,attr"`
	Record      FastFind_FileMatchRecord `xml:"record"`
}

// FastFind_RegValue - a registry key description
type FastFind_RegValue struct {
	Key             string `xml:"key,attr"`
	Value           string `xml:"value,attr"`
	Type            string `xml:"type,attr"`
	LastModifiedKey string `xml:"lastmodified_key,attr"`
	DataSize        uint64 `xml:"data_size,attr"`
}

// FastFind_Hive - a registry hive description, potentially containing matches
type FastFind_Hive struct {
	VolumeID   string `xml:"volume_id,attr"`
	SnapshotID string `xml:"snapshot_id,attr"`
	HivePath   string `xml:"hive_path,attr"`
	RegMatches []struct {
		Description string              `xml:"description,attr"`
		Values      []FastFind_RegValue `xml:"value"`
	}
}

/*

auld:
	Registry struct {
		Hive []struct {
			VolumeID   string `xml:"volume_id,attr"`
			SnapshotID string `xml:"snapshot_id,attr"`
			HivePath   string `xml:"hive_path,attr"`
			RegMatches []struct {
				Description string `xml:"description,attr"`
				Values      []struct {
					Key             string `xml:"key,attr"`
					SubkeysCount    uint   `xml:"subkeys_count,attr"`
					ValuesCount     uint   `xml:"values_count,attr"`
					Value           string `xml:"value,attr"`
					Type            string `xml:"type,attr"`
					Size            uint64 `xml:"size,attr"`
					LastmodifiedKey string `xml:"lastmodified_key,attr"`
				} `xml:"key"`
			} `xml:"regfind_match"`
		} `xml:"hive"`
	} `xml:"registry"`


*/

// FastFind_ObjectValue - a system object description
type FastFind_ObjectValue struct {
	Type             string `xml:"type,attr"`
	Name             string `xml:"name,attr"`
	Path             string `xml:"path,attr"`
	LinkTarget       string `xml:"link_target,attr"`
	LinkCreationtime string `xml:"link_creationtime,attr"`
}

// FastFindResult - Structure to load XML from FastFind_result.xml for ORC versions >10.2.2
type FastFindResultNg struct {
	//XMLName    xml.Name `xml:"fast_find"`
	Computer   string `xml:"computer,attr"`
	OS         string `xml:"os,attr"`
	Role       string `xml:"role,attr"`
	Filesystem struct {
		FSMatches []FastFind_FileMatch `xml:"filefind_match"`
	} `xml:"filesystem"`
	Registry struct {
		Hives []struct {
			VolumeID       string `xml:"volume_id,attr"`
			SnapshotID     string `xml:"snapshot_id,attr"`
			HivePath       string `xml:"hive_path,attr"`
			RegfindMatches []struct {
				Description string `xml:"description,attr"`
				Values      []struct {
					Key             string `xml:"key,attr"`
					Value           string `xml:"value,attr"`
					Type            string `xml:"type,attr"`
					LastmodifiedKey string `xml:"lastmodified_key,attr"`
					DataSize        uint64 `xml:"data_size,attr"`
				} `xml:"value"`
				Keys []struct {
					Key             string `xml:"key,attr"`
					SubkeysCount    uint   `xml:"subkeys_count,attr"`
					ValuesCount     uint   `xml:"values_count,attr"`
					LastmodifiedKey string `xml:"lastmodified_key,attr"`
				} `xml:"key"`
			} `xml:"regfind_match"`
		} `xml:"hive"`
	} `xml:"registry"`
	Object struct {
		ObjectMatches []FastFind_FileMatch `xml:"object_match"`
	} `xml:"object"`
}

// FastFindResult - Structure to load XML from FastFind_result.xml for ORC versions <10.2.2
type FastFindResultLegacy struct {
	//    XMLName   xml.Name `xml:"fast_find"`
	Computer  string               `xml:"computer,attr"`
	OS        string               `xml:"os,attr"`
	Role      string               `xml:"role,attr"`
	FSMatches []FastFind_FileMatch `xml:"filesystem"`
	Registry  []struct {
		VolumeID   string `xml:"volume_id,attr"`
		SnapshotID string `xml:"snapshot_id,attr"`
		HivePath   string `xml:"hive_path,attr"`
		RegMatch   struct {
			Description string            `xml:"description,attr"`
			Value       FastFind_RegValue `xml:"value"`
		} `xml:"regfind_match"`
	} `xml:"registry"`
}

var resultsFname = "FastFind_result.xml"
