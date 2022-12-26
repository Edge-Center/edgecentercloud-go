package projects

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/project/v1/types"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

func List(c *edgecloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, listURL(c), func(r pagination.PageResult) pagination.Page {
		return ProjectPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific project based on its unique ID.
func Get(c *edgecloud.ServiceClient, id int) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil) // nolint
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToProjectCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a project.
type CreateOpts struct {
	ClientID    int                `json:"client_id,omitempty"`
	State       types.ProjectState `json:"state,omitempty"`
	Name        string             `json:"name" required:"true" validate:"required"`
	Description string             `json:"description,omitempty"`
}

// ToProjectCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToProjectCreateMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

// Validate
func (opts CreateOpts) Validate() error {
	return edgecloud.Validate.Struct(opts)
}

// Create accepts a CreateOpts struct and creates a new project using the values provided.
func Create(c *edgecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToProjectCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil) // nolint
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToProjectUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a project.
type UpdateOpts struct {
	Name        string `json:"name,omitempty" validate:"required_without_all=Description"`
	Description string `json:"description,omitempty" validate:"required_without_all=Name"`
}

// ToProjectUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToProjectUpdateMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

// Validate
func (opts UpdateOpts) Validate() error {
	return edgecloud.Validate.Struct(opts)
}

// Update accepts a UpdateOpts struct and updates an existing project using the values provided.
func Update(c *edgecloud.ServiceClient, id int, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToProjectUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, id), b, &r.Body, &edgecloud.RequestOpts{ // nolint
		OkCodes: []int{200, 201},
	})
	return
}

// ListAll is a convenience function that returns all projects.
func ListAll(client *edgecloud.ServiceClient) ([]Project, error) {
	pages, err := List(client).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractProjects(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// Delete a project
func Delete(client *edgecloud.ServiceClient, id int) (r tasks.Result) {
	url := deleteURL(client, id)
	_, r.Err = client.DeleteWithResponse(url, &r.Body, nil) // nolint
	return
}
