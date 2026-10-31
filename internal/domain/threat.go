package domain

type Threat struct {
	ID          int64
	MD5         string
	SHA256      *string
	Name        string
	Family      *string
	Type        string
	Severity    string
	Confidence  *int16
	Description *string
	Tags        []string
	Source      *string
}
