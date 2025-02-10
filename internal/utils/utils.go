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
		log.Warning(fmt.Sprintf("Failed to stats '%s': %s", filename, err.Error()))
		return false, err
	}
	if fi.Mode().IsDir() {
		return true, nil
	}
	return false, nil

}

func isArchive(filename string) (bool, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		log.Warning(fmt.Sprintf("Failed to stats '%s': %s", filename, err.Error()))
		return false, err
	}
	if fi.Mode().IsRegular() && strings.HasSuffix(filename, ".7z") {
		return true, nil
	}
	return false, nil
}

// ExpandArchiveFilepath - take a filelist, descend in directories, and returns a list of 7z files
func ExpandArchiveFilePaths(filenames []string) ([]string, error) {

	files := []string{}
	var err error
	for _, filename := range filenames {
		log.Trace(fmt.Sprintf("Examining : '%s'", filename))
		isdir, err := isDir(filename)
		if err != nil {
			log.Errorf("Failed to examine '%s'", filename)
			return files, err
		} else {
			if isdir {
				// directory, globing files
				log.Trace(fmt.Sprintf("Exploring directory : '%s'", filename))
				dirfiles := []string{}
				subdirs := []string{}
				err := filepath.Walk(filename,
					func(path string, info os.FileInfo, err error) error {
						if path != filename {
							log.Trace(fmt.Sprintf("Walking into '%s'", path))
							isdir, _ := isDir(path)
							isarc, _ := isArchive(path)
							if isdir {
								subdirs = append(subdirs, path)
							} else if isarc {
								log.Trace(fmt.Sprintf("Adding subdir archive %s", path))
								dirfiles = append(dirfiles, path)
							}
						}
						return nil
					})
				if err != nil {
					log.Warningf("Failed to examine directory '%s'", filename)
				}
				if len(subdirs) > 0 {
					subfiles, err := ExpandArchiveFilePaths(subdirs)
					if err != nil {
						log.Warning(fmt.Sprintf("Failed to process subdirs: %v", err))
						return nil, err
					}
					log.Trace("Appending dir parsing result")
					dirfiles = append(dirfiles, subfiles...)
				}
				files = append(files, dirfiles...)
			} else {
				isarc, err := isArchive(filename)
				if err != nil {
					log.Warningf("Failed to examine '%s'", filename)
				} else if isarc {
					log.Trace(fmt.Sprintf("Adding archive %s", filename))
					files = append(files, filename)
				} else {
					log.Trace(fmt.Sprintf("Skipping '%s': neither an archive nor a directory", filename))
				}
			}
		}
	}
	log.Trace(fmt.Sprintf("Examination of %v returned %d results", filenames, len(files)))
	files = Uniq(files)
	return files, err
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

func SetLogLevel(debug bool, trace bool) {

	log.SetReportCaller(false)
	log.SetLevel(log.InfoLevel)
	if debug {
		log.SetLevel(log.TraceLevel)
	}
	if trace {
		log.SetLevel(log.TraceLevel)
	}
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
