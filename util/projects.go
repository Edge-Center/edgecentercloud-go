package util

import (
	"context"
	"fmt"
	"log"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

// findProjectByName searches for a project with the specified name in the provided project slice.
// Returns the project ID if found, otherwise returns an error.
func findProjectByName(arr []edgecloud.Project, name string) (int, error) {
	for _, el := range arr {
		if el.Name == name {
			return el.ID, nil
		}
	}
	return 0, fmt.Errorf("project with name %s not found", name)
}

// GetProject returns a valid project ID.
// If the projectID is provided, it will be returned directly.
// If projectName is provided instead, the function will search for the project by name and return its ID.
// Returns an error if the project is not found or there is an issue with the client.
func GetProject(ctx context.Context, client *edgecloud.Client, projectName string) (int, error) {
	log.Println("[DEBUG] Try to get project ID")

	projectsList, _, err := client.Projects.List(ctx, &edgecloud.ProjectListOptions{})
	if err != nil {
		return 0, err
	}
	log.Printf("[DEBUG] Projects: %v", projectsList)
	projectID, err := findProjectByName(projectsList, projectName)
	if err != nil {
		return 0, err
	}
	log.Printf("[DEBUG] The attempt to get the project is successful: projectID=%d", projectID)

	return projectID, nil
}
