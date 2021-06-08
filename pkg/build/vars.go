package build

type BuildInfo struct {
	Group            string
	Artifact         string
	Name             string
	Version          string
	BuildVersion     string
	Time             string
	GitRevision      string
	GitShortRevision string
	Arch             string
}

var Group string
var Artifact string
var Name string
var Version string
var BuildVersion string
var Time string
var GitRevision string
var GitShortRevision string
var Arch string

var Info BuildInfo

func init() {
	Info = BuildInfo{
		Group:            Group,
		Artifact:         Artifact,
		Name:             Name,
		Version:          Version,
		BuildVersion:     BuildVersion,
		Time:             Time,
		GitRevision:      GitRevision,
		GitShortRevision: GitShortRevision,
		Arch:             Arch,
	}
}
