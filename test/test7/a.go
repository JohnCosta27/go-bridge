package main

import "github.com/JohnCosta27/go-bridge/test/test7/nested"

type A struct {
	a struct {
		b nested.B
		c struct {
			d nested.B
		}
	}
}
