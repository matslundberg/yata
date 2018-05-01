package main

import (
	"fmt"
	"regexp"
)

type Tag struct {
	name string
}

func (t Tag) Id() dbEntryId {
	return dbEntryId(t.name)
}

func (t Tag) Source() string {
	return ""
}

func (t Tag) LineNum() int {
	return 0
}

func (t Tag) update(command string, value string) dbEntry {
	return t
}

func (t Tag) print() {
	fmt.Println(t.name)
}

func (t Tag) toString() string {
	return t.name
}

func (t Tag) filter(filter []string) bool {
	match := true

	return match
}

func (t Tag) loadFromString(content string, sourceFile string, lineNum int) dbEntry {
	ret := Tag{name: content}
	return ret
}

type TagDataType struct {
}

func (dt TagDataType) findString(content string) []string {
	re := regexp.MustCompile("[+][a-zA-Z0-9]+")
	return re.FindAllString(content, -1)
}

func (dt TagDataType) find(db NotesDatabase, filter []string) dbResultSet {
	tags := make(dbResultSet)

	for _, note := range db.notes {
		tagsStrings := dt.findString(note.content)

		for _, tagString := range tagsStrings {
			tag := Tag{}.loadFromString(tagString, note.filename, 0)

			if tag.filter(filter) {
				tags[tag.Id()] = tag
			}
		}
	}

	return tags
}

func (dt TagDataType) findById(db NotesDatabase, id dbEntryId) dbEntry {
	panic("Not implemented")
	return nil
}
