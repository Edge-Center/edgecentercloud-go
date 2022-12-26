package stacks

import (
	"fmt"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

type ErrInvalidEnvironment struct {
	edgecloud.BaseError
	Section string
}

func (e ErrInvalidEnvironment) Error() string {
	return fmt.Sprintf("environment has wrong section: %s", e.Section)
}

type ErrInvalidDataFormat struct {
	edgecloud.BaseError
}

func (e ErrInvalidDataFormat) Error() string {
	return "data in neither json nor yaml format"
}

type ErrInvalidTemplateFormatVersion struct {
	edgecloud.BaseError
	Version string
}

func (e ErrInvalidTemplateFormatVersion) Error() string {
	return "template format version not found"
}

type ErrTemplateRequired struct {
	edgecloud.BaseError
}

func (e ErrTemplateRequired) Error() string {
	return "template required for this function"
}
