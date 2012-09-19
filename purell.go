/*
Package purell offers URL normalization as described on the wikipedia page:
http://en.wikipedia.org/wiki/URL_normalization
*/
package purell

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

// A set of normalization flags determines how a URL will
// be normalized.
type NormalizationFlags int

const (
	// Safe normalizations
	FlagLowercaseScheme NormalizationFlags = 1 << iota
	FlagLowercaseHost
	FlagUppercaseEscapes
	FlagDecodeUnnecessaryEscapes
	FlagRemoveDefaultPort
	FlagRemoveEmptyQuerySeparator

	// Usually safe normalizations
	FlagRemoveTrailingSlash // Should choose one or the other (in add-remove slash)
	FlagAddTrailingSlash
	FlagRemoveDotSegments

	// Unsafe normalizations
	FlagRemoveDirectoryIndex
	FlagRemoveFragment
	FlagForceHttp
	FlagRemoveDuplicateSlashes
	FlagRemoveWWW // Should choose one or the other (in add-remove www)
	FlagAddWWW
	FlagSortQuery

	FlagsSafe = FlagLowercaseHost | FlagLowercaseScheme | FlagUppercaseEscapes | FlagDecodeUnnecessaryEscapes | FlagRemoveDefaultPort | FlagRemoveEmptyQuerySeparator

	FlagsUsuallySafe = FlagsSafe | FlagRemoveTrailingSlash | FlagRemoveDotSegments

	FlagsUnsafe = FlagsUsuallySafe | FlagRemoveDirectoryIndex | FlagRemoveFragment | FlagForceHttp | FlagRemoveDuplicateSlashes | FlagRemoveWWW | FlagSortQuery
)

var rxPort = regexp.MustCompile(`(:\d+)/?$`)
var rxDirIndex = regexp.MustCompile(`(^|/)((?:default|index)\.\w{1,4})$`)
var rxDupSlashes = regexp.MustCompile(`/{2,}`)

// MustNormalizeURLString returns the normalized URL as a string. It panics if
// the URL cannot be parsed.
func MustNormalizeURLString(u string, f NormalizationFlags) string {
	s, err := NormalizeURLString(u, f)
	if err != nil {
		panic(err)
	}
	return s
}

// NormalizeURLString returns the returns the normalized URL as
// as a string.
func NormalizeURLString(u string, f NormalizationFlags) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	NormalizeURL(parsed, f)
	return parsed.String(), nil
}

var transforms = []struct {
	flag      NormalizationFlags
	normalize func(*url.URL)
}{
	{FlagLowercaseScheme, lowercaseScheme},
	{FlagLowercaseHost, lowercaseHost},
	{FlagRemoveDefaultPort, removeDefaultPort},
	{FlagRemoveTrailingSlash, removeTrailingSlash},
	{FlagRemoveDirectoryIndex, removeDirectoryIndex}, // Must be before add trailing slash
	{FlagAddTrailingSlash, addTrailingSlash},
	{FlagRemoveDotSegments, removeDotSegments},
	{FlagRemoveFragment, removeFragment},
	{FlagForceHttp, forceHttp},
	{FlagRemoveDuplicateSlashes, removeDuplicateSlashes},
	{FlagRemoveWWW, removeWWW},
	{FlagAddWWW, addWWW},
	{FlagSortQuery, sortQuery},
}

// NormalizeURL normalizes the given URL according to the
// given flags.
func NormalizeURL(u *url.URL, f NormalizationFlags) {
	for _, t := range transforms {
		if f&t.flag == t.flag {
			t.normalize(u)
		}
	}
}

func lowercaseScheme(u *url.URL) {
	u.Scheme = strings.ToLower(u.Scheme)
}

func lowercaseHost(u *url.URL) {
	u.Host = strings.ToLower(u.Host)
}

func removeDefaultPort(u *url.URL) {
	if len(u.Host) > 0 {
		u.Host = rxPort.ReplaceAllStringFunc(u.Host, func(val string) string {
			if val == ":80" {
				return ""
			}
			return val
		})
	}
}

func removeTrailingSlash(u *url.URL) {
	if l := len(u.Path); l > 0 && strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path[:l-1]
	} else if l = len(u.Host); l > 0 && strings.HasSuffix(u.Host, "/") {
		u.Host = u.Host[:l-1]
	}
}

func addTrailingSlash(u *url.URL) {
	if l := len(u.Path); l > 0 && !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	} else if l = len(u.Host); l > 0 && !strings.HasSuffix(u.Host, "/") {
		u.Host += "/"
	}
}

func removeDotSegments(u *url.URL) {
	var dotFree []string

	if len(u.Path) > 0 {
		sections := strings.Split(u.Path, "/")
		for _, s := range sections {
			if s == ".." {
				if len(dotFree) > 0 {
					dotFree = dotFree[:len(dotFree)-1]
				}
			} else if s != "." {
				dotFree = append(dotFree, s)
			}
		}
		// Special case if host does not end with / and new path does not begin with /
		u.Path = strings.Join(dotFree, "/")
		if !strings.HasSuffix(u.Host, "/") && !strings.HasPrefix(u.Path, "/") {
			u.Path = "/" + u.Path
		}
	}
}

func removeDirectoryIndex(u *url.URL) {
	if len(u.Path) > 0 {
		u.Path = rxDirIndex.ReplaceAllString(u.Path, "$1")
	}
}

func removeFragment(u *url.URL) {
	u.Fragment = ""
}

func forceHttp(u *url.URL) {
	if strings.ToLower(u.Scheme) == "https" {
		u.Scheme = "http"
	}
}

func removeDuplicateSlashes(u *url.URL) {
	if len(u.Path) > 0 {
		u.Path = rxDupSlashes.ReplaceAllString(u.Path, "/")
	}
}

func removeWWW(u *url.URL) {
	if len(u.Host) > 0 && strings.HasPrefix(strings.ToLower(u.Host), "www.") {
		u.Host = u.Host[4:]
	}
}

func addWWW(u *url.URL) {
	if len(u.Host) > 0 && !strings.HasPrefix(strings.ToLower(u.Host), "www.") {
		u.Host = "www." + u.Host
	}
}

func sortQuery(u *url.URL) {
	q := u.Query()
	if len(q) == 0 {
		return
	}
	arKeys := make([]string, len(q))
	i := 0
	for k, _ := range q {
		arKeys[i] = k
		i++
	}
	sort.Strings(arKeys)
	buf := new(bytes.Buffer)
	for _, k := range arKeys {
		sort.Strings(q[k])
		for _, v := range q[k] {
			if buf.Len() > 0 {
				buf.WriteRune('&')
			}
			buf.WriteString(fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
		}
	}

	// Rebuild the raw query string
	u.RawQuery = buf.String()
}