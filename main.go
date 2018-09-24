package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func printHelp() {
	fmt.Printf(
		`
yata $action $type $filter

Example commands
yata list tasks status is open +tag @mention
yata list tags
yata list mentions
yata list projects
yata complete tasks status is open
yata complete tasks these // References previous search result


Better alternative?
yata list open
yata list open @mention
yata list ongoing @mention
yata list completed
yata list rejected
yata list projects
yata list mentions
yata list tags
yata edit $task_id
yata complete open
yata complete these
yata add New thing todo @mention +tag   --> Appends to yyyy-mm-dd.todo
yata reject $task_id
yata 

NOTE! Projects and mentions are the same thing...
`)
}

type Command string
type DataType string
type Content string
type FilterField int

const (
	filterField_id FilterField = iota
	filterField_status
	filterField_mention
	filterField_tag
	filterField_reference
)

type ComparisonType int

const (
	compType_exactMatch ComparisonType = iota
	compType_noMatch
	compType_partlyMatch
)

type Filter struct {
	field FilterField
	comp  ComparisonType
	value string
}

func newFilter(field FilterField, comp ComparisonType, value string) Filter {
	return Filter{
		field: field,
		comp:  comp,
		value: value,
	}
}

func parseCommand(command []string) (Command, DataType, []Filter, Content, error) {
	if len(command) >= 2 {
		cmd := Command(command[0])

		if(cmd == "add") {
			data := DataType("tasks")
			content := Content( strings.Join(command[1:], " ") )
			return cmd, data, make([]Filter, 0), content, nil
		}

		data := DataType(command[1])

		filter_count := 2
		//fmt.Println(command)
		switch data {
		case "tasks", "projects", "mentions", "tags":
			// do nothing
		default:
			data = DataType("tasks")
			filter_count = 1
		}

		filters_strings := make([]string, len(command)-filter_count)
		for k, filter := range command[filter_count:] {
			filters_strings[k] = filter
		}
		filters, err := parseFilters(filters_strings)
		if err != nil {
			return Command(""), DataType(""), make([]Filter, 0), Content(""), fmt.Errorf("Failed to parse filters")
		}
		return cmd, data, filters, Content(""), nil
	} else {
		return Command(""), DataType(""), make([]Filter, 0), Content(""), fmt.Errorf("Failed to parse command")
	}
}

func parseFilters(strings []string) ([]Filter, error) {
	re_mentions := regexp.MustCompile("^[@][a-zA-Z0-9]+")
	re_tags := regexp.MustCompile("^[+][a-zA-Z0-9]+")
	re_todo_id := regexp.MustCompile("^(id\\:)[a-zA-Z0-9]{5}")
	filters := make([]Filter, 0)

	for i := 0; i < len(strings); i++ {
		word := strings[i]
		switch {
		case word == "id":
			value := strings[i+2]
			i = i + 2

			filters = append(filters, Filter{
				field: filterField_id,
				comp:  compType_exactMatch,
				value: value,
			})

		case word == "status":
			value := strings[i+2]
			i = i + 2

			filters = append(filters, Filter{
				field: filterField_status,
				comp:  compType_exactMatch,
				value: value,
			})
		case re_mentions.MatchString(word):
			filters = append(filters, Filter{
				field: filterField_mention,
				comp:  compType_exactMatch,
				value: word,
			})

		case re_tags.MatchString(word):
			filters = append(filters, Filter{
				field: filterField_tag,
				comp:  compType_exactMatch,
				value: word,
			})
		case re_todo_id.MatchString(word):
			filters = append(filters, Filter{
				field: filterField_id,
				comp:  compType_exactMatch,
				value: word[3:],
			})
		case StringToTodoStatus(word) != unknown:
			filters = append(filters, Filter{
				field: filterField_status,
				comp:  compType_exactMatch,
				value: word,
			})
		case word == "these":
			filters = append(filters, Filter{
				field: filterField_reference,
				comp:  compType_exactMatch,
				value: word,
			})
		}
	}

	return filters, nil

}

func run() error {
	path := os.Getenv("YATA_PATH")

	if path == "" {
		return fmt.Errorf("Env variable YATA_PATH not set")
	}

	db, err := LoadDatabase(path)
	if err != nil {
		return fmt.Errorf("Failed to open database <= %s", err)
	}

	args := os.Args[1:]
	command, data, filter, content, err := parseCommand(args)
	if err != nil {
		fmt.Println(args, err)
		printHelp()
		return nil
		//return fmt.Errorf("Failed to parse commmand", args)
	}

	fmt.Println("Running command", command, data, filter)

	switch command {
	case "add":
		err := db.add(data, content)
		if err != nil {
			return fmt.Errorf("Failed to add %s <= %s", data, err)
		} else {
			fmt.Printf("New todo added")
		}
	case "list":
		list, err := db.find(data, filter)
		if err != nil {
			return fmt.Errorf("Failed to run find on database <= %s", err)
		}

		for _, entry := range list {
			entry.print()
		}

		db.saveResultSet(list)
	case "complete":
		list, err := db.find(data, filter)
		if err != nil {
			return fmt.Errorf("Failed to run find on database <= %s", err)
		}

		for _, entry := range list {
			if err := db.update(entry, command); err != nil {
				return fmt.Errorf("Failed to update dbEntry with %s <= %s", command, err)
			} else {
				entry.print()
			}
		}
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
