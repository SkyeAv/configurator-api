package main

import "log"

func CheckErr(code uint16, err error) {
	if err != nil {
		log.Fatalf("%d | ERROR | %v", code, err)
	}
}
