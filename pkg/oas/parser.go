package oas

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

type OAS3Source struct {
	path          string
	specification *OAS3Specification
}

type OAS3Specification struct {
	Info struct {
		Title          string `yaml:"title" json:"title"`
		Version        string `yaml:"version" json:"version"`
		Description    string `yaml:"description" json:"description"`
		TermsOfService string `yaml:"termsOfService" json:"termsOfService"`

		Contact struct {
			Name  string `yaml:"name" json:"name"`
			Email string `yaml:"email" json:"email"`
			Url   string `yaml:"url" json:"url"`
		} `yaml:"contact" json:"contact"`

		License struct {
			Name string `yaml:"name" json:"name"`
			Url  string `yaml:"url" json:"url"`
		} `yaml:"license" json:"license"`

		ExtraInfo struct {
			BusinessCategory string   `yaml:"businessCategory" json:"businessCategory"`
			Deprecated       bool     `yaml:"deprecated" json:"deprecated"`
			DisplayName      string   `yaml:"displayName" json:"displayName"`
			IconUrl          string   `yaml:"iconUrl" json:"iconUrl"`
			Keywords         []string `yaml:"keywords" json:"keywords"`
			LogoUrl          string   `yaml:"logoUrl" json:"logoUrl"`
			LongDescription  string   `yaml:"longDescription" json:"longDescription"`
			Starred          bool     `yaml:"starred" json:"starred"`
			Tags             []string `yaml:"tags" json:"tags"`
			ThumbnailUrl     string   `yaml:"thumbnailUrl" json:"thumbnailUrl"`
			VcsGitRevision   string   `yaml:"vcsGitRevision" json:"vcsGitRevision"`
			VcsGitUrl        string   `yaml:"vcsGitUrl" json:"vcsGitUrl"`
		} `yaml:"x-extra-info" json:"x-extra-info"`
	} `yaml:"info" json:"info"`
}

func ParseFile(path string) (*OAS3Source, error) {
	// Initialize repo index.
	log.Debugf("Parse OAS specifiction <%s>.", path)

	// Read file
	candidateFileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshall
	var candidateFileSpecification OAS3Specification
	if filepath.Ext(path) == ".json" {
		// Use JSON loader
		err = json.Unmarshal(candidateFileBytes, &candidateFileSpecification)
		if err != nil {
			return nil, err
		}
	} else if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
		// Use YAML loader.
		err = yaml.Unmarshal(candidateFileBytes, &candidateFileSpecification)
		if err != nil {
			return nil, err
		}
	}

	return &OAS3Source{
		path:          path,
		specification: &candidateFileSpecification,
	}, nil
}
