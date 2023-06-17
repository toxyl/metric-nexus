package main

import "os"

func fileExists(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	_, err = file.Stat()
	return err == nil
}
