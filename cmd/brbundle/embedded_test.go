package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSplitString(t *testing.T) {
	result := splitByte([]byte("abcd\x00efghijklmn"), 8)
	assert.Equal(t, 2, len(result))
}

func TestSplitString2(t *testing.T) {
	// escape sequence is on the line limit
	result := splitByte([]byte("abcdef\x00ghijklmn"), 8)
	assert.Equal(t, 2, len(result))
}

func TestSplitString3(t *testing.T) {
	// escape sequence is the last
	result := splitByte([]byte("abcdefghijklmn\x00"), 8)
	assert.Equal(t, 2, len(result))
}

func TestSplitString4(t *testing.T) {
	// too short
	result := splitByte([]byte("abc"), 8)
	assert.Equal(t, 1, len(result))
}

func TestSplitString5(t *testing.T) {
	// empty
	result := splitByte([]byte(""), 8)
	assert.Equal(t, 1, len(result))
}

func TestSplitString6(t *testing.T) {
	// other escape sequences
	result := splitByte([]byte("abc\"def"), 8)
	assert.Equal(t, 1, len(result))
}

func TestSplitString7(t *testing.T) {
	// other escape sequences on the end of line
	result := splitByte([]byte("abcdefg\nhijk"), 8)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "hijk", result[1])
}

func TestFormatContent(t *testing.T) {
	result := formatContent([]byte("abc\"\\\\def"), 7)
	assert.True(t, strings.Contains(result, `\"`))
	assert.True(t, strings.Contains(result, `\\\\`))
}

func TestFormatContent2(t *testing.T) {
	result := formatContent([]byte("abcdefg\nhijk"), 7)
	assert.True(t, strings.Contains(result, `\n`))
	assert.True(t, strings.Contains(result, `hijk`))
}
