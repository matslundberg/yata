package main

import (
	"fmt"
    "io/ioutil"
    "strings"
    "encoding/json"
    "os"
)

const (
	RS_FILE = "/tmp/miin.rs_file"
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

type dbResultSet map[dbEntryId]dbEntry

type dbEntry interface {
    print()
    toString() string
    filter([]string) bool
    update(string, string) dbEntry
    loadFromString(string, string, int) dbEntry
    Id() dbEntryId
    Source() string
    LineNum() int
}

type dbDataType interface {
    find(db NotesDatabase, filter []string) (dbResultSet)
    findString(content string) []string
    findById(db NotesDatabase, id dbEntryId) dbEntry
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

func (db NotesDatabase) find(data string, filter []string) (dbResultSet, error) {
    dt := LoadDataType(data)
    // TODO Fix this hack
    if dt == nil && data == "these" {
    	rs, err := db.getResultSet()
    	if err != nil {
    		return nil, fmt.Errorf("Failed to get resultSet <= %s", err)
    	}
    	return rs, nil
    }
    return dt.find(db, filter), nil
}

func (db NotesDatabase) getResultSet() (dbResultSet, error) {
	if _, err := os.Stat(RS_FILE); os.IsNotExist(err) {
		return make(dbResultSet, 0), nil
	}

	//fmt.Println(resultSet)
	data, err := ioutil.ReadFile(RS_FILE)

	if err != nil {
		return make(dbResultSet, 0), fmt.Errorf("Failed to read from resultSet file <= %s", err)
	}

	var unmarshalled map[string]interface{}

	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		//fmt.Println(err)
		return make(dbResultSet, 0), fmt.Errorf("Failed to deserialize resultSet file <= %s", err)
	}

	ret := make(dbResultSet, 0)
	for id, _ := range unmarshalled {
		ret[dbEntryId(id)] = TodoDataType{}.findById(db, dbEntryId(id))
	}

	return ret, nil
}

func (db NotesDatabase) saveResultSet(resultSet dbResultSet) (error) {
	//fmt.Println(resultSet)
	str, err := json.Marshal(resultSet)
	if err != nil {
		return fmt.Errorf("Failed to convert resultSet to json <= %s", err)
	}
	//fmt.Println(string(str))

	if err := ioutil.WriteFile(RS_FILE, str, 0644); err != nil {
		return fmt.Errorf("Failed to write contents to resultSet file <= %s", err)
	}

	return nil
}

func (db NotesDatabase) update(entry dbEntry, action string) error {
    switch action {
    case "complete":
        //fmt.Println("Updating", entry)
        data, err := ioutil.ReadFile(entry.Source())
        if err != nil {
            return fmt.Errorf("Failed to open file %s <= %s", entry.Source(), err)
        }

        lines := strings.Split(string(data), "\n")

        entry = entry.update("status", "completed")
        //fmt.Println(entry)
        lines[ entry.LineNum() ] = entry.toString()

        // for number, line := range lines {
	       //  fmt.Println(number, line)
        // }

        new_content := strings.Join(lines, "\n")
        if err := ioutil.WriteFile(entry.Source(), []byte(new_content), 0644); err != nil {
            return fmt.Errorf("Failed to write contents to file %s <= %s", entry.Source(), err)
        }
    }
    return nil
}