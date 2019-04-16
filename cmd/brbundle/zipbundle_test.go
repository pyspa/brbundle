package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanPath(t *testing.T) {
	testcases := []struct {
		name   string
		prefix string
		file   string
		result string
	}{
		{"only file", "", "test.png", "test.png"},
		{"only file with dir", "", "dir/test.png", "dir/test.png"},
		{"prefix and file", "dir", "test.png", "dir/test.png"},
		{"prefix and file with dir", "dir1", "dir2/test.png", "dir1/dir2/test.png"},
		{"windows path", "dir1", "dir2\\test.png", "dir1/dir2/test.png"},
		{"omit root slash", "/dir1", "dir2/test.png", "dir1/dir2/test.png"},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			result := cleanPath(testcase.prefix, testcase.file)
			assert.Equal(t, testcase.result, result)
		})
	}
}
