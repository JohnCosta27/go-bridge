package main

import "testing"

func TestBasicParseSimpleStruct1(t *testing.T) {
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

func TestBasicParseSimpleStruct2(t *testing.T) {
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

func TestBasicParseSimpleStruct3(t *testing.T) {
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

func TestBasicParseMultipleSimpleStructs(t *testing.T) {
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

func TestBasicParseNestedStruct(t *testing.T) {
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

func TestBasicParseMultipleNestedStructs(t *testing.T) {
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

func TestBasicParseEmbeddedStruct(t *testing.T) {
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

func TestBasicParseEmbeddedStructComplex(t *testing.T) {
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

func TestBasicParseArrayTypes(t *testing.T) {
	simpleStruct := `
package types

type WithArray struct {
	Hello []string
  World []float32
}
  `

	valibotValidator := `
import { object, array, string, number } from 'valibot';

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

func TestBasicParseStructArrayTypes(t *testing.T) {
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

func TestBasicParseMapTypesSimple(t *testing.T) {
	simpleStruct := `
package types

type ForMap struct {
  Hello map[string]string
}
`

	valibotValidator := `
import { object, record, string } from 'valibot';

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

func TestBasicParseMapTypesToStruct(t *testing.T) {
	simpleStruct := `
package types

type A struct {
  Field int
}

type ForMap struct {
  Hello map[string]A
}
`

	valibotValidator := `
import { object, number, record } from 'valibot';

const A = object({
  Field: number(),
});

const ForMap = object({
  Hello: record(A),
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

func TestBasicParseMapTypesArraySimple(t *testing.T) {
	simpleStruct := `
package types

type ForMap struct {
  Hello map[string][]string
}
`

	valibotValidator := `
import { object, record, array, string } from 'valibot';

const ForMap = object({
  Hello: record(array(string())),
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

func TestBasicParseChaos(t *testing.T) {
	simpleStruct := `
package types

type ForMap struct {
  Hello map[string][]string
}

type BigType struct {
  A []map[string]ForMap
  B map[string]map[string][][][]map[string]ForMap
}
`

	valibotValidator := `
import { object, record, array, string } from 'valibot';

const ForMap = object({
  Hello: record(array(string())),
});

const BigType = object({
  A: array(record(ForMap)),
  B: record(record(array(array(array(record(ForMap)))))),
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

func TestBasicParseAnonStruct1(t *testing.T) {
	simpleStruct := `
package types

type A struct {
  Hello struct {
    World string
  }
}
`

	valibotValidator := `
import { object, string } from 'valibot';

const A = object({
  Hello: object({
    World: string(),
  }),
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

func TestBasicParseAnonStruct2(t *testing.T) {
	simpleStruct := `
package types

type A struct {
  Hello struct {
    World string
    A map[string][]int
    AnotherStruct struct{
      B string
    }
    B []struct{
      C string
    }
  }
}
`

	valibotValidator := `
import { object, string, record, array, number } from 'valibot';

const A = object({
  Hello: object({
    World: string(),
    A: record(array(number())),
    AnotherStruct: object({
      B: string(),
    }),
    B: array(object({
        C: string(),
      })),
  }),
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

func TestPointers(t *testing.T) {
	t.Run("Simple pointer to string", func(t *testing.T) {
		simpleStruct := `
package types

type A struct {
  Pointer *string
}
`

		valibotValidator := `
import { object, string } from 'valibot';

const A = object({
  Pointer: string(),
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
	})

	t.Run("Pointers to structs", func(t *testing.T) {
		simpleStruct := `
package types

type A struct {
  Hello string
}

type B struct {
  APointer *A
  AArrayPointer []*A
  AMapPointer map[string]*A
}
`

		valibotValidator := `
import { object, string, array, record } from 'valibot';

const A = object({
  Hello: string(),
});

const B = object({
  APointer: A,
  AArrayPointer: array(A),
  AMapPointer: record(A),
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
	})
}
