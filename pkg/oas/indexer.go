package oas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

type IndexOptions struct {
	Directory  string
	Extensions []string

	Url string
}

func NewIndexOptions() *IndexOptions {
	return &IndexOptions{
		Extensions: []string{".json", ".yaml", ".yml"},
		Url:        "",
	}
}

type V1_RepositoryIndex struct {
	ApiVersion int                                `yaml:"apiVersion" json:"apiVersion"`
	Entries    map[string][]V1_SpecificationEntry `yaml:"entries" json:"entries"`
}

func NewV1_RepositoryIndex() *V1_RepositoryIndex {
	return &V1_RepositoryIndex{
		ApiVersion: 1,
		Entries:    make(map[string][]V1_SpecificationEntry),
	}
}

func (r V1_RepositoryIndex) AddSpecificationEntry(e *V1_SpecificationEntry) {
	r.Entries[e.Name] = append(r.Entries[e.Name], *e)
}

type V1_SpecificationEntry struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}

func NewV1_SpecificationEntry() *V1_SpecificationEntry {
	return &V1_SpecificationEntry{}
}

func Index(o *IndexOptions) error {
	log.Infof("Indexing specifications located in directory: %s.", o.Directory)

	log.Debugf("Index directory: %s", o.Directory)
	log.Debugf("Public URL: %s", o.Url)
	log.Debugf("File extensions: %s", o.Extensions)

	// Check if directory exists
	log.Debugf("Verify that input directory exists.")
	_, err := os.Stat(o.Directory)
	if err != nil {
		return err
	}

	// Scan files candidates.
	candidateFiles, err := scanFiles(o.Directory, o.Extensions)
	if err != nil {
		return err
	}
	log.Infof("%d candidate files found within the directory.", len(candidateFiles))

	// Build repository index
	log.Debugf("Read each file in directory to build repository data.")
	repositoryData, err := buildRepositoryData(candidateFiles)
	if err != nil {
		return err
	}

	// Marshall repository content
	log.Debugf("Marshall repository data into index and yaml files.")
	err = marshallIndex(o.Directory, repositoryData)
	if err != nil {
		return err
	}

	// Log indexation is complete.
	log.Infof("Indexation complete.")

	return nil
}

func scanFiles(directory string, extensions []string) ([]string, error) {
	var files []string

	filepath.Walk(directory, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			// Skip root index.yaml
			if path == filepath.Join(directory, "index.yaml") {
				return nil
			}

			// Skip root index.json
			if path == filepath.Join(directory, "index.json") {
				return nil
			}

			// Analyze subfiles.
			log.Debugf("> Checking file <%s>.", path)
			extension := strings.ToLower(filepath.Ext(path))
			for _, v := range extensions {
				if v == extension {
					log.Infof("> Retaining candidate file <%s>.", path)
					files = append(files, path)
				}
			}
		}
		return nil
	})

	return files, nil
}

func buildRepositoryData(candidateFiles []string) (*V1_RepositoryIndex, error) {
	// Initialize repo index.
	repositoryIndex := NewV1_RepositoryIndex()

	// Scan each file and accumulate content in the structure.
	for _, candidateFile := range candidateFiles {
		fmt.Println(candidateFile)

		// Build entry from file.
		specificationEntry := NewV1_SpecificationEntry()
		specificationEntry.Name = "some-api"
		specificationEntry.Version = "1.0.0"

		// Add spec entry to repo index.
		repositoryIndex.AddSpecificationEntry(specificationEntry)
	}

	// Return result.
	return repositoryIndex, nil
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
