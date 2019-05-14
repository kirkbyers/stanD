package main

import "github.com/boltdb/bolt"

// State is the state object for stanD
type State struct {
	DB *bolt.DB
}

// AddTask adds a task to the saved state for the current day (local time)
func AddTask() {}

// GetDay gets task that were saved at the time stamp given
func GetDay() {}
