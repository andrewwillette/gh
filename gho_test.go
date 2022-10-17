package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUrlFromGitRemote(t *testing.T) {
	res := getUrlFromGitRemote()
	println(res)
	assert.Equal(t, "http://github.com/andrewwillette/gho", res)
}
