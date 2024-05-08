package main

import "testing"

func TestBackendSimpleType(t *testing.T) {
	validators := make(map[string]uint)
	nameMap := make(map[string]string)

	var counter uint = 0

	basicField := BasicStructField{name: "Name", Type: "string"}

	output, err := getStructFieldType(validators, nameMap, &counter, basicField, 0)
	expected := "string()"

	t.Log(output)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(validators) != 1 {
		t.Errorf("Validators should have length 1, got %+v\n", validators)
		t.FailNow()
	}

	_, exists := validators["string"]
	if !exists {
		t.Error("Should have gotten string in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	if output != expected {
		t.Error("Should have gotten expected value.\n")
		t.FailNow()
	}
}

func TestBackendStructType(t *testing.T) {
	validators := make(map[string]uint)
	var counter uint = 0

	nameMap := make(map[string]string)
	nameMap["AnotherStruct"] = "AnotherStruct"

	basicField := BasicStructField{name: "Name", Type: "AnotherStruct"}

	output, err := getStructFieldType(validators, nameMap, &counter, basicField, 0)
	expected := "AnotherStruct"

	t.Log(output)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(validators) != 0 {
		t.Errorf("Validators should have length 0, got %+v\n", validators)
		t.FailNow()
	}

	if output != expected {
		t.Error("Should have gotten expected value.\n")
		t.FailNow()
	}
}

func TestBackendArrayType(t *testing.T) {
	validators := make(map[string]uint)
	var counter uint = 0

	arrayField := ArrayStructField{name: "Name", Type: BasicStructField{name: "Name", Type: "int64"}}
	nameMap := make(map[string]string)

	output, err := getStructFieldType(validators, nameMap, &counter, arrayField, 0)
	expected := "array(number())"

	t.Log(output)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(validators) != 2 {
		t.Errorf("Validators should have length 2, got %+v\n", validators)
		t.FailNow()
	}

	_, exists := validators["array"]
	if !exists {
		t.Error("Should have gotten array in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	_, exists = validators["number"]
	if !exists {
		t.Error("Should have gotten number in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	if output != expected {
		t.Error("Should have gotten expected value.\n")
		t.FailNow()
	}
}

func TestBackendArrayStructType(t *testing.T) {
	validators := make(map[string]uint)
	var counter uint = 0

	nameMap := make(map[string]string)
	nameMap["SomeStruct"] = "SomeStruct"

	arrayField := ArrayStructField{name: "Name", Type: BasicStructField{name: "Name", Type: "SomeStruct"}}

	output, err := getStructFieldType(validators, nameMap, &counter, arrayField, 0)
	expected := "array(SomeStruct)"

	t.Log(output)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(validators) != 1 {
		t.Errorf("Validators should have length 1, got %+v\n", validators)
		t.FailNow()
	}

	_, exists := validators["array"]
	if !exists {
		t.Error("Should have gotten array in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	if output != expected {
		t.Error("Should have gotten expected value.\n")
		t.FailNow()
	}
}

func TestBackendArrayArrayType(t *testing.T) {
	validators := make(map[string]uint)
	var counter uint = 0

	arrayField := ArrayStructField{name: "Name", Type: ArrayStructField{name: "Name", Type: BasicStructField{name: "Name", Type: "bool"}}}

	nameMap := make(map[string]string)
	nameMap["Name"] = "Name"

	output, err := getStructFieldType(validators, nameMap, &counter, arrayField, 0)
	expected := "array(array(boolean()))"

	t.Log(output)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(validators) != 2 {
		t.Errorf("Validators should have length 2, got %+v\n", validators)
		t.FailNow()
	}

	_, exists := validators["array"]
	if !exists {
		t.Error("Should have gotten array in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	_, exists = validators["boolean"]
	if !exists {
		t.Error("Should have gotten boolean in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	if output != expected {
		t.Error("Should have gotten expected value.\n")
		t.FailNow()
	}
}

func TestBackendMapType(t *testing.T) {
	validators := make(map[string]uint)
	var counter uint = 0

	arrayField := MapStructField{name: "Name", KeyType: "string", Value: BasicStructField{name: "Name", Type: "uint"}}

	nameMap := make(map[string]string)
	nameMap["Name"] = "Name"

	output, err := getStructFieldType(validators, nameMap, &counter, arrayField, 0)
	expected := "record(number())"

	t.Log(output)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(validators) != 2 {
		t.Errorf("Validators should have length 2, got %+v\n", validators)
		t.FailNow()
	}

	_, exists := validators["record"]
	if !exists {
		t.Error("Should have gotten record in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	_, exists = validators["number"]
	if !exists {
		t.Error("Should have gotten number in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	if output != expected {
		t.Error("Should have gotten expected value.\n")
		t.FailNow()
	}
}

func TestBackendMapArrayStructType(t *testing.T) {
	validators := make(map[string]uint)
	var counter uint = 0

	arrayField := MapStructField{name: "Name", KeyType: "string", Value: ArrayStructField{name: "Name", Type: BasicStructField{name: "Name", Type: "string"}}}

	nameMap := make(map[string]string)
	nameMap["Name"] = "Name"

	output, err := getStructFieldType(validators, nameMap, &counter, arrayField, 0)
	expected := "record(array(string()))"

	t.Log(output)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(validators) != 3 {
		t.Errorf("Validators should have length 3, got %+v\n", validators)
		t.FailNow()
	}

	_, exists := validators["record"]
	if !exists {
		t.Error("Should have gotten record in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	_, exists = validators["array"]
	if !exists {
		t.Error("Should have gotten array in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	_, exists = validators["string"]
	if !exists {
		t.Error("Should have gotten string in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	if output != expected {
		t.Error("Should have gotten expected value.\n")
		t.FailNow()
	}
}

func TestBackendChaos(t *testing.T) {
	validators := make(map[string]uint)
	var counter uint = 0

	arrayField := ArrayStructField{
		name: "Name",
		Type: MapStructField{
			name:    "Name",
			KeyType: "string",
			Value: ArrayStructField{
				name: "Name", Type: MapStructField{
					name:    "Name",
					KeyType: "string",
					Value: MapStructField{
						name:    "Name",
						KeyType: "string",
						Value: ArrayStructField{
							name: "Name",
							Type: BasicStructField{
								name: "Name",
								Type: "string",
							},
						},
					},
				},
			},
		},
	}

	nameMap := make(map[string]string)
	nameMap["Name"] = "Name"

	output, err := getStructFieldType(validators, nameMap, &counter, arrayField, 0)
	expected := "array(record(array(record(record(array(string()))))))"

	t.Log(output)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(validators) != 3 {
		t.Errorf("Validators should have length 3, got %+v\n", validators)
		t.FailNow()
	}

	_, exists := validators["record"]
	if !exists {
		t.Error("Should have gotten record in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	_, exists = validators["array"]
	if !exists {
		t.Error("Should have gotten array in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	_, exists = validators["string"]
	if !exists {
		t.Error("Should have gotten string in validators\n")
		for k := range validators {
			t.Error(k + "\n")
		}
		t.FailNow()
	}

	if output != expected {
		t.Error("Should have gotten expected value.\n")
		t.FailNow()
	}
}
