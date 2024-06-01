package main

import "github.com/JohnCosta27/go-bridge/test/test7/nested"

type Simple struct {
	Hello string
}

type A struct {
	nested.B
	Simple

	Filed struct {
		nested.B
		Simple
	}
}
