package types

// CodeSample represents a single code example
type CodeSample struct {
	Lang   Language `yaml:"lang" json:"lang"`     // language of the sample
	Source string   `yaml:"source" json:"source"` // the actual source code
	Label  string   `yaml:"label" json:"label"`   // displayed language name
}
