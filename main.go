package main

import (
    "fmt"
    "os"
)


func parseCommand(command []string) (string, string, []string) {
    cmd := command[0]
    data := command[1]
    filter := command[2:]
    return cmd, data, filter
}


func main() {
    //path := "/home/matslundberg/Dropbox/notes/";
    path := "./tests/";

    db := LoadDatabase(path)

    args := os.Args[1:]
    command, data, filter := parseCommand(args)
    fmt.Println(command, data, filter)

    switch(command) {
    case "list":
        list := FindInDatabase(db, data, filter)
        //fmt.Println(db)
        for _, entry := range list {
            entry.print()
        }
    }
    
}
