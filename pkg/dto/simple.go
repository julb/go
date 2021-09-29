package dto

// ValueDTO is a simple object which contains a value
type ValueDTO struct {
	Value interface{} `yaml:"value" json:"value"`
}

// Large content DTO enables to hold large static content.
type LargeContentDTO struct {
	MimeType string `yaml:"mimeType" json:"mimeType"`
	Content  string `yaml:"content" json:"content"`
}

// Large content DTO enables to hold large static content.
type LargeContentWithSubjectDTO struct {
	MimeType string `yaml:"mimeType" json:"mimeType"`
	Subject  string `yaml:"subject" json:"subject"`
	Content  string `yaml:"content" json:"content"`
}
