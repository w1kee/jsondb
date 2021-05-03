package jsondb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"reflect"
)

type (
	Database struct {
		fs fs.FS
	}
)

// Opens a database located at sysfp
// This works on the real os filesystem
func Open(sysfp string) (*Database, error) {
	f, err := os.Open(sysfp)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file stat: %w", err)
	}
	if !s.IsDir() {
		return nil, fmt.Errorf("file is not dir")
	}
	return &Database{os.DirFS(sysfp)}, nil
}

// WARNING: not implemented
// OpenFS opens a database on fsys
func OpenFS(fsys fs.FS) (*Database, error) {
	return nil, fmt.Errorf("not implemented")
}

// Get gets an item from fps and stores it to v
//
// v must be a pointer or Get panics
func (db *Database) Get(v interface{}, fps ...string) error {
	if reflect.ValueOf(v).Type().Kind() != reflect.Ptr {
		return fmt.Errorf("v is not a pointer")
	}

	f, err := db.fs.Open(path.Join(fps...) + ".json")
	if err != nil {
		return err
	}
	defer f.Close()

	obj, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(obj, v)
	if err != nil {
		return err
	}
	return nil
}

// Set stores v into fps
func (db *Database) Set(v interface{}, fps ...string) error {
	// TODO: check if v is of a type that can be stored in json file (struct{})

	f, err := db.fs.Open(path.Join(fps...))
	if err != nil {
		return err
	}
	defer f.Close()

	fw, ok := f.(io.Writer)
	if !ok {
		return fmt.Errorf("file doesn't implement io.Writer")
	}
	obj, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = fw.Write(obj)
	if err != nil {
		return err
	}

	return nil
}
