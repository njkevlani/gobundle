package main

import "fmt"

func main() {
	t := NewTrie()
	t.Add("lorem")
	fmt.Println("trie has 'lorem' -> ", t.Has("lorem"))
	fmt.Println("trie has 'ipsum' -> ", t.Has("ipsum"))
}

type Trie struct{ Next map[rune]Trie }

func (t Trie) Add(word string) {
	cur := t
	for _, ch := range word {
		if _, ok := cur.Next[ch]; !ok {
			cur.Next[ch] = Trie{Next: make(map[rune]Trie)}
		}
		cur = cur.Next[ch]
	}
	var nt Trie
	cur.Next['$'] = nt
}
func (t Trie) Has(word string) bool {
	cur := t
	var ok bool
	for _, ch := range word {
		if cur, ok = cur.Next[ch]; !ok {
			return false
		}
	}
	_, endsHere := cur.Next['$']
	return endsHere
}
func NewTrie() Trie {
	return Trie{Next: make(map[rune]Trie)}
}
