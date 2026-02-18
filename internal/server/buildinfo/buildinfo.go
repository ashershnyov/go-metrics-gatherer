package buildinfo

import "strings"

const na = "N/A"

// BuildInfo represents servers's build info.
type BuildInfo struct {
	Version string
	Date    string
	Commit  string
}

// New returns new servers's build information object.
func New(version, date, commit string) BuildInfo {
	bi := BuildInfo{
		Version: version,
		Date:    date,
		Commit:  commit,
	}
	if bi.Version == "" {
		bi.Version = na
	}
	if bi.Date == "" {
		bi.Date = na
	}
	if bi.Commit == "" {
		bi.Commit = na
	}
	return bi
}

// String returns string representation of BuildInfo.
func (bi *BuildInfo) String() string {
	var sb strings.Builder
	sb.WriteString("Build version: ")
	sb.WriteString(bi.Version)
	sb.WriteString("\nBuild date: ")
	sb.WriteString(bi.Date)
	sb.WriteString("\nBuild commit: ")
	sb.WriteString(bi.Commit)
	return sb.String()
}
