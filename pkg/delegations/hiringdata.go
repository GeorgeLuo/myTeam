package delegations

import "encoding/json"

// HiringData represents the top-level structure of the hiring information.
type HiringData struct {
	HiringEmployeeID string `json:"hiring_employee_id"`
	Roles            []Role `json:"roles"`
}

// Role represents information about a specific role recommended for hiring.
type Role struct {
	Title               string           `json:"title"`
	TmpID               string           `json:"tmp_id"`
	TopLevelRequirement string           `json:"top_level_requirement"`
	ReportsTo           []string         `json:"reports_to"`
	Responsibilities    []Responsibility `json:"responsibilities"`
}

// Responsibility represents a specific responsibility associated with a role.
type Responsibility struct {
	Description string `json:"description"`
}

// Example usage: Marshal an instance of HiringData to JSON.
func main() {
	// Example data for marshalling
	data := HiringData{
		HiringEmployeeID: "your employee id",
		Roles: []Role{
			{
				Title:               "title of first recommendation",
				TmpID:               "localized id for this recommendation document with prefix tmp_id_",
				TopLevelRequirement: "the axis upon which success is determined",
				ReportsTo:           []string{"employee id or tmp_id"},
				Responsibilities: []Responsibility{
					{
						Description: "description of responsibility",
					},
				},
			},
		},
	}

	// Marshal the data to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	// Print the JSON representation
	println(string(jsonData))
}
