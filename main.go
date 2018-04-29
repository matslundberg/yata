package main

import (
    "fmt"
    "io/ioutil"
    "regexp"
    "log"
    "net/http"
    "strings"
    "github.com/logrusorgru/aurora"
    "crypto/sha256"
    "encoding/hex"
    "os"
)

type TodoStatus int
// type TodoId struct {
//     hashId string
//     ref string
// }

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

type Mention struct {
    name string
}

func (m Mention) id() dbEntryId {
    return dbEntryId(m.name)
}

func (t Mention) print() {
    fmt.Println(t.name)
}

func (t Mention) filter(filter []string) (bool) {
    match := true

    return match
}

func (t Mention) loadFromString(content string, sourceFile string) dbEntry {
    ret := Mention{name: content}
    return ret
}

func parseCommand(command []string) (string, string, []string) {
    cmd := command[0]
    data := command[1]
    filter := command[2:]
    return cmd, data, filter
}

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

type MentionDataType struct {
}

func (dt MentionDataType) findString(content string) []string {
    re := regexp.MustCompile("[@][a-zA-Z0-9]+")
    return re.FindAllString(content, -1)
}


func (dt MentionDataType) find(db NotesDatabase, filter []string) (map[dbEntryId]dbEntry) {
    mentions := make(map[dbEntryId]dbEntry)

    for _, note := range db.notes {
        mentionStrings := dt.findString(note.content)

        for _, mentionString := range mentionStrings {
            mention := Mention{}.loadFromString(mentionString, note.filename)

            if(mention.filter(filter)) {
                mentions[mention.id()] = mention;
            }
        }
    }

    return mentions
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

func GetFileContentType(buffer []byte) (string, error) {

    // Only the first 512 bytes are used to sniff the content type.
    //buffer := make([]byte, 512)

    //_, err := out.Read(buffer)
    //if err != nil {
    //    return "", err
    //}

    // Use the net/http package's handy DectectContentType function. Always returns a valid 
    // content-type by returning "application/octet-stream" if no others seemed to match.
    contentType := http.DetectContentType(buffer)

    return contentType, nil
}
