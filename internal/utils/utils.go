package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func isDir(filename string) (bool, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		log.Warningf("Failed to stats '%s': %s", filename, err.Error())
		return false, err
	}
	if fi.Mode().IsDir() {
		return true, nil
	}
	return false, nil

}

func isDataFileWithExtension(filename string, extensions []string) (bool, string, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		log.Warningf("Failed to stats '%s': %s", filename, err.Error())
		return false, "", err
	}
	if fi.Mode().IsRegular() {
		for _, ext := range extensions {
			if strings.HasSuffix(strings.ToUpper(filename), strings.ToUpper(ext)) {
				return true, ext, nil
			}
		}
	}
	return false, "", nil
}

func isArchive(filename string) (bool, string, error) {
	extensions := []string{"7z", "p7b"}
	return isDataFileWithExtension(filename, extensions)
}

// ExpandArchiveFilepath - take a filelist, descend in directories, and returns a list of 7z archives and a list of p7b containers
func ExpandArchiveFilePaths(filenames []string) ([]string, []string, error) {

	archives := []string{}
	containers := []string{}
	var err error
	for _, filename := range filenames {
		log.Tracef("Examining : '%s'", filename)
		isdir, err := isDir(filename)
		if err != nil {
			log.Errorf("Failed to examine '%s'", filename)
			return archives, containers, err
		} else {

			if isdir {
				// directory, globing files
				log.Tracef("Exploring directory : '%s'", filename)
				dirArchives := []string{}
				dirContainers := []string{}
				subdirs := []string{}
				err := filepath.Walk(filename,
					func(path string, info os.FileInfo, err error) error {
						if path != filename {
							log.Tracef("Walking into '%s'", path)
							isDir, _ := isDir(path)
							isArc, ext, err := isArchive(path)
							if err != nil {
								log.Warningf("Failed to examine '%s'", filename)
							} else if isDir {
								subdirs = append(subdirs, path)
							} else if isArc {
								switch ext {
								case "7z":
									log.Tracef("Adding subdir archive '%s'", path)
									dirArchives = append(dirArchives, path)
								case "p7b":
									log.Tracef("Adding subdir encrypted archive '%s'", path)
									dirContainers = append(dirContainers, path)
								default:
									log.Tracef("ignoring unsupported archive file type '%s'", path)
								}

							} else {
								log.Tracef("ignoring non archive file '%s'", path)
							}
						}
						return nil
					})
				if err != nil {
					log.Warningf("Failed to examine directory '%s'", filename)
				}
				if len(subdirs) > 0 {
					subArcs, subEncArcs, err := ExpandArchiveFilePaths(subdirs)
					if err != nil {
						log.Warningf("Failed to process subdirs: %v", err)
						return nil, nil, err
					}
					log.Trace("Appending dir parsing result")
					dirArchives = append(dirArchives, subArcs...)
					dirContainers = append(dirContainers, subEncArcs...)
				}
				// append recursion results
				archives = append(archives, dirArchives...)
				containers = append(containers, dirContainers...)
			} else {
				isArc, ext, err := isArchive(filename)
				if err != nil {
					log.Warningf("Failed to examine '%s'", filename)
				} else if isArc {

					switch ext {
					case "7z":
						log.Tracef("Adding subdir archive '%s'", filename)
						archives = append(archives, filename)
					case "p7b":
						log.Tracef("Adding subdir encrypted archive '%s'", filename)
						containers = append(containers, filename)
					default:
						log.Tracef("ignoring unsupported archive file type '%s'", filename)
					}

				} else {
					log.Tracef("Skipping '%s': neither an archive nor a directory", filename)
				}
			}
		}
	}
	log.Tracef("Examination of %v returned %d results", filenames, len(archives))
	archives = Uniq(archives)
	containers = Uniq(containers)
	return archives, containers, err
}

func Uniq(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func PrintAndLog(level log.Level, format string, args ...interface{}) {

	message := fmt.Sprintf(format, args...)

	switch level {
	case log.DebugLevel:
		log.Debug(message)
	case log.InfoLevel:
		log.Info(message)
	case log.WarnLevel:
		log.Warn(message)
	case log.ErrorLevel:
		log.Error(message)
	case log.FatalLevel:
		log.Fatal(message)
	case log.PanicLevel:
		log.Panic(message)
	default:
	}

	fmt.Println(message)
}

func SetLogLevel(info bool, debug bool, trace bool) {

	log.SetReportCaller(false)

	if trace {
		log.SetLevel(log.TraceLevel)
	} else if debug {
		log.SetLevel(log.DebugLevel)
	} else if info {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
	log.Debugf("Log level is %s", log.GetLevel())
}

/*
 * ReScanString
 */
func ReScanStrings(lines string, pattern string) ([]string, error) {
	// regexp compilation
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Errorf("invalid regexp '%s': %v", pattern, err)
		return nil, err
	}

	var results []string
	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		m := re.FindStringSubmatch(line)
		if len(m) > 1 {
			results = append(results, m[1])
		}
	}

	if err := scanner.Err(); err != nil {
		log.Errorf("error parsing string: %v", err)
		return nil, err
	}

	return results, nil
}

// Build a comparable version of a software version string
// From <https://stackoverflow.com/questions/18409373/how-to-compare-two-version-number-strings-in-golang>
func VersionOrdinal(version string) string {
	// ISO/IEC 14651:2011
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			panic("VersionOrdinal: invalid version")
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo)
}
