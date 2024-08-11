package main

import "johncosta.tech/go-bridge/test/test7/nested"

type A struct {
	a struct {
		b nested.B
		c struct {
			d nested.B
		}
	}
}
