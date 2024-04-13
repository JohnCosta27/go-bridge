package main

import "testing"

func TestGoFile(t *testing.T) {
	valibotString, err := ParseV2("./test/test1/a.go")

	valibotValidator := `
import { object, string } from 'valibot';

const A = object({
  Hello: string(),
});

const B = object({
  Hello: string(),
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

	valibotString, err := ParseV2("./test/test2/a.go")

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

	valibotString, err := ParseV2("./test/test2/a.go")

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
