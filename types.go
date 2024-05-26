package main

type StructField interface {
	Name() string
}

type BasicStructField struct {
	/* Type can be golang type or a golang struct type */
	Type string

	name string
}

type UnknownStructField struct {
	FullType string

	name string
}

type ArrayStructField struct {
	Type StructField

	name string
}

type MapStructField struct {
	KeyType string
	Value   StructField

	name string
}

type AnonStructField struct {
	Fields []StructField

	name string
}

func (s BasicStructField) Name() string {
	return s.name
}

func (s UnknownStructField) Name() string {
	return s.name
}

func (s ArrayStructField) Name() string {
	return s.name
}

func (s MapStructField) Name() string {
	return s.name
}

func (s AnonStructField) Name() string {
	return s.name
}

type Struct struct {
	Order uint

	Name        string
	PackagePath string

	Fields []StructField
}

type StructList []Struct
