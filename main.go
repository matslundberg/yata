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
type TodoId struct {
    hashId string
    ref string
}

func New(description string, sourceFile string) TodoId {
    hash := sha256.Sum256([]byte(description+sourceFile))
    hashId := hex.EncodeToString( hash[:] )
    ref := hashId[:8]
    return TodoId{hashId, ref}
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
    id TodoId
}

func (t Todo) print() {
    PrintTodo(t)
}

func (t Todo) filter(filter []string) (bool) {
    match := true

    for i := 0; i < len(filter); i++ {
        word := filter[i]

        switch(word) {
        case "status":
            value := filter[i+2]
            i = i+2
            if(TodoStatusToString(t.status) != value) {
                match = false
            }
        }
    }

    return match
}

func LoadTodoFromString(todoString string, sourceFile string) Todo {
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

    ret := Todo{status: status, description: description, source: sourceFile, id: id}
    //fmt.Println(ret)
    return ret
}

func PrintTodo(todo Todo) {
    switch todo.status {
    case open:
        fmt.Println(todo.id.ref, "[ ] "+todo.description, aurora.Gray(todo.source))
    case completed:
        fmt.Println(todo.id.ref, aurora.Green("[x] "+todo.description), aurora.Gray(todo.source))
    case ongoing:
        fmt.Println(todo.id.ref, aurora.Brown("[/] "+todo.description), aurora.Gray(todo.source))
    case rejected:
        fmt.Println(todo.id.ref, aurora.Black("[-] "+todo.description), aurora.Gray(todo.source))
    }
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

type dbEntry interface {
    print()
    filter([]string) bool
}

func FindInDatabase(db NotesDatabase, data string, filter []string) ([]dbEntry) {

    todos := make([]dbEntry, 0)

    for _, note := range db.notes {
        re := regexp.MustCompile("(\\[\\s*( |x|X|/|[-])\\s*\\](.+))")
        todoStrings := re.FindAllString(note.content, -1)

        for _, todoString := range todoStrings {
            todo := LoadTodoFromString(todoString, note.filename)

            if(todo.filter(filter)) {
                todos = append(todos, todo)
            }
        }
    }

    return todos
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
