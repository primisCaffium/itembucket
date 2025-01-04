package common

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func SuppressStackTraceOnPanic() {
	if r := recover(); r != nil {
		fmt.Println(r)
		os.Exit(1)
	}
}

func Panic(err interface{}) {
	if err != nil {
		panic(err)
	}
}

func PInt64(v int64) *int64 {
	return &v
}

func PStr(v string) *string {
	return &v
}

func PTime(v time.Time) *time.Time {
	return &v
}

func Unmarshal(bytes []byte, object interface{}) {
	err := json.Unmarshal(bytes, object)
	if err != nil {
		panic(err)
	}
}

func Marshal(object interface{}) []byte {
	b, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	return b
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func WriteToFile(outputPath string, text string) {
	err := os.Remove(outputPath)
	Panic(err)

	file, err := os.Create(outputPath)
	Panic(err)

	_, err = file.WriteString(text)
	Panic(err)
}
