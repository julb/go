package oas

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

type IndexOpts struct {
	Directory  string
	Extensions []string

	Url string
}

func NewIndexOpts() *IndexOpts {
	return &IndexOpts{
		Extensions: []string{".json", ".yaml", ".yml"},
		Url:        "",
	}
}

type V1_RepositoryIndex struct {
	ApiVersion int                                               `yaml:"apiVersion" json:"apiVersion"`
	Entries    map[string][]V1_RepositoryIndexSpecificationEntry `yaml:"entries" json:"entries"`
}

func NewV1_RepositoryIndex() *V1_RepositoryIndex {
	return &V1_RepositoryIndex{
		ApiVersion: 1,
		Entries:    make(map[string][]V1_RepositoryIndexSpecificationEntry),
	}
}

func (r V1_RepositoryIndex) AddSpecificationEntry(e *V1_RepositoryIndexSpecificationEntry) {
	// Add entry
	r.Entries[e.Name] = append(r.Entries[e.Name], *e)
}

func (r V1_RepositoryIndex) SortByVersionDesc() {
	for _, entries := range r.Entries {
		// Sort array by reverse version (latest on top)
		sort.SliceStable(entries, func(i, j int) bool {
			a, _ := semver.NewVersion(entries[i].Version)
			b, _ := semver.NewVersion(entries[j].Version)
			return a.GreaterThan(b)
		})
	}
}

type V1_RepositoryIndexSpecificationEntry struct {
	ApiVersion       int      `yaml:"apiVersion" json:"apiVersion"`
	BusinessCategory string   `yaml:"businessCategory" json:"businessCategory"`
	Deprecated       bool     `yaml:"deprecated" json:"deprecated"`
	Description      string   `yaml:"description" json:"description"`
	DisplayName      string   `yaml:"displayName" json:"displayName"`
	Keywords         []string `yaml:"keywords" json:"keywords"`
	LongDescription  string   `yaml:"longDescription" json:"longDescription"`
	Name             string   `yaml:"name" json:"name"`
	Starred          bool     `yaml:"starred" json:"starred"`
	Tags             []string `yaml:"tags" json:"tags"`
	TermsOfService   string   `yaml:"termsOfService" json:"termsOfService"`
	Url              string   `yaml:"url" json:"url"`
	Version          string   `yaml:"version" json:"version"`

	Contact struct {
		Name  string `yaml:"name" json:"name"`
		Email string `yaml:"email" json:"email"`
		Url   string `yaml:"url" json:"url"`
	} `yaml:"contact" json:"contact"`

	Image struct {
		Icon      string `yaml:"icon" json:"icon"`
		Logo      string `yaml:"logo" json:"logo"`
		Thumbnail string `yaml:"thumbnail" json:"thumbnail"`
	} `yaml:"image" json:"image"`

	License struct {
		Name string `yaml:"name" json:"name"`
		Url  string `yaml:"url" json:"url"`
	} `yaml:"license" json:"license"`

	Vcs struct {
		GitUrl      string `yaml:"gitUrl" json:"gitUrl"`
		GitRevision string `yaml:"gitRevision" json:"gitRevision"`
	} `yaml:"vcs" json:"vcs"`
}

func NewV1_RepositoryIndexSpecificationEntry() *V1_RepositoryIndexSpecificationEntry {
	return &V1_RepositoryIndexSpecificationEntry{
		ApiVersion: 1,
	}
}

func Index(opts *IndexOpts) error {
	log.Infof("Indexing specifications located in directory: %s.", opts.Directory)

	log.Debugf("Index directory: %s", opts.Directory)
	log.Debugf("Public URL: %s", opts.Url)
	log.Debugf("File extensions: %s", opts.Extensions)

	// Check if directory exists
	log.Debugf("Verify that input directory exists.")
	_, err := os.Stat(opts.Directory)
	if err != nil {
		return err
	}

	// Scan files candidates.
	candidateFiles, err := scanFiles(opts)
	if err != nil {
		return err
	}
	log.Debugf("%d candidate files found within the directory.", len(candidateFiles))

	// Build repository index
	log.Debugf("Read each file in directory to build repository data.")
	repositoryData, err := buildRepositoryData(opts, candidateFiles)
	if err != nil {
		return err
	}

	// Marshall repository content
	log.Debugf("Marshall repository data into index and yaml files.")
	err = marshallIndex(opts.Directory, repositoryData)
	if err != nil {
		return err
	}

	// Log indexation is complete.
	log.Infof("Indexation complete.")

	return nil
}

