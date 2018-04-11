package main

import (
    "fmt"
    "io/ioutil"
    "regexp"
    "log"
    "net/http"
    "strings"
//    "github.com/fatih/color"
//    "gopkg.in/h2non/filetype.v1"
    "github.com/logrusorgru/aurora"
)

type TodoStatus int

const (
    unknown TodoStatus = iota
    open 
    ongoing
    completed
    rejected
)

type Todo struct {
    status TodoStatus
    description string
    source string
    id string
    ref string
}

func NewTodoFromString(todoString string, sourceFile string) Todo {
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

    if(statusChar == "") {
        status = open
    } else if(statusChar == "x" || statusChar == "X") {
        status = completed
    } else if(statusChar == "/") {
        status = ongoing
    } else if(statusChar == "-") {
        status = rejected
    }

    ret := Todo{status: status, description: description, source: sourceFile, id: "", ref: ""}
    //fmt.Println(ret)
    return ret
}

func PrintTodo(todo Todo) {
    switch todo.status {
    case open:
        fmt.Println(" [ ] "+todo.description, aurora.Gray(todo.source))
    case completed:
        fmt.Println(aurora.Green(" [x] "+todo.description), aurora.Gray(todo.source))
    case ongoing:
        fmt.Println(aurora.Brown(" [/] "+todo.description), aurora.Gray(todo.source))
    case rejected:
        fmt.Println(aurora.Black(" [-] "+todo.description), aurora.Gray(todo.source))
    }
}

func main() {
    path := "/home/matslundberg/Dropbox/notes/";
    //path := "./tests/";
    
    files, err := ioutil.ReadDir(path)
    if err != nil {
        log.Fatal(err)
    }

    todos := make([]Todo, 0)

    for _, f := range files {
        //fmt.Println(f.Name())
        
        todo_file := path+f.Name()
        b, err := ioutil.ReadFile(todo_file) // just pass the file name
        if err != nil {
            fmt.Print(err)
        }

        contentType, err := GetFileContentType(b)
        if err != nil {
            panic(err)
        }

        //fmt.Println("Content Type: " + contentType)
        if strings.Contains(contentType, "text/plain") {
            //fmt.Println(b) // print the content as 'bytes'

            str := string(b) // convert content to a 'string'

            //fmt.Println(str) // print the content as a 'string'

            re := regexp.MustCompile("(\\[\\s*( |x|X|/|[-])\\s*\\](.+))")
            todoStrings := re.FindAllString(str, -1)
            //fmt.Println(todos)
            for _, todoString := range todoStrings {
                //fmt.Println( todoString )
                todo := NewTodoFromString(todoString, f.Name())
                todos = append(todos, todo)
            }

        }

    }

    for _, todo := range todos {
        if(todo.status != completed &&  todo.status != rejected) {
            PrintTodo(todo)
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
