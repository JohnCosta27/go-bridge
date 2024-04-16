package main

import "testing"

func TestSamePackage(t *testing.T) {
	valibotString, err := MainParse("./test/test1/a.go", "johncosta.tech/struct-to-types")

	valibotValidator := `
import { object, string } from 'valibot';

const B = object({
  Hello: string(),
});

const A = object({
  Hello: B,
});
`

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(valibotString)

	if valibotString != valibotValidator {
		t.FailNow()
	}
}

func TestNestedDeps(t *testing.T) {
	valibotString, err := MainParse("./test/test2/a.go", "johncosta.tech/struct-to-types")

	valibotValidator := `
import { object, number, string } from 'valibot';

const D = object({
  VeryNested: number(),
});

const NestedStruct = object({
  Hello: string(),
  Bruh: D,
});

const MainStruct = object({
  World: NestedStruct,
});
`
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(valibotString)

	if valibotString != valibotValidator {
		t.FailNow()
	}
}
