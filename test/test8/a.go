package main

import "github.com/JohnCosta27/go-bridge/test/test8/nested"

type A struct {
	NormalField string
	nested.B

	MyStruct struct {
		nested.B

		MyNestedStruct map[string]struct {
			nested.B
		}

		MyAnonArrayStruct []struct {
			nested.B
		}
	}
}
