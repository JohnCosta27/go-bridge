package main

import "testing"

func TestSamePackage(t *testing.T) {
	valibotString, err := Parse("./test/test1/a.go", "johncosta.tech/struct-to-types")

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
	t.Skip()

	valibotString, err := Parse("./test/test2/a.go", "johncosta.tech/struct-to-types")

	valibotValidator := `
import { object, string } from 'valibot';

const NestedStruct = object({
  Hello: string(),
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

func TestDeepNestedDeps(t *testing.T) {
	t.Skip()

	valibotString, err := Parse("./test/test2/a.go", "johncosta.tech/struct-to-types")

	valibotValidator := `
import { object, string } from 'valibot';

const NestedStruct = object({
  Hello: string(),
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
