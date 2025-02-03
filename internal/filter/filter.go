package filter

import (
	"strings"

	fastfound "dfir-orc/fastburnt/internal/fastfind"

	log "github.com/sirupsen/logrus"
)

func (c *criteriaEntry) Match(sha256 string, sha1 string, md5 string, fullname string) bool {
	log.Trace("Trying to match sha256:" + sha256 + " sha1:" + sha1 + " md5:" + md5 + " filename" + fullname)
	log.Trace("With sha256:" + c.SHA256 + " sha1:" + c.SHA1 + " md5:" + c.MD5 + " filere" + c.FileRE)
	if c.SHA256 != "" && strings.EqualFold(c.SHA256, sha256) {
		return true
	}
	if c.SHA1 != "" && strings.EqualFold(c.SHA1, sha1) {
		return true
	}
	if c.MD5 != "" && strings.EqualFold(c.MD5, md5) {
		return true
	}
	if c.Regexp != nil {
		matched := c.Regexp.MatchString(fullname)
		if matched {
			return true
		}
	}
	return false
}

func isListed(criterias *[]criteriaEntry, m *fastfound.FastFindMatch) (bool, string) {
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
func (f *Filter) IsWhitelisted(m *fastfound.FastFindMatch) (bool, string) {
	log.Trace("Matching in whitelist: " + m.Fullname + " " + m.SHA256)
	return isListed(&f.whitelistCriteria, m)
}

// Match - method returning whether a FastFind Result looks like a specialy interresting file (for IOC loaded in filter_data.go)
func (f *Filter) IsBlacklisted(m *fastfound.FastFindMatch) (bool, string) {
	log.Trace("Matching in blacklist: " + m.Fullname + " " + m.SHA256)
	return isListed(&f.blacklistCriteria, m)
}
