package stacks

import (
	"fmt"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

type InvalidEnvironmentError struct {
	edgecloud.BaseError
	Section string
}

func (e InvalidEnvironmentError) Error() string {
	return fmt.Sprintf("environment has wrong section: %s", e.Section)
}

type InvalidDataFormatError struct {
	edgecloud.BaseError
}

func (e InvalidDataFormatError) Error() string {
	return "data in neither json nor yaml format"
}

type InvalidTemplateFormatVersionError struct {
	edgecloud.BaseError
	Version string
}

func (e InvalidTemplateFormatVersionError) Error() string {
	return "template format version not found"
}

type TemplateRequiredError struct {
	edgecloud.BaseError
}

func (e TemplateRequiredError) Error() string {
	return "template required for this function"
}
