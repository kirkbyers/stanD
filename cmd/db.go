package cmd

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	eventsBucket = []byte("events")
)

// State is the state object for stanD
type State struct {
	dbPath string
}

// Init should be called before stat is used
func (s *State) Init(fp string) (err error) {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	if fp == "" {
		s.dbPath = home + "/.stand.db"
	} else {
		s.dbPath = fp
	}

	return s.update(func(tx *bolt.Tx) (err error) {
		// Create Buckets
		_, err = tx.CreateBucketIfNotExists(eventsBucket)
		if err != nil {
			return err
		}
		return nil
	})
}

// update is a helper function to read+write to bolt
func (s *State) update(ufn func(*bolt.Tx) error) (err error) {
	db, err := bolt.Open(s.dbPath, 0600, nil)
	defer db.Close()
	if err != nil {
		return err
	}
	err = db.Update(ufn)
	if err != nil {
		return err
	}
	return
}

// A Task is a single entry
type Task struct {
	Timestamp time.Time `json:"timestamp"`
	Body      string    `json:"body"`
}

// AddTask adds a task to the saved state for the current day (local time)
func (s *State) AddTask(body string) (err error) {
	newTask := Task{
		Timestamp: time.Now(),
		Body:      body,
	}
	newTaskB, err := json.Marshal(newTask)
	if err != nil {
		return err
	}
	return s.update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket(eventsBucket)
		return b.Put([]byte(newTask.Timestamp.Format(("Jan 2 2006"))), newTaskB)
	})
}

// GetDay gets task that were saved at the time stamp given
func GetDay() {}
