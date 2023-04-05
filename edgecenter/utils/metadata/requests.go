package metadata

import (
	"net/http"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

func ResourceMetadataList(client *edgecloud.ServiceClient, id string) pagination.Pager {
	url := ResourceMetadataURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ResourceMetadataPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func ResourceMetadataListAll(client *edgecloud.ServiceClient, id string) ([]Metadata, error) {
	pages, err := ResourceMetadataList(client, id).AllPages()
	if err != nil {
		return nil, err
	}
	all, err := ExtractMetadata(pages)
	if err != nil {
		return nil, err
	}
	return all, nil
}

// ResourceMetadataCreateOrUpdate creates or update a metadata for a resource.
func ResourceMetadataCreateOrUpdate(client *edgecloud.ServiceClient, id string, opts map[string]string) (r ActionResultMetadata) {
	_, r.Err = client.Post(ResourceMetadataURL(client, id), opts, nil, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// ResourceMetadataReplace replace a metadata for a resource.
func ResourceMetadataReplace(client *edgecloud.ServiceClient, id string, opts map[string]string) (r ActionResultMetadata) {
	_, r.Err = client.Put(ResourceMetadataURL(client, id), opts, nil, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// ResourceMetadataDelete deletes defined metadata key for a resource.
func ResourceMetadataDelete(client *edgecloud.ServiceClient, id string, key string) (r ActionResultMetadata) {
	_, r.Err = client.Delete(ResourceMetadataItemURL(client, id, key), &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// ResourceMetadataGet gets defined metadata key for a resource.
func ResourceMetadataGet(client *edgecloud.ServiceClient, id string, key string) (r ResourceMetadataResult) {
	url := ResourceMetadataItemURL(client, id, key)

	_, r.Err = client.Get(url, &r.Body, nil)
	return
}
