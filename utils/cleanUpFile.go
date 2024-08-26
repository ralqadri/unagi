package utils

import (
	"log"
	"os"
)

func CleanUpFile(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		log.Fatalf("error deleting file!: %s", err)
	} else {
		log.Printf("successfully deleted file: %s\n", filepath)	
	}
}