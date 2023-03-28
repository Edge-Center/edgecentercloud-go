package types

import (
	"encoding/json"
	"fmt"
)

type Visibility string

// HwMachineType virtual chipset type.
type HwMachineType string

// SSHKeyType whether the image supports SSH key or not.
type SSHKeyType string

// OSType the operating system installed on the image.
type OSType string

// HwFirmwareType specifies the type of firmware with which to boot the guest.
type HwFirmwareType string

type ImageSourceType string

const (
	VisibilityPrivate Visibility = "private"
	VisibilityShared  Visibility = "shared"
	VisibilityPublic  Visibility = "public"

	HwMachineI440 HwMachineType = "i440"
	HwMachineQ35  HwMachineType = "q35"

	SSHKeyAllow    SSHKeyType = "allow"
	SSHKeyDeny     SSHKeyType = "deny"
	SSHKeyRequired SSHKeyType = "required"

	OsLinux   OSType = "linux"
	OsWindows OSType = "windows"

	HwFirmwareBIOS HwFirmwareType = "bios"
	HwFirmwareUEFI HwFirmwareType = "uefi"

	ImageSourceVolume ImageSourceType = "volume"
)

func (v Visibility) IsValid() error {
	switch v {
	case VisibilityPrivate, VisibilityShared, VisibilityPublic:
		return nil
	}
	return fmt.Errorf("invalid Visibility type: %v", v)
}

func (v Visibility) ValidOrNil() (*Visibility, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v Visibility) String() string {
	return string(v)
}

func (v Visibility) List() []Visibility {
	return []Visibility{VisibilityPrivate, VisibilityShared, VisibilityPublic}
}

func (v Visibility) StringList() []string {
	lst := v.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}
	return strings
}

// UnmarshalJSON - implements Unmarshaler interface for Visibility.
func (v *Visibility) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := Visibility(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for Visibility.
func (v *Visibility) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v HwMachineType) IsValid() error {
	switch v {
	case HwMachineI440, HwMachineQ35:
		return nil
	}
	return fmt.Errorf("invalid HwMachineType type: %v", v)
}

func (v HwMachineType) ValidOrNil() (*HwMachineType, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v HwMachineType) String() string {
	return string(v)
}

func (v HwMachineType) List() []HwMachineType {
	return []HwMachineType{HwMachineI440, HwMachineQ35}
}

func (v HwMachineType) StringList() []string {
	lst := v.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}
	return strings
}

// UnmarshalJSON - implements Unmarshaler interface for HwMachineType.
func (v *HwMachineType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := HwMachineType(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for HwMachineType.
func (v *HwMachineType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v SSHKeyType) IsValid() error {
	switch v {
	case SSHKeyAllow, SSHKeyDeny, SSHKeyRequired:
		return nil
	}
	return fmt.Errorf("invalid SSHKeyType type: %v", v)
}

func (v SSHKeyType) ValidOrNil() (*SSHKeyType, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v SSHKeyType) String() string {
	return string(v)
}

func (v SSHKeyType) List() []SSHKeyType {
	return []SSHKeyType{SSHKeyAllow, SSHKeyDeny, SSHKeyRequired}
}

func (v SSHKeyType) StringList() []string {
	lst := v.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}
	return strings
}

// UnmarshalJSON - implements Unmarshaler interface for SSHKeyType.
func (v *SSHKeyType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := SSHKeyType(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for SSHKeyType.
func (v *SSHKeyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v OSType) IsValid() error {
	switch v {
	case OsLinux, OsWindows:
		return nil
	}
	return fmt.Errorf("invalid OSType type: %v", v)
}

func (v OSType) ValidOrNil() (*OSType, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v OSType) String() string {
	return string(v)
}

func (v OSType) List() []OSType {
	return []OSType{OsLinux, OsWindows}
}

func (v OSType) StringList() []string {
	lst := v.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}
	return strings
}

// UnmarshalJSON - implements Unmarshaler interface for OSType.
func (v *OSType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := OSType(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for OSType.
func (v *OSType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v HwFirmwareType) IsValid() error {
	switch v {
	case HwFirmwareBIOS, HwFirmwareUEFI:
		return nil
	}
	return fmt.Errorf("invalid HwFirmwareType type: %v", v)
}

func (v HwFirmwareType) ValidOrNil() (*HwFirmwareType, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v HwFirmwareType) String() string {
	return string(v)
}

func (v HwFirmwareType) List() []HwFirmwareType {
	return []HwFirmwareType{HwFirmwareBIOS, HwFirmwareUEFI}
}

func (v HwFirmwareType) StringList() []string {
	lst := v.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}
	return strings
}

// UnmarshalJSON - implements Unmarshaler interface for HwFirmwareType.
func (v *HwFirmwareType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := HwFirmwareType(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for HwFirmwareType.
func (v *HwFirmwareType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v ImageSourceType) IsValid() error {
	if v == ImageSourceVolume {
		return nil
	}
	return fmt.Errorf("invalid ImageSourceType type: %v", v)
}

func (v ImageSourceType) ValidOrNil() (*ImageSourceType, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v ImageSourceType) String() string {
	return string(v)
}

func (v ImageSourceType) List() []ImageSourceType {
	return []ImageSourceType{ImageSourceVolume}
}

func (v ImageSourceType) StringList() []string {
	lst := v.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}
	return strings
}

// UnmarshalJSON - implements Unmarshaler interface for ImageSourceType.
func (v *ImageSourceType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := ImageSourceType(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for ImageSourceType.
func (v *ImageSourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}
