package main

import (
	"fmt"
    "io/ioutil"
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

type dbEntryId string

type dbEntry interface {
    print()
    filter([]string) bool
    loadFromString(string, string) dbEntry
    Id() dbEntryId
}

type dbDataType interface {
    find(db NotesDatabase, filter []string) (map[dbEntryId]dbEntry)
    findString(content string) []string
}

func LoadDatabase(path string) (NotesDatabase, error) {
    notes := make([]Note,0)

    files, err := ioutil.ReadDir(path)
    if err != nil {
        return NotesDatabase{}, fmt.Errorf("Failed to open directory %s <= %s", path, err)
    }

    for _, f := range files {
        
        filename := path+f.Name()
        b, err := ioutil.ReadFile(filename)
        if err != nil {
            // Just ignore files which cannot be read.
            continue
        }

        contentType, err := GetFileContentType(b)
        if err != nil {
            panic(err)
        }

        if strings.Contains(contentType, "text/plain") {
            str := string(b) // convert content to a 'string'
            note := Note{filename: filename, content: str, contentType: contentType}

            notes = append(notes, note)
        }
    }

    return NotesDatabase{path: path,notes: notes}, nil
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

func (db NotesDatabase) find(data string, filter []string) (map[dbEntryId]dbEntry) {
    dt := LoadDataType(data)
    return dt.find(db, filter)
}