func scanFiles(o *IndexOpts) ([]string, error) {
	var files []string

	filepath.Walk(o.Directory, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			// Skip root index.yaml
			if path == filepath.Join(o.Directory, "index.yaml") {
				return nil
			}

			// Skip root index.json
			if path == filepath.Join(o.Directory, "index.json") {
				return nil
			}

			// Analyze subfiles.
			log.Tracef("> Checking file <%s>.", path)
			extension := strings.ToLower(filepath.Ext(path))
			for _, v := range o.Extensions {
				if v == extension {
					log.Debugf("> Include file <%s>.", path)
					files = append(files, path)
				}
			}
		}
		return nil
	})

	return files, nil
}

func buildRepositoryData(o *IndexOpts, candidateFiles []string) (*V1_RepositoryIndex, error) {
	// Initialize repo index.
	repositoryIndex := NewV1_RepositoryIndex()

	// Scan each file and accumulate content in the structure.
	for _, candidateFile := range candidateFiles {
		log.Debugf("Processing file <%s>.", candidateFile)

		// Parse specification
		oas3Source, err := ParseFile(candidateFile)
		if err != nil {
			return nil, err
		}

		// Convert to entry
		specificationEntry, err := buildSpecificationEntry(o, oas3Source)
		if err != nil {
			return nil, err
		}

		// Add spec entry to repo index.
		if specificationEntry != nil {
			repositoryIndex.AddSpecificationEntry(specificationEntry)
		}
	}

	// Force sort
	repositoryIndex.SortByVersionDesc()

	// Return result.
	return repositoryIndex, nil
}

func buildSpecificationEntry(o *IndexOpts, oas3Source *OAS3Source) (*V1_RepositoryIndexSpecificationEntry, error) {
	// Compute specification URL
	specificationUrl := strings.TrimPrefix(oas3Source.path, o.Directory)
	if o.Url != "" {
		specificationUrl = strings.Join([]string{o.Url, specificationUrl}, "/")
	}

	// Build entry from file.
	specificationEntry := NewV1_RepositoryIndexSpecificationEntry()
	specificationEntry.BusinessCategory = oas3Source.specification.Info.ExtraInfo.BusinessCategory
	specificationEntry.Contact.Email = oas3Source.specification.Info.Contact.Email
	specificationEntry.Contact.Name = oas3Source.specification.Info.Contact.Name
	specificationEntry.Contact.Url = oas3Source.specification.Info.Contact.Url
	specificationEntry.Deprecated = oas3Source.specification.Info.ExtraInfo.Deprecated
	specificationEntry.Description = oas3Source.specification.Info.Description
	specificationEntry.DisplayName = oas3Source.specification.Info.ExtraInfo.DisplayName
	specificationEntry.Image.Icon = oas3Source.specification.Info.ExtraInfo.IconUrl
	specificationEntry.Image.Logo = oas3Source.specification.Info.ExtraInfo.LogoUrl
	specificationEntry.Image.Thumbnail = oas3Source.specification.Info.ExtraInfo.ThumbnailUrl
	specificationEntry.Keywords = oas3Source.specification.Info.ExtraInfo.Keywords
	specificationEntry.License.Name = oas3Source.specification.Info.License.Name
	specificationEntry.License.Url = oas3Source.specification.Info.License.Url
	specificationEntry.LongDescription = oas3Source.specification.Info.ExtraInfo.LongDescription
	specificationEntry.Name = oas3Source.specification.Info.Title
	specificationEntry.Starred = oas3Source.specification.Info.ExtraInfo.Starred
	specificationEntry.Tags = oas3Source.specification.Info.ExtraInfo.Tags
	specificationEntry.Url = specificationUrl
	specificationEntry.Vcs.GitRevision = oas3Source.specification.Info.ExtraInfo.VcsGitRevision
	specificationEntry.Vcs.GitUrl = oas3Source.specification.Info.ExtraInfo.VcsGitUrl
	specificationEntry.Version = oas3Source.specification.Info.Version
	return specificationEntry, nil
}

func marshallIndex(directory string, data *V1_RepositoryIndex) error {
	// Marshalling into JSON
	indexJsonFilePath := filepath.Join(directory, "index.json")

	log.Debugf("Marshalling repository index to %s.", indexJsonFilePath)
	jsonMarshalled, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	log.Tracef("index.json dump: \n%s\n", string(jsonMarshalled))
	err = ioutil.WriteFile(filepath.Join(directory, "index.json"), jsonMarshalled, 0644)
	if err != nil {
		return err
	}

	// Marshalling into YAML
	indexYamlFilePath := filepath.Join(directory, "index.yaml")

	log.Debugf("Marshalling repository index into %s.", indexYamlFilePath)
	yamlMarshalled, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}
	log.Tracef("index.yaml dump: \n%s\n", string(yamlMarshalled))
	err = ioutil.WriteFile(indexYamlFilePath, yamlMarshalled, 0644)
	if err != nil {
		return err
	}

	log.Debug("Marshalling complete.")

	return nil
}
