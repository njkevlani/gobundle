package algo

var isSeen map[int]bool

func init() {
	isSeen = make(map[int]bool)
}

func IsSeen(num int) bool {
	ret := isSeen[num]
	isSeen[num] = true

	return ret
}
