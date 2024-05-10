package main

import "johncosta.tech/go-bridge/test/test3/nested"

type TestingStruct struct {
	NestedArray    []nested.IAmNested
	SomeOtherField float32
}
