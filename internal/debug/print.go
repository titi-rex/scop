package debug

import "log"

func Dprint(a, b any) {
	log.Printf("A:\n%v\n", a)
	log.Printf("B:\n%v\n", b)
}
