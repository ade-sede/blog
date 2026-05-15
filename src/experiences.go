package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// ExperienceEntry represents a single work or education entry on the resume.
type ExperienceEntry struct {
	Title             string   `json:"title"`
	Company           string   `json:"company"`
	Begin             string   `json:"begin"`
	End               string   `json:"end"`
	Description       string   `json:"description"`
	BulletPointsIntro string   `json:"bulletPointsIntro"`
	BulletPoints      []string `json:"bulletPoints"`
}

// ExperiencesData holds all resume experience entries, separated by category.
type ExperiencesData struct {
	WorkExperiences   []ExperienceEntry `json:"workExperiences"`
	SchoolExperiences []ExperienceEntry `json:"schoolExperiences"`
}

// loadExperiencesFromJSON reads and unmarshals the experiences JSON file
// at filename into an ExperiencesData.
func loadExperiencesFromJSON(filename string) (*ExperiencesData, error) {
	jsonBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	var experiences ExperiencesData
	err = json.Unmarshal(jsonBytes, &experiences)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}
	return &experiences, nil
}
