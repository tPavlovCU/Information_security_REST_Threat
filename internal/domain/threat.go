package domain

type Threat struct {
	ID          int64    `json:"id"`
	MD5         string   `json:"md5"`
	SHA256      *string  `json:"sha256"`
	Name        string   `json:"name"`
	Family      *string  `json:"family"`
	Type        string   `json:"type"`
	Severity    string   `json:"severity"`
	Confidence  *int16   `json:"confidence"`
	Description *string  `json:"description"`
	Tags        []string `json:"tags"`
	Source      *string  `json:"source"`
}

const (
	ThreatTypeTrojan     = "trojan"
	ThreatTypeRansomware = "ransomware"
	ThreatTypeVirus      = "virus"
	ThreatTypeSpyware    = "spyware"
)

const (
	ThreatSeverityLow      = "low"
	ThreatSeverityMedium   = "medium"
	ThreatSeverityHigh     = "high"
	ThreatSeverityCritical = "critical"
)
