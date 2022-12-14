package main

//
// plugin-to compile run- go build -buildmode=plugin wordcount.go
//

import (
	"strconv"
	"strings"
	"unicode"
)

// The Map function is called once for each file of input. The first
// The return value is a slice of structs of key/value pairs.
func Map(filename string, contents string) []KeyValue {
	// tokenize contents into an array of words.
	ff := func(r rune) bool { return !unicode.IsLetter(r) }
	words := strings.FieldsFunc(contents, ff)

	kva := []KeyValue{}
	for _, w := range words {
		kv := KeyValue{w, "1"}
		kva = append(kva, kv)
	}
	return kva
}

// The Reduce function is called once for each key generated by the
// map tasks, with a list of all the values created for that key by
// any map task.
func Reduce(key string, values []string) string {
	// return the number of occurrences of this "key".
	return strconv.Itoa(len(values))
}
