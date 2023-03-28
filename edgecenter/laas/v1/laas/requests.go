package laas

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

// UpdateOptsBuilder allows extensions to add additional parameters to the request.
type UpdateOptsBuilder interface {
	ToStatusUpdateMap() (map[string]interface{}, error)
}

// UpdateStatusOpts represents options used to update a laas status.
type UpdateStatusOpts struct {
	IsInitialized bool `json:"is_initialized"`
}

// ToStatusUpdateMap builds a request body from UpdateStatusOpts.
func (opts UpdateStatusOpts) ToStatusUpdateMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

// GetStatus retrieves laas status.
func GetStatus(c *edgecloud.ServiceClient) (r StatusResult) {
	_, r.Err = c.Get(statusURL(c), &r.Body, nil)
	return
}

// UpdateStatus update LaaS status.
func UpdateStatus(c *edgecloud.ServiceClient, opts UpdateOptsBuilder) (r StatusResult) {
	url := statusURL(c)
	b, err := opts.ToStatusUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(url, b, &r.Body, &edgecloud.RequestOpts{OkCodes: []int{200}})
	return
}

// RegenerateUser regenerate LaaS credentials.
func RegenerateUser(c *edgecloud.ServiceClient) (r UserResult) {
	_, r.Err = c.Post(usersURL(c), nil, &r.Body, nil)
	return
}

// ListTopic list LaaS Kafka topics within client namespace.
func ListTopic(c *edgecloud.ServiceClient) pagination.Pager {
	url := topicsURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return TopicPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListTopicAll list LaaS Kafka topics within client namespace.
func ListTopicAll(c *edgecloud.ServiceClient) ([]Topic, error) {
	page, err := ListTopic(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractTopics(page)
}

// CreateTopicOptsBuilder allows extensions to add additional parameters to the request.
type CreateTopicOptsBuilder interface {
	ToTopicCreateMap() (map[string]interface{}, error)
}

// CreateTopicOpts represents options used to create a topic.
type CreateTopicOpts struct {
	Name string `json:"name"`
}

// ToTopicCreateMap builds a request body from CreateTopicOpts.
func (opts CreateTopicOpts) ToTopicCreateMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

// CreateTopic create LaaS topic.
func CreateTopic(c *edgecloud.ServiceClient, opts CreateTopicOptsBuilder) (r TopicResult) {
	url := topicsURL(c)
	b, err := opts.ToTopicCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(url, b, &r.Body, &edgecloud.RequestOpts{OkCodes: []int{201}})
	return
}

// DeleteTopic delete LaaS Kafka topic within client namespace.
func DeleteTopic(c *edgecloud.ServiceClient, name string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteTopicURL(c, name), nil)
	return
}

// ListKafkaHosts retrieves LaaS kafka hosts.
func ListKafkaHosts(c *edgecloud.ServiceClient) (r HostsResult) {
	_, r.Err = c.Get(kafkaURL(c), &r.Body, nil)
	return
}

// ListOpenSearchHosts retrieves LaaS opensearch hosts.
func ListOpenSearchHosts(c *edgecloud.ServiceClient) (r HostsResult) {
	_, r.Err = c.Get(openSearchURL(c), &r.Body, nil)
	return
}
