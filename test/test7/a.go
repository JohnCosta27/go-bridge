package main

import "github.com/JohnCosta27/go-bridge/test/test7/nested"

type Simple struct {
	A string
}

type A struct {
	nested.B
}
