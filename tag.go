package main

import (
    "fmt"
    "regexp"
)

type Tag struct {
    name string
}

func (t Tag) id() dbEntryId {
    return dbEntryId(t.name)
}

func (t Tag) print() {
    fmt.Println(t.name)
}

func (t Tag) filter(filter []string) (bool) {
    match := true

    return match
}

func (t Tag) loadFromString(content string, sourceFile string) dbEntry {
    ret := Tag{name: content}
    return ret
}

type TagDataType struct {
}

func (dt TagDataType) findString(content string) []string {
    re := regexp.MustCompile("[+][a-zA-Z0-9]+")
    return re.FindAllString(content, -1)
}

func (dt TagDataType) find(db NotesDatabase, filter []string) (map[dbEntryId]dbEntry) {
    tags := make(map[dbEntryId]dbEntry)

    for _, note := range db.notes {
        tagsStrings := dt.findString(note.content)

        for _, tagString := range tagsStrings {
            tag := Tag{}.loadFromString(tagString, note.filename)

            if(tag.filter(filter)) {
                tags[tag.id()] = tag
            }
        }
    }

    return tags
}
