package config

import (
	"regexp"

	cconfig "github.com/coreos/coreos-cloudinit/config"
)

var emptyString = regexp.MustCompile(`(?m)^[^:]+: ""\n`)
var emptyArray = regexp.MustCompile(`(?m)^[^:]+: \[\]\n`)
var emptyKey1 = regexp.MustCompile(`(?m)^[^:]+:\n([^ \t\n\-])`)
var emptyKey2 = regexp.MustCompile(`(?m)^[ \t]{2}[^:]+:\n([ \t]{0,2}[^ \t\n\-])`)
var emptyKey3 = regexp.MustCompile(`(?m)^[ \t]{4}[^:]+:\n([ \t]{0,4}[^ \t\n\-])`)
var emptyKey4 = regexp.MustCompile(`(?m)^[ \t]{6}[^:]+:\n([ \t]{0,6}[^ \t\n\-])`)
var emptyKey5 = regexp.MustCompile(`(?m)^[ \t]{7}[^:]+:\n([ \t]{0,7}[^ \t\n\-])`)

func ProcessedCloudConfig(cc cconfig.CloudConfig) string {
	out := cc.String()
	out = emptyString.ReplaceAllString(out, "")
	out = emptyArray.ReplaceAllString(out, "")
	out = emptyKey1.ReplaceAllString(out, "$1")
	out = emptyKey2.ReplaceAllString(out, "$1")
	out = emptyKey3.ReplaceAllString(out, "$1")
	out = emptyKey4.ReplaceAllString(out, "$1")
	out = emptyKey5.ReplaceAllString(out, "$1")
	out = "#cloud-config\n" + out

	return out
}
