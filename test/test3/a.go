package main

import "johncosta.tech/struct-to-types/test/test3/nested"

type TestingStruct struct {
	NestedArray    []nested.IAmNested
	SomeOtherField float32
}
