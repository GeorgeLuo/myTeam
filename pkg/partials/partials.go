package partials

import (
	"fmt"
	"os"
)

func IdAssignment(num int) string {
	const IdAssignmentFmt = "You are employee %d."
	return fmt.Sprintf(IdAssignmentFmt, num)
}

func HiringCapabilities(num int) string {
	const HiringCapabilitiesFmt = "You can make up to %d hires."
	return fmt.Sprintf(HiringCapabilitiesFmt, num)
}

func LoadFromFile(filename string) string {
	// Read the contents of the file specified by filename
	data, err := os.ReadFile(filename)
	if err != nil {
		// If there's an error reading the file, return an empty string
		return ""
	}
	// Return the contents of the file as a string
	return string(data)
}
