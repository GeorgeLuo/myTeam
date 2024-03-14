package promptbuilder

// PromptOutline is used to display the components of a prompt
type PromptOutline struct {
	TopLevelRequirement string         `json:"top_level_requirement"`
	OrgMetadata         OrgMetaData    `json:"org_metadata"`
	Functions           []Function     `json:"functions"`
	Files               []FileMetadata `json:"files"`
}

type OrgMetaData struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	MaxDirectReports string   `json:"max_direct_reports"`
	ReportingTo      []string `json:"reporting_to"`
	DirectReports    []string `json:"direct_reports"`
}

type FileMetadata struct {
	Category string `json:"category"`
	Filename string `json:"filename"`
}

type Function struct {
	Description string `json:"desc"`
}
