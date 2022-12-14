package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"plugin"
	"sort"

	mr "github.com/Aadithya-V/mapreduce"
)

// bykey is used for implementing sort.Interface to sort by key values.
type byKey []mr.KeyValue

// Implementing the interface for sorting by key.
func (a byKey) Len() int           { return len(a) }
func (a byKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: main.exe ex.so inputfiles...\n")
		fmt.Fprintf(os.Stderr, "File *.so is the plugin implementing the Map and Reduce functions.\n")
		os.Exit(1)
	}

	mapf, reducef := loadPlugin(os.Args[1])

	// read each input file, pass to Map function,
	// accumulate the intermediate Map output.
	intermediate := []mr.KeyValue{}
	for _, filename := range os.Args[2:] {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("cannot open %v", filename)
		}
		content, err := io.ReadAll(file)
		if err != nil {
			log.Fatalf("cannot read %v", filename)
		}
		file.Close()
		kva := mapf(filename, string(content))
		intermediate = append(intermediate, kva...)
	}

	sort.Sort(byKey(intermediate))

	ofilename := "output"
	ofile, _ := os.Create(ofilename)

	// call Reduce on each distinct key in intermediate[],
	// and print the result to ofilename.
	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := reducef(intermediate[i].Key, values)

		// print result to output file.
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}

	ofile.Close()
}

// loadPlugin loads the application Map and Reduce functions
// from a plugin file, ex- wc.so
func loadPlugin(filename string) (func(string, string) []mr.KeyValue, func(string, []string) string) {
	p, err := plugin.Open(filename)
	if err != nil {
		log.Fatalf("cannot load plugin %v", filename)
	}
	xmapf, err := p.Lookup("Map")
	if err != nil {
		log.Fatalf("cannot find Map in %v", filename)
	}
	mapf := xmapf.(func(string, string) []mr.KeyValue)
	xreducef, err := p.Lookup("Reduce")
	if err != nil {
		log.Fatalf("cannot find Reduce in %v", filename)
	}
	reducef := xreducef.(func(string, []string) string)

	return mapf, reducef
}
