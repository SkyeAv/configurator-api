package main

import "log"

func checkErr(code uint16, err error) {
	if err != nil {
		log.Fatalf("%d | ERROR | %v", code, err)
	}
}
