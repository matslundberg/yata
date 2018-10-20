package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func assertCommandParams(t *testing.T,
	args string, err error,
	given_command Command, expected_command Command,
	given_data DataType, expected_data DataType,
	given_filter []Filter, expected_filter []Filter) {
	assert.Nil(t, err, fmt.Sprintf("parseCommand failed for %s", args))
	assert.Equal(t, given_command, expected_command, fmt.Sprintf("parseCommand failed to parse command for %s", args))
	assert.Equal(t, given_data, expected_data, fmt.Sprintf("parseCommand failed to parse data for %s", args))
	assert.ElementsMatch(t, given_filter, expected_filter, fmt.Sprintf("parseCommand didn't return expected filter for %s", args))
}

func TestParseCommand(t *testing.T) {
	args := ""

	args = "list tasks status is open"
	command, data, filter, err := parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("tasks"),
		filter, []Filter{newFilter(filterField_status, compType_exactMatch, "open")})

	args = "list tasks +tag"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("tasks"),
		filter, []Filter{newFilter(filterField_tag, compType_exactMatch, "+tag")})

	args = "list tasks @mention"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("tasks"),
		filter, []Filter{newFilter(filterField_mention, compType_exactMatch, "@mention")})

	args = "list open"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("tasks"),
		filter, []Filter{newFilter(filterField_status, compType_exactMatch, "open")})

	args = "list open @mention"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("tasks"),
		filter, []Filter{
			newFilter(filterField_status, compType_exactMatch, "open"),
			newFilter(filterField_mention, compType_exactMatch, "@mention"),
		})

	args = "list dummy123 @mention"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("tasks"),
		filter, []Filter{
			newFilter(filterField_mention, compType_exactMatch, "@mention"),
		})

	args = "list ongoing @mention"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("tasks"),
		filter, []Filter{
			newFilter(filterField_status, compType_exactMatch, "ongoing"),
			newFilter(filterField_mention, compType_exactMatch, "@mention"),
		})

	args = "list completed"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("tasks"),
		filter, []Filter{
			newFilter(filterField_status, compType_exactMatch, "completed"),
		})

	args = "list rejected"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("tasks"),
		filter, []Filter{
			newFilter(filterField_status, compType_exactMatch, "rejected"),
		})

	args = "list projects"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("projects"),
		filter, []Filter{})

	args = "list mentions"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("mentions"),
		filter, []Filter{})

	args = "list tags"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("list"),
		data, DataType("tags"),
		filter, []Filter{})

	args = "edit id:ABC12"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("edit"),
		data, DataType("tasks"),
		filter, []Filter{
			newFilter(filterField_id, compType_exactMatch, "ABC12"),
		})

	args = "complete open"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("complete"),
		data, DataType("tasks"),
		filter, []Filter{
			newFilter(filterField_status, compType_exactMatch, "open"),
		})

	args = "complete id:abc12"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("complete"),
		data, DataType("tasks"),
		filter, []Filter{
			newFilter(filterField_id, compType_exactMatch, "abc12"),
		})

	args = "complete id is abc12"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("complete"),
		data, DataType("tasks"),
		filter, []Filter{
			newFilter(filterField_id, compType_exactMatch, "abc12"),
		})

	args = "complete these"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("complete"),
		data, DataType("tasks"),
		filter, []Filter{
			newFilter(filterField_reference, compType_exactMatch, "these"),
		})

	//args = "add New thing todo @mention +tag"
	//command, data, filter, err = parseCommand(strings.Split(args, " "))
	//assertCommandParams(t, args, err,
	//	command, Command("add"), d
	//	ata, DataType("tasks"),
	//	filter, []Filter{Filter("New"), Filter("thing"), Filter("todo"), Filter("@mention"), Filter("+tag")})

	args = "reject id:ABC12"
	command, data, filter, err = parseCommand(strings.Split(args, " "))
	assertCommandParams(t, args, err,
		command, Command("reject"),
		data, DataType("tasks"),
		filter, []Filter{
			newFilter(filterField_id, compType_exactMatch, "ABC12"),
		})
}
