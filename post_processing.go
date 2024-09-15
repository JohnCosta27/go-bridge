package main

import "slices"

// Naive solution
//
// 1. Search through every single field of struct.
// 2. Recursively search through fields that can contain more fields.
// 3. Replace the ones that are embedded.
//
// A better approach would involve keeping some map of where the embedded
// structs are, instead of having to search for them.
func rebuildStruct(structs []Struct, processingStruct Struct) Struct {
	processingStruct.Fields = rebuildStructFields(structs, processingStruct.Fields)

	return processingStruct
}

func rebuildStructFields(structs []Struct, fields []StructField) []StructField {
	processedFields := make([]StructField, 0)

	for _, field := range fields {
		switch t := field.(type) {
		case UnknownStructField:
			processedFields = append(processedFields, t)
		case BasicStructField:
			if t.name != EMBEDDED_DEP {
				processedFields = append(processedFields, t)
				continue
			}

			embeddedStructIndex := slices.IndexFunc(structs, func(s Struct) bool {
				return s.Name == t.Type
			})

			processedFields = append(processedFields, structs[embeddedStructIndex].Fields...)
		default:
			processedFields = append(processedFields, recReplaceEmbeddedStruct(structs, t))
		}
	}

	return processedFields
}

func recReplaceEmbeddedStruct(structs []Struct, field StructField) StructField {
	switch t := field.(type) {
	case AnonStructField:
		return AnonStructField{
			name: t.name,

			Fields: rebuildStructFields(structs, t.Fields),
		}
	case MapStructField:
		return recReplaceEmbeddedStruct(structs, t.Value)
	case ArrayStructField:
		return recReplaceEmbeddedStruct(structs, t.Type)
	default:
		panic("unreachable")
	}
}
