// TODO: this
// reference: https://github.com/Corsace/corsace/blob/master/Server/utils/link.ts#L4

package utils

import (
	"fmt"
	"strings"
)

func SanitizeFileName(filename string) string {
	if strings.Contains(filename, "?") {
		filename = strings.Split(filename, "?")[0]
	}

	if strings.Contains(filename, "#") {
		filename = strings.Split(filename, "#")[0]
	}

	fmt.Println("cleaned filename: ", filename)
	
	return filename
}