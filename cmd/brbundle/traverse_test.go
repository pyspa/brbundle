package main

import (
	"testing"
)

func TestSplitBuildTag(t *testing.T) {
	tests := []struct {
		name           string
		fileName       string
		buildTag       string
		wantPsuedoName string
		wantMatch      bool
	}{
		{
			"regular file with no build tag",
			"test",
			"",
			"test",
			true,
		},
		{
			"regular file with build tag and matched",
			"test__linux.txt",
			"linux",
			"test.txt",
			true,
		},
		{
			"regular file with build tag and not matched",
			"test__linux.txt",
			"darwin",
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPsuedoName, gotMatch := splitBuildTag(tt.fileName, tt.buildTag)
			if gotPsuedoName != tt.wantPsuedoName {
				t.Errorf("splitBuildTag() gotPsuedoName = %v, want %v", gotPsuedoName, tt.wantPsuedoName)
			}
			if gotMatch != tt.wantMatch {
				t.Errorf("splitBuildTag() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
		})
	}
}
