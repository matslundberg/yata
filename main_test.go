package main

import (
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func assertCommandParams(t *testing.T, 
		args string, err error,
		given_command Command, expected_command Command, 
		given_data DataType, expected_data DataType,
		given_filter []Filter, expected_filter []Filter) {
	assert.Nil(t, err, fmt.Sprintf("parseCommand failed for %s", args))
	assert.Equal(t, given_command, expected_command, "parseCommand failed to parse command")
	assert.Equal(t, given_data, expected_data, "parseCommand failed to parse data")
	assert.Equal(t, len(given_filter), len(expected_filter), "parseCommand failed to parse filter")
	assert.ElementsMatch(t, given_filter, expected_filter, "parseCommand didn't return expected filter");
}

func TestParseCommand(t *testing.T) {
	args := ""

	args = "list tasks status is open"
	command, data, filter, err := parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("list"), data, DataType("tasks"), filter, []Filter{Filter("status"), Filter("is"), Filter("open")})

	args = "list tasks +tag"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("list"), data, DataType("tasks"), filter, []Filter{Filter("+tag")})

	args = "list tasks @mention"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("list"), data, DataType("tasks"), filter, []Filter{Filter("@mention")})

	args = "list open"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("list"), data, DataType("tasks"), filter, []Filter{Filter("open")})

	args = "list open @mention"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("list"), data, DataType("tasks"), filter, []Filter{Filter("open"), Filter("@mention")})

	args = "list ongoing @mention"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("list"), data, DataType("tasks"), filter, []Filter{Filter("ongoing"), Filter("@mention")})

	args = "list completed"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("list"), data, DataType("tasks"), filter, []Filter{Filter("completed")})

	args = "list rejected"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("list"), data, DataType("tasks"), filter, []Filter{Filter("rejected")})

	args = "list projects"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("list"), data, DataType("projects"), filter, []Filter{})

	args = "list mentions"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("list"), data, DataType("mentions"), filter, []Filter{})

	args = "list tags"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("list"), data, DataType("tags"), filter, []Filter{})

	args = "edit $task_id"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("edit"), data, DataType("tasks"), filter, []Filter{Filter("$task_id")})

	args = "complete open"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("complete"), data, DataType("tasks"), filter, []Filter{Filter("open")})

	args = "complete these"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("complete"), data, DataType("tasks"), filter, []Filter{Filter("these")})

	args = "add New thing todo @mention +tag"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("add"), data, DataType("tasks"), filter, []Filter{Filter("New"), Filter("thing"), Filter("todo"), Filter("@mention"), Filter("+tag")})

	args = "reject $task_id"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err, command, Command("reject"), data, DataType("tasks"), filter, []Filter{Filter("$task_id")})
}
