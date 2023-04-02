package utils

import (
	"os"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}
func WriteToFile(outputPath string, text string) {
	os.Remove(outputPath)

	file, err := os.Create(outputPath)
	Panic(err)

	_, err = file.WriteString(text)
	Panic(err)
}
