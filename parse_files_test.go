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

func TestNestedStructArrayTypes(t *testing.T) {
	valibotString, err := MainParse("./test/test3/a.go", "johncosta.tech/struct-to-types")

	valibotValidator := `
import { object, string, number, array } from 'valibot';

const IAmNested = object({
  Some: string(),
  Field: number(),
  AnotherOneForLuck: array(string()),
});

const TestingStruct = object({
  NestedArray: array(IAmNested),
  SomeOtherField: number(),
});
`

	t.Log(valibotString)

	if err != nil {
		t.Log("Error is not null")
		t.Log(err)
		t.FailNow()
	}

	if valibotString != valibotValidator {
		t.FailNow()
	}
}

func TestFileMapTypes(t *testing.T) {
	valibotString, err := MainParse("./test/test4/a.go", "johncosta.tech/struct-to-types")

	valibotValidator := `
import { object, string, record } from 'valibot';

const IAmNested = object({
  Hello: string(),
});

const WithMap = object({
  Map: record(IAmNested),
});
`

	t.Log(valibotString)

	if err != nil {
		t.Log("Error is not null")
		t.Log(err)
		t.FailNow()
	}

	if valibotString != valibotValidator {
		t.FailNow()
	}
}

func TestDuplicateNames(t *testing.T) {
	t.Skip()
	valibotString, err := MainParse("./test/test5/a.go", "johncosta.tech/struct-to-types")

	valibotValidator := `
import { object, string, record } from 'valibot';

const MyStruct = object({
  Main: string(),
});

const NestedMyStruct = object({
  Nested: string(),
});
`

	t.Log(valibotString)

	if err != nil {
		t.Log("Error is not null")
		t.Log(err)
		t.FailNow()
	}

	if valibotString != valibotValidator {
		t.FailNow()
	}
}
