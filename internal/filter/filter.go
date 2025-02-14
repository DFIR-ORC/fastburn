package filter

/*
 *
 * Filtering of results with blacklists and whitelists
 *
 */
import (
	"strings"

	fbn "fastburn/internal/fastfind"

	log "github.com/sirupsen/logrus"
)

func (c *criteriaEntry) Match(sha256 string, sha1 string, md5 string, fullname string) bool {
	log.Tracef("Trying to match sha256:'%s' sha1:'%s' md5:'%s' filename:'%s'", sha256, sha1, md5, fullname)
	log.Tracef("With sha256:'%s' sha1:'%s' md5:'%s' filere:[%s]", sha256, sha1, md5, c.FileRE)
	if c.SHA256 != "" && strings.EqualFold(c.SHA256, sha256) {
		return true
	} else {
		log.Trace("no sha256 to match")
	}
	if c.SHA1 != "" && strings.EqualFold(c.SHA1, sha1) {
		return true
	} else {
		log.Trace("no sha1 to match")
	}

	if c.MD5 != "" && strings.EqualFold(c.MD5, md5) {
		return true
	} else {
		log.Trace("no md5 to match")
	}
	if c.Regexp != nil {
		matched := c.Regexp.MatchString(fullname)
		if matched {
			return true
		} else {
			log.Tracef("path '%s' does not match regexp [%s]", fullname, c.Regexp)
		}
	} else {
		log.Trace("no regexp to match")
	}
	return false
}

func isListed(criterias *[]criteriaEntry, m *fbn.FastFindMatch) (bool, string) {
	log.Trace("Matching " + m.Fullname)

	fullname := strings.ToLower(m.Fullname)
	for _, c := range *criterias {
		if c.Match(m.SHA256, m.SHA1, m.MD5, fullname) {
			log.Trace("!Match for '" + fullname + "'")
			return true, c.Description
		}
	}
	log.Trace("!Nomatch for '" + fullname + "'")
	return false, ""
}

// Match - method returning whether a FastFind Result looks like a specialy interresting file (for IOC loaded in filter_data.go)
func (f *Filter) IsWhitelisted(m *fbn.FastFindMatch) (bool, string) {
	log.Trace("Matching in whitelist: " + m.Fullname + " " + m.SHA256)
	return isListed(&f.whitelistCriteria, m)
}

// Match - method returning whether a FastFind Result looks like a specialy interresting file (for IOC loaded in filter_data.go)
func (f *Filter) IsBlacklisted(m *fbn.FastFindMatch) (bool, string) {
	log.Trace("Matching in blacklist: " + m.Fullname + " " + m.SHA256)
	return isListed(&f.blacklistCriteria, m)
}
