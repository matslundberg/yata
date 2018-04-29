package main

import (
    "io/ioutil"
    "log"
    "strings"
)

type Note struct {
    filename string
    content string
    contentType string
}

type NotesDatabase struct {
    path string
    notes []Note
}

func LoadDatabase(path string) (NotesDatabase) {
    notes := make([]Note,0)

    files, err := ioutil.ReadDir(path)
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        //fmt.Println(f.Name())
        
        todo_file := path+f.Name()
        b, err := ioutil.ReadFile(todo_file) // just pass the file name
        if err != nil {
            //fmt.Print(err)
            continue
        }

        contentType, err := GetFileContentType(b)
        if err != nil {
            panic(err)
        }

        //fmt.Println("Content Type: " + contentType)
        if strings.Contains(contentType, "text/plain") {
            str := string(b) // convert content to a 'string'
            note := Note{filename: todo_file, content: str, contentType: contentType}

            notes = append(notes, note)
        }
    }

    return NotesDatabase{path: path,notes: notes}
}

type dbEntryId string

type dbEntry interface {
    print()
    filter([]string) bool
    loadFromString(string, string) dbEntry
    id() dbEntryId
}

type dbDataType interface {
    find(db NotesDatabase, filter []string) (map[dbEntryId]dbEntry)
    findString(content string) []string
}

func LoadDataType(data string) dbDataType {
    switch(data) {
    case "tasks":
        return TodoDataType{}
    case "tags":
        return TagDataType{}
    case "projects", "mentions":
        return MentionDataType{}
    }

    return nil
}

func FindInDatabase(db NotesDatabase, data string, filter []string) (map[dbEntryId]dbEntry) {
    dt := LoadDataType(data)
    return dt.find(db, filter)
}

