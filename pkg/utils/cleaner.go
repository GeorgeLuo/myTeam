package utils

import "strings"

func CleanResponse(raw string) (strData string) {
	strData = string(raw)
	if strings.HasPrefix(strData, "```json") && strings.HasSuffix(strData, "```") {
		// Remove the very first and last characters,
		// ensuring removal of the code fences even if the JSON is indented or formatted strangely
		strData = strData[7 : len(strData)-3] // Removes starting ```json and ending ```
		strData = strings.TrimSpace(strData)
	}
	return
}
