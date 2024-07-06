package goeffective

import "fmt"

// in this lession we will custom the printer of sequence, and custom sort
type Sequence []int

func (s Sequence) String() string {
	return "hello custom string"
}

func Interface() {
	s := Sequence{1, 2, 3, 4, 5}
	fmt.Println(s)
}
