package delegations

import (
	"encoding/json"
	"fmt"
	"myTeam/pkg/utils"
)

// HiringData represents the top-level structure of the hiring information.
type HiringData struct {
	HiringAgentID string `json:"hiring_agent_id"`
	Roles         []Role `json:"roles"`
}

// Role represents information about a specific role recommended for hiring.
type Role struct {
	Title               string           `json:"title"`
	TmpID               string           `json:"tmp_id"`
	Pseudonym           string           `json:"pseudonym"`
	TopLevelRequirement string           `json:"top_level_requirement"`
	ReportsTo           []string         `json:"reports_to"`
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
	strData = utils.CleanResponse(strData)
	fmt.Println("Cleaned JSON string:", strData)
	err := json.Unmarshal([]byte(strData), &hiringData)
	if err != nil {
		return HiringData{}, err
	}
	return hiringData, nil
}
