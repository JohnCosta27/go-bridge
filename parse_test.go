package main

import "testing"

func TestSimpleStruct1(t *testing.T) {
	simpleStruct := `
package types

type SimpleStruct struct {
	Hello string
}
  `

	valibotValidator := `
import { object, string } from 'valibot';

const SimpleStruct = object({
  Hello: string(),
});
`

	outputParse, err := Parse(simpleStruct)
	t.Log(outputParse)

	if err != nil {
		t.Log("Error is not null")
		t.Log(err)
		t.FailNow()
	}

	if outputParse != valibotValidator {
		t.FailNow()
	}
}

func TestSimpleStruct2(t *testing.T) {
	simpleStruct := `
package types

type NotAsSimple struct {
	Hello int
  World int64
  A string
  B bool
}
  `

	valibotValidator := `
import { object, number, string, boolean } from 'valibot';

const NotAsSimple = object({
  Hello: number(),
  World: number(),
  A: string(),
  B: boolean(),
});
`

	outputParse, err := Parse(simpleStruct)

	t.Log(outputParse)

	if err != nil {
		t.Log("Error is not null")
		t.Log(err)
		t.FailNow()
	}

	if outputParse != valibotValidator {
		t.FailNow()
	}
}

func TestSimpleStruct3(t *testing.T) {
	simpleStruct := `
package types

type SimpleButComplex struct {
	A int
  B float32
  C bool
  D string
  E string
  F int8
  G int16
  H int32
  I int64
  J float64
}
`

	valibotValidator := `
import { object, number, boolean, string } from 'valibot';

const SimpleButComplex = object({
  A: number(),
  B: number(),
  C: boolean(),
  D: string(),
  E: string(),
  F: number(),
  G: number(),
  H: number(),
  I: number(),
  J: number(),
});
`

	outputParse, err := Parse(simpleStruct)

	t.Log(outputParse)

	if err != nil {
		t.Log("Error is not null")
		t.Log(err)
		t.FailNow()
	}

	if outputParse != valibotValidator {
		t.FailNow()
	}
}
