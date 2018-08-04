package main

import (
	"fmt"
	"os"
)

func printHelp() {
	fmt.Printf(
`
muistiin $action $type $filter

Example commands
muistiin list tasks status is open +tag @mention
muistiin list tags
muistiin list mentions
muistiin list projects
muistiin complete tasks status is open
muistiin complete tasks these // References previous search result

NOTE! Projects and mentions are the same thing...
`)
}

type Command string
type DataType string
type Filter string

func parseCommand(command []string) (Command, DataType, []Filter, error) {
	if(len(command) > 2) {
		cmd := Command(command[0])
		data := DataType(command[1])
		filters := make([]Filter, len(command)-2)
		for k, filter := range command[2:] {
			filters[k-2] = Filter(filter)
		}
		return cmd, data, filters, nil
	} else {
		return Command(""), DataType(""), make([]Filter, 0), fmt.Errorf("Failed to parse command")
	}
}

func run() error {
	path := os.Getenv("MIIN_PATH")

	if path == "" {
		return fmt.Errorf("Env variable MIIN_PATH not set")
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
