package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PersonnelMetadata holds the metadata for an employee.
type PersonnelMetadata struct {
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Prompt        string            `json:"prompt"`
	ModelVendor   string            `json:"model_vendor"`
	ModelMetadata map[string]string `json:"model_metadata"`
}

// Workspace represents a workspace containing personnel information.
type Workspace struct {
	Personnel map[string]PersonnelMetadata `json:"personnel"`
	Directory string                       `json:"-"`
}

// NewWorkspace creates a new Workspace instance and attempts to load personnel data from the specified directory.
func NewWorkspace(directory string) Workspace {
	workspace := Workspace{
		Personnel: make(map[string]PersonnelMetadata),
		Directory: directory,
	}
	// Ensure the directory exists
	if err := os.MkdirAll(directory, 0755); err != nil {
		panic("Failed to create workspace directory: " + err.Error())
	}
	workspace.load()
	return workspace
}

// AddPersonnel adds a new employee's metadata to the workspace and updates the JSON file.
func (w *Workspace) AddPersonnel(id string, name string, description string, prompt string,
	vendor string, modelMetadata map[string]string) {
	// Create the PersonnelMetadata for the employee.
	employee := PersonnelMetadata{
		Name:          name,
		Description:   description,
		Prompt:        prompt,
		ModelVendor:   vendor,
		ModelMetadata: modelMetadata,
	}

	// Add or update the employee in the map using the provided id.
	w.Personnel[id] = employee
	w.save()
}

// SetModelMetaDataByID updates or adds a new metadata key-value pair for a specific personnel by ID.
func (w *Workspace) SetModelMetaDataByID(id string, key string, value string) error {
	// Check if the personnel with the given ID exists
	personnel, exists := w.Personnel[id]
	if !exists {
		return fmt.Errorf("personnel with ID %s does not exist", id)
	}
	// Initialize the map if it's nil
	if personnel.ModelMetadata == nil {
		personnel.ModelMetadata = make(map[string]string)
	}
	// Set the key-value pair in model_metadata
	personnel.ModelMetadata[key] = value
	// Update the map entry
	w.Personnel[id] = personnel
	// Save changes
	w.save()
	return nil
}

// File returns the workspace as a file
func (w *Workspace) File() string {
	return filepath.Join(w.Directory, "workspace.json")
}

// save writes the current state of Workspace to a JSON file in the specified Directory.
func (w *Workspace) save() {
	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		panic("Failed to marshal workspace: " + err.Error())
	}

	filePath := filepath.Join(w.Directory, "workspace.json")
	// Ensure the directory exists before trying to save the file.
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		panic("Failed to create directory for workspace file: " + err.Error())
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		panic("Failed to write workspace file: " + err.Error())
	}
}

// load attempts to load personnel data from a JSON file in the specified Directory.
func (w *Workspace) load() {
	filePath := filepath.Join(w.Directory, "workspace.json")
	_, err := os.Stat(filePath)
	// If the file does not exist, create it with an empty workspace.
	if os.IsNotExist(err) {
		w.save() // This will create an empty file.
		return
	}
	if err != nil {
		panic("Failed to stat workspace file: " + err.Error())
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic("Failed to read workspace file: " + err.Error())
	}

	if err := json.Unmarshal(data, w); err != nil {
		panic("Failed to unmarshal workspace data: " + err.Error())
	}
}
