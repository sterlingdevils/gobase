package errorutils

import "log"

func CheckFatalError[T any](t T, e error) T {
	if e != nil {
		log.Fatalln(e)
	}
	return t
}

func CheckError[T any](t T, e error) T {
	if e != nil {
		log.Println(e)
	}
	return t
}
