package main

import (
	"fmt"
	"regexp"
)

type Mention struct {
	name string
}

func (m Mention) Id() dbEntryId {
	return dbEntryId(m.name)
}

func (m Mention) Source() string {
	return ""
}

func (m Mention) LineNum() int {
	return 0
}

func (m Mention) update(command string, value string) dbEntry {
	return m
}

func (t Mention) print() {
	fmt.Println(t.name)
}

func (t Mention) toString() string {
	return t.name
}

func (t Mention) filter(filter []string) bool {
	match := true

	return match
}

func (t Mention) loadFromString(content string, sourceFile string, lineNum int) dbEntry {
	ret := Mention{name: content}
	return ret
}

type MentionDataType struct {
}

func (dt MentionDataType) findString(content string) []string {
	re := regexp.MustCompile("[@][a-zA-Z0-9]+")
	return re.FindAllString(content, -1)
}

func (dt MentionDataType) find(db NotesDatabase, filter []string) dbResultSet {
	mentions := make(dbResultSet)

	for _, note := range db.notes {
		mentionStrings := dt.findString(note.content)

		for _, mentionString := range mentionStrings {
			mention := Mention{}.loadFromString(mentionString, note.filename, 0)

			if mention.filter(filter) {
				mentions[mention.Id()] = mention
			}
		}
	}

	return mentions
}

func (dt MentionDataType) findById(db NotesDatabase, id dbEntryId) dbEntry {
	panic("Not implemented")
	return nil
}
