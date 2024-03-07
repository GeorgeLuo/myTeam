package workspace

// EmployeeMetadata holds the metadata for an employee.
type EmployeeMetadata struct {
	Name        string
	Description string
	Prompt      string
}

// Workspace represents a workspace containing personnel information.
type Workspace struct {
	Personnel map[string]EmployeeMetadata
}

// NewWorkspace creates a new Workspace instance.
func NewWorkspace(directory string) Workspace {
	return Workspace{
		Personnel: make(map[string]EmployeeMetadata),
	}
}

// AddPersonnel adds a new employee's metadata to the workspace.
func (w *Workspace) AddPersonnel(id string, name string, description string, prompt string) {
	// Create the EmployeeMetadata for the employee.
	employee := EmployeeMetadata{
		Name:        name,
		Description: description,
		Prompt:      prompt,
	}

	// Add or update the employee in the map using the provided id.
	w.Personnel[id] = employee
}
