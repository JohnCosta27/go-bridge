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

	outputParse, err := CodeParse(simpleStruct)
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

	outputParse, err := CodeParse(simpleStruct)

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

	outputParse, err := CodeParse(simpleStruct)

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

func TestMultipleSimpleStructs(t *testing.T) {
	simpleStruct := `
package types

type SimpleStruct struct {
	Hello string
}

type AlsoSimpleStruct struct {
	World string
}
  `

	valibotValidator := `
import { object, string } from 'valibot';

const SimpleStruct = object({
  Hello: string(),
});

const AlsoSimpleStruct = object({
  World: string(),
});
`

	outputParse, err := CodeParse(simpleStruct)
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

func TestNestedStruct(t *testing.T) {
	simpleStruct := `
package types

type OtherStruct struct {
	Hello string
}

type MyStruct struct {
	Nested OtherStruct
}
  `

	valibotValidator := `
import { object, string } from 'valibot';

const OtherStruct = object({
  Hello: string(),
});

const MyStruct = object({
  Nested: OtherStruct,
});
`

	outputParse, err := CodeParse(simpleStruct)
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

func TestMultipleNestedStructs(t *testing.T) {
	simpleStruct := `
package types

type A struct {
	A string
  B C
  D E
}

type C struct {
  Hello string
  A E
}

type E struct {
	World bool
  Woooo float64
}
`

	valibotValidator := `
import { object, boolean, number, string } from 'valibot';

const E = object({
  World: boolean(),
  Woooo: number(),
});

const C = object({
  Hello: string(),
  A: E,
});

const A = object({
  A: string(),
  B: C,
  D: E,
});
`

	outputParse, err := CodeParse(simpleStruct)
	t.Log(outputParse)
	t.Log(len(outputParse), len(valibotValidator))

	if err != nil {
		t.Log("Error is not null")
		t.Log(err)
		t.FailNow()
	}

	if outputParse != valibotValidator {
		t.FailNow()
	}
}

func TestEmbeddedStruct(t *testing.T) {
	simpleStruct := `
package types

type B struct {
  A
}

type A struct {
	Hello string
}
`

	valibotValidator := `
import { object, string } from 'valibot';

const B = object({
  Hello: string(),
});

const A = object({
  Hello: string(),
});
`

	outputParse, err := CodeParse(simpleStruct)
	t.Log(outputParse)
	t.Log(len(outputParse), len(valibotValidator))

	if err != nil {
		t.Log("Error is not null")
		t.Log(err)
		t.FailNow()
	}

	if outputParse != valibotValidator {
		t.FailNow()
	}
}

func TestEmbeddedStructComplex(t *testing.T) {
	simpleStruct := `
package types

type A struct {
  Hello float64
  B
  World string
  C
}

type B struct {
	C
  MyField bool
}

type C struct {
  D
  MyNum int
}

type D struct {
  FieldD string
}
`

	valibotValidator := `
import { object, number, string, boolean } from 'valibot';

const A = object({
  Hello: number(),
  FieldD: string(),
  MyNum: number(),
  MyField: boolean(),
  World: string(),
  FieldD: string(),
  MyNum: number(),
});

const B = object({
  FieldD: string(),
  MyNum: number(),
  MyField: boolean(),
});

const C = object({
  FieldD: string(),
  MyNum: number(),
});

const D = object({
  FieldD: string(),
});
`

	outputParse, err := CodeParse(simpleStruct)
	t.Log(outputParse)
	t.Log(len(outputParse), len(valibotValidator))

	if err != nil {
		t.Log("Error is not null")
		t.Log(err)
		t.FailNow()
	}

	if outputParse != valibotValidator {
		t.FailNow()
	}
}

func TestArrayTypes(t *testing.T) {
	simpleStruct := `
package types

type WithArray struct {
	Hello []string
  World []float32
}
  `

	valibotValidator := `
import { object, string, array, number } from 'valibot';

const WithArray = object({
  Hello: array(string()),
  World: array(number()),
});
`

	outputParse, err := CodeParse(simpleStruct)
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

func TestStructArrayTypes(t *testing.T) {
	simpleStruct := `
package types

type ForArray struct {
  SomeField string
}

type WithArray struct {
	Hello []ForArray
}
  `

	valibotValidator := `
import { object, string, array } from 'valibot';

const ForArray = object({
  SomeField: string(),
});

const WithArray = object({
  Hello: array(ForArray),
});
`

	outputParse, err := CodeParse(simpleStruct)
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

func TestMapTypes(t *testing.T) {
	simpleStruct := `
package types

type ForMap struct {
  Hello map[string]string
}
`

	valibotValidator := `
import { object, string, record } from 'valibot';

const ForMap = object({
  Hello: record(string()),
});
`

	outputParse, err := CodeParse(simpleStruct)
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
