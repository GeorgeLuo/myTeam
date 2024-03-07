package delegations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// HiringData represents the top-level structure of the hiring information.
type HiringData struct {
	HiringEmployeeID int    `json:"hiring_employee_id"`
	Roles            []Role `json:"roles"`
}

// Role represents information about a specific role recommended for hiring.
type Role struct {
	Title               string           `json:"title"`
	TmpID               string           `json:"tmp_id"`
	Pseudonym           string           `json:"pseudonym"`
	TopLevelRequirement string           `json:"top_level_requirement"`
	ReportsTo           []int            `json:"reports_to"`
	Responsibilities    []Responsibility `json:"responsibilities"`
}

// Responsibility represents a specific responsibility associated with a role.
type Responsibility struct {
	Description string `json:"description"`
}

// Improved UnmarshalHiringData function.
func UnmarshalHiringData(data []byte) (HiringData, error) {
	var hiringData HiringData
	strData := string(data)
	// Enhanced cleaning logic
	if strings.HasPrefix(strData, "```json") && strings.HasSuffix(strData, "```") {
		// Remove the very first and last characters,
		// ensuring removal of the code fences even if the JSON is indented or formatted strangely
		strData = strData[7 : len(strData)-3] // Removes starting ```json and ending ```
		strData = strings.TrimSpace(strData)
	}
	// Debugging line: This can help ensure the string is clean before unmarshalling.
	fmt.Println("Cleaned JSON string:", strData)
	// Attempt to unmarshal the cleaned-up JSON string into the struct.
	err := json.Unmarshal([]byte(strData), &hiringData)
	if err != nil {
		// Additional logging here could help with debugging.
		//fmt.Println("Unmarshal error:", err)
		return HiringData{}, err
	}
	return hiringData, nil
}
