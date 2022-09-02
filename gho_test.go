package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGithub(t *testing.T) {
	// var tests = []struct {
	// 	provided string
	// 	expected string
	// }{
	// 	{
	// 		provided: "andrewwillette.com",
	// 		expected: "andrewwillette.com",
	// 	},
	// }

	// for _, tt := range tests {

	// }

}

func TestGetUrlFromGitRemote(t *testing.T) {
	res := getUrlFromGitRemote()
	println(res)
	assert.Equal(t, "http://github.com/andrewwillette/gho", res)
}
