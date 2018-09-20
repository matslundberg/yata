package main

import (
	"fmt"
	"os"
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
type Filter string

func parseCommand(command []string) (Command, DataType, []Filter, error) {
	if(len(command) >= 2) {
		cmd := Command(command[0])
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

		filters := make([]Filter, len(command)-filter_count)
		for k, filter := range command[filter_count:] {
			//fmt.Println(k, filter_count)
			filters[k] = Filter(filter)
		}
		return cmd, data, filters, nil
	} else {
		return Command(""), DataType(""), make([]Filter, 0), fmt.Errorf("Failed to parse command")
	}
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
	command, data, filter, err := parseCommand(args)
	if err != nil {
		printHelp();
		return nil;
		//return fmt.Errorf("Failed to parse commmand", args)
	}

	fmt.Println("Running command", command, data, filter)

	switch command {
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
