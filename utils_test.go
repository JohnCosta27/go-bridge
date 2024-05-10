package main

import "testing"

func TestSplit(t *testing.T) {
	t.Run("Easy test", func(t *testing.T) {

		name := getNameWithPackage("nested/morenested-StructName", 0)
		t.Log(name)
		if name != "StructName" {
			t.Fail()
		}
	})

	t.Run("Single package name test", func(t *testing.T) {

		name := getNameWithPackage("morenested-StructName", 0)
		t.Log(name)
		if name != "StructName" {
			t.Fail()
		}
	})
}
