package promptbuilder

// PromptOutline is used to display the components of a prompt
type PromptOutline struct {
	TopLevelRequirement string
	OrgMetadata         OrgMetaData
	Functions           []Function
	Files               []FileMetadata
}

type OrgMetaData struct {
	ID               string
	Name             string
	MaxDirectReports string
	ReportingTo      []string
	DirectReports    []string
}

type FileMetadata struct {
	Category string
	Filename string
}

type Function struct {
	Description string
}
