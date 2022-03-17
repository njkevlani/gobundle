package main

import (
	"fmt"

	"github.com/njkevlani/gobundle/test_files/test_project0/ds"
)

func main() {
	t := ds.NewTrie()
	t.Add("lorem")
	fmt.Println("trie has 'lorem' -> ", t.Has("lorem"))
	fmt.Println("trie has 'ipsum' -> ", t.Has("ipsum"))
}
