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
	Body string `json:"body"`
}

// Day is a collection of tasks
type Day struct {
	Timestamp time.Time `json:"timestamp"`
	Tasks     []Task    `json:"tasks"`
}

// AddTask adds a task to the saved state for the current day (local time)
func (s *State) AddTask(body string) (err error) {
	now := time.Now()
	newTask := Task{
		Body: body,
	}
	var day Day
	if err != nil {
		return err
	}
	return s.update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket(eventsBucket)
		dayK := []byte(now.Format(("Jan 2 2006")))
		dayV := b.Get(dayK)
		err = json.Unmarshal(dayV, &day)
		if err != nil {
			day = Day{
				Timestamp: now,
				Tasks:     make([]Task, 0),
			}
		}
		day.Tasks = append(day.Tasks, newTask)
		dayUpdate, err := json.Marshal(day)
		if err != nil {
			return err
		}
		return b.Put([]byte(now.Format("Jan 2 2006")), dayUpdate)
	})
}

// GetDay gets task that were saved at the time stamp given
func (s *State) GetDay(offset int) (res Day, err error) {
	now := time.Now()
	targetDate := now.Add(time.Minute * 60 * 24 * time.Duration(offset) * -1)
	var dayV []byte
	err = s.update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket(eventsBucket)
		dayK := []byte(targetDate.Format("Jan 2 2006"))
		dayV = b.Get(dayK)
		if err = json.Unmarshal(dayV, &res); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return res, err
	}

	return res, err
}
