package gobase

import "log"

func CheckFatalError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func CheckError(e error) {
	if e != nil {
		log.Println(e)
	}
}

func CheckFatalError2[T any](t T, e error) T {
	CheckFatalError(e)
	return t
}

func CheckError2[T any](t T, e error) T {
	CheckError(e)
	return t
}

func CheckFatalError3[T, U any](t T, u U, e error) (T, U) {
	CheckFatalError(e)
	return t, u
}

func CheckError3[T, U any](t T, u U, e error) (T, U) {
	CheckError(e)
	return t, u
}
