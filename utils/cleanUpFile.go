package utils

import (
	"log"
	"os"
)

func CleanUpFile(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		log.Printf("error deleting file!: %s", err)
		return
	} else {
		log.Printf("successfully deleted file: %s\n", filepath)	
	}
}