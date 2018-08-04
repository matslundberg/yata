package main

import (
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
	"fmt"
)


func TestParseCommand(t *testing.T) {
	args := ""
	assert := assert.New(t)

	args = "list tasks status is open"
	command, data, filter, err := parseCommand(strings.Split(args, " "))
	assert.Nil(err, fmt.Sprintf("parseCommand failed for %s", args))
	assert.Equal(command, Command("list"), "parseCommand failed to parse command")
	assert.Equal(data, DataType("tasks"), "parseCommand failed to parse data")
	assert.Equal(len(filter), 3, "parseCommand failed to parse filter")

	args = "list tasks +tag"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assert.Nil(err, fmt.Sprintf("parseCommand failed for %s", args))
	assert.Equal(command, Command("list"), "parseCommand failed to parse command")
	assert.Equal(data, DataType("tasks"), "parseCommand failed to parse data")
	assert.Equal(len(filter), 1, "parseCommand failed to parse filter")

	args = "list tasks @mention"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assert.Nil(err, fmt.Sprintf("parseCommand failed for %s", args))
	assert.Equal(command, Command("list"), "parseCommand failed to parse command")
	assert.Equal(data, DataType("tasks"), "parseCommand failed to parse data")
	assert.Equal(len(filter), 1, "parseCommand failed to parse filter")
}