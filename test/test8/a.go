package main

import "github.com/JohnCosta27/go-bridge/test/test8/nested"

type A struct {
	NormalField string
	nested.B
}
