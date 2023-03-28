package edgecloud

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

var defaultSleepTimeout = 1

// WaitFor polls a predicate function, once per second, up to a timeout limit.
// This is useful to wait for a resource to transition to a certain state.
// To handle situations when the predicate might hang indefinitely, the
// predicate will be prematurely cancelled after the timeout.
// Resource packages will wrap this in a more convenient function that's
// specific to a certain resource, but it can also be useful on its own.
func WaitFor(timeout int, predicate func() (bool, error)) error {
	type WaitForResult struct {
		Success bool
		Error   error
	}

	start := time.Now().Unix()

	for {
		// If a timeout is set, and that's been exceeded, shut it down.
		if timeout >= 0 && time.Now().Unix()-start >= int64(timeout) {
			return fmt.Errorf("a timeout occurred")
		}

		time.Sleep(time.Duration(defaultSleepTimeout) * time.Second)

		var result WaitForResult
		ch := make(chan bool, 1)
		go func() {
			defer close(ch)
			satisfied, err := predicate()
			result.Success = satisfied
			result.Error = err
		}()

		select {
		case <-ch:
			if result.Error != nil {
				return result.Error
			}
			if result.Success {
				return nil
			}
		// If the predicate has not finished by the timeout, cancel it.
		case <-time.After(time.Duration(timeout) * time.Second):
			return fmt.Errorf("a timeout occurred")
		}
	}
}

// NormalizeURL is an internal function to be used by provider clients.
//
// It ensures that each endpoint URL has a closing `/`, as expected by
// ServiceClient's methods.
func NormalizeURL(url string) string {
	if !strings.HasSuffix(url, "/") {
		return url + "/"
	}
	return url
}

// NormalizePathURL is used to convert rawPath to a fqdn, using basePath as
// a reference in the filesystem, if necessary. basePath is assumed to contain
// either '.' when first used, or the file:// type fqdn of the parent resource.
// e.g. myFavScript.yaml => file://opt/lib/myFavScript.yaml
func NormalizePathURL(basePath, rawPath string) (string, error) {
	u, err := url.Parse(rawPath)
	if err != nil {
		return "", err
	}
	// if a scheme is defined, it must be a fqdn already
	if u.Scheme != "" {
		return u.String(), nil
	}
	// if basePath is a url, then child resources are assumed to be relative to it
	bu, err := url.Parse(basePath)
	if err != nil {
		return "", err
	}
	var basePathSys, absPathSys string
	if bu.Scheme != "" {
		basePathSys = filepath.FromSlash(bu.Path)
		absPathSys = filepath.Join(basePathSys, rawPath)
		bu.Path = filepath.ToSlash(absPathSys)
		return bu.String(), nil
	}

	absPathSys = filepath.Join(basePath, rawPath)
	u.Path = filepath.ToSlash(absPathSys)
	u.Scheme = "file"
	return u.String(), nil

}

// StripLastSlashURL removes last slash symbols from url.
func StripLastSlashURL(url string) string {
	if len(url) == 0 {
		return url
	}
	i := len(url) - 1
	for i >= 0 && url[i] == '/' {
		i--
	}
	return url[:i+1]
}

// NativeMapToStruct converts from map to struct.
func NativeMapToStruct(nativeMap interface{}, obj interface{}) error {

	var md mapstructure.Metadata
	config := &mapstructure.DecoderConfig{
		Metadata: &md,
		Result:   obj,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if err := decoder.Decode(nativeMap); err != nil {
		return err
	}
	return nil
}

func FailOnErrorF(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Fatalf("%s: %+v", fmt.Sprintf(msg, args...), err)
	}
}

func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
