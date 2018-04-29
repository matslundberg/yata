package main

import (
    "fmt"
    "regexp"
    "strings"
    "github.com/logrusorgru/aurora"
    "crypto/sha256"
    "encoding/hex"
)

type TodoStatus int

func New(description string, sourceFile string) string {
    hash := sha256.Sum256([]byte(description+sourceFile))
    hashId := hex.EncodeToString( hash[:] )
    //ref := hashId[:8]
    return hashId
}

const (
    unknown TodoStatus = iota
    open 
    ongoing
    completed
    rejected
)

func TodoStatusToString(status TodoStatus) (string) {
    var constLookup = map[uint16]string{
        uint16(unknown): "unknown",
        uint16(open): "open",
        uint16(ongoing): "ongoing",
        uint16(completed): "completed",
        uint16(rejected): "rejected" }

    return constLookup[uint16(status)]
}

type Todo struct {
    status TodoStatus
    description string
    source string
    mentions []Mention
    tags []Tag
    _id dbEntryId
}

func (t Todo) id() dbEntryId {
    return t._id
}

func (todo Todo) print() {
    switch todo.status {
    case open:
        fmt.Println(todo.id(), "[ ] "+todo.description, aurora.Gray(todo.source))
    case completed:
        fmt.Println(todo.id(), aurora.Green("[x] "+todo.description), aurora.Gray(todo.source))
    case ongoing:
        fmt.Println(todo.id(), aurora.Brown("[/] "+todo.description), aurora.Gray(todo.source))
    case rejected:
        fmt.Println(todo.id(), aurora.Black("[-] "+todo.description), aurora.Gray(todo.source))
    }
}

func (t Todo) hasMention(mention Mention) (bool) {
    for _, m := range t.mentions {
        if(m == mention) {
            return true
        }
    }

    return false
}

func (t Todo) hasTag(tag Tag) (bool) {
    for _, t := range t.tags {
        if(t == tag) {
            return true
        }
    }

    return false
}

func (t Todo) filter(filter []string) (bool) {
    match := true

    re_mentions := regexp.MustCompile("^[@][a-zA-Z0-9]+")
    re_tags := regexp.MustCompile("[+][a-zA-Z0-9]+")

    for i := 0; i < len(filter); i++ {
        word := filter[i]

        switch {
        case word == "status":
            value := filter[i+2]
            i = i+2
            if(TodoStatusToString(t.status) != value) {
                match = false
            }
        case re_mentions.MatchString(word):
            value := word
            if(!t.hasMention(Mention{name: value})) {
                match = false
            }
        case re_tags.MatchString(word):
            value := word
            if(!t.hasTag(Tag{name: value})) {
                match = false
            }
        }
    }

    return match
}

func (t Todo) loadFromString(todoString string, sourceFile string) dbEntry {
    statusChar := ""
    description := ""
    i := strings.Index(todoString, "[")
    if i >= 0 {
        j := strings.Index(todoString[i:], "]")
        if j >= 0 {
            statusChar = strings.TrimSpace(todoString[i+1 : j-i])
            description = strings.TrimSpace(todoString[j+1:])
            //fmt.Println(description)
        }
    }

    status := unknown

    id := New(description, sourceFile)

    if(statusChar == "") {
        status = open
    } else if(statusChar == "x" || statusChar == "X") {
        status = completed
    } else if(statusChar == "/") {
        status = ongoing
    } else if(statusChar == "-") {
        status = rejected
    }

    mentions := make([]Mention, 0)
    mdt := MentionDataType{}
    for _, mention := range mdt.findString(description) {
        mentions = append(mentions, Mention{name: mention})
    }

    tags := make([]Tag, 0)
    tdt := TagDataType{}
    for _, tag := range tdt.findString(description) {
        tags = append(tags, Tag{name: tag})
    }

    ret := Todo{status: status, description: description, source: sourceFile, _id: dbEntryId(id), mentions: mentions, tags: tags}
    //fmt.Println(ret)
    return ret
}

type TodoDataType struct {
}

func (dt TodoDataType) findString(content string) []string {
    re := regexp.MustCompile("(\\[\\s*( |x|X|/|[-])\\s*\\](.+))")
    return re.FindAllString(content, -1)
}

func (dt TodoDataType) find(db NotesDatabase, filter []string) (map[dbEntryId]dbEntry) {
    todos := make(map[dbEntryId]dbEntry)

    for _, note := range db.notes {
        todoStrings :=  dt.findString(note.content)

        for _, todoString := range todoStrings {
            todo := Todo{}.loadFromString(todoString, note.filename)

            if(todo.filter(filter)) {
                todos[todo.id()] = todo
            }
        }
    }

    return todos
}
