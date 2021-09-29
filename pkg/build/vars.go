package build

type BuildInfo struct {
	Group            string `yaml:"group" json:"group"`
	Artifact         string `yaml:"artifact" json:"artifact"`
	Name             string `yaml:"name" json:"name"`
	Version          string `yaml:"version" json:"version"`
	BuildVersion     string `yaml:"buildVersion" json:"buildVersion"`
	Time             string `yaml:"time" json:"time"`
	GitRevision      string `yaml:"gitRevision" json:"gitRevision"`
	GitShortRevision string `yaml:"gitShortRevision" json:"gitShortRevision"`
	Arch             string `yaml:"arch" json:"arch"`
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
