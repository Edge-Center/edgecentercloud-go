package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString_ErrorNoChars(t *testing.T) {
	result, err := GenerateRandomString(5, false, false, false)
	assert.Error(t, err)
	assert.Equal(t, "no characters specified", err.Error())
	assert.Empty(t, result)
}

func TestGenerateRandomString_OnlyLetters(t *testing.T) {
	result, err := GenerateRandomString(10, true, false, false)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(result))
	for _, c := range result {
		assert.Contains(t, letters, string(c))
	}
}

func TestGenerateRandomString_OnlyDigits(t *testing.T) {
	result, err := GenerateRandomString(10, false, true, false)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(result))
	for _, c := range result {
		assert.Contains(t, digits, string(c))
	}
}

func TestGenerateRandomString_OnlySpecial(t *testing.T) {
	result, err := GenerateRandomString(10, false, false, true)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(result))
	for _, c := range result {
		assert.Contains(t, special, string(c))
	}
}

func TestGenerateRandomString_Combined(t *testing.T) {
	result, err := GenerateRandomString(10, true, true, true)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(result))
	chars := letters + digits + special
	for _, c := range result {
		assert.Contains(t, chars, string(c))
	}
}
