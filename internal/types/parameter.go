package types

// Parameters defines all available parameters
type Parameters struct {
	Query  []*Parameter
	Header []*Parameter
	Cookie []*Parameter
	Path   []*Parameter
}

// Parameter holds a parameter and example value
type Parameter struct {
	Name  string
	Value interface{}
}
