package main

import (
	"fmt"
	"os"
)

func parseCommand(command []string) (string, string, []string, error) {
	cmd := command[0]
	data := command[1]
	filter := command[2:]
	return cmd, data, filter, nil
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
		return fmt.Errorf("Failed to parse commmand", args)
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
