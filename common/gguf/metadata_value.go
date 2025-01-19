package gguf

import "fmt"

// -----------------------------------------------------------------------------
// GGUFString: corresponds to gguf_string_t
// -----------------------------------------------------------------------------

type GGUFString struct {
	// length in bytes
	Len uint64
	// actual data (not null-terminated in file)
	Data string
}

// GGUFArray corresponds to the array sub-struct in union:
//
//	struct {
//	  gguf_metadata_value_type type;
//	  uint64_t len;
//	  gguf_metadata_value_t array[len];
//	}
type GGUFArray struct {
	ElemType GgufMetadataValueType
	Len      uint64
	// each element is itself a union
	Elements []GGUFMetadataValue
}

// -----------------------------------------------------------------------------
// GGUFMetadataValue: corresponds to union gguf_metadata_value_t
//    We must handle all possible subfields, including nested arrays (recursive).
// -----------------------------------------------------------------------------

type GGUFMetadataValue struct {
	// We store which type it actually is (uint8, int16, string, array, etc.)
	Type GgufMetadataValueType

	// For basic scalar types (uint8, int8, int16, etc.), we store them in these fields.
	// Only one will be used, depending on 'Type'.
	U8  uint8
	I8  int8
	U16 uint16
	I16 int16
	U32 uint32
	I32 int32
	F32 float32
	U64 uint64
	I64 int64
	F64 float64
	B   bool // for GGUF_METADATA_VALUE_TYPE_BOOL

	// For strings:
	StrVal *GGUFString

	// For arrays (recursive):
	ArrVal *GGUFArray
}

// String implements fmt.Stringer for GGUFMetadataValue.
func (mv GGUFMetadataValue) String() string {
	// Use mv.Type to decide how to print
	switch mv.Type {
	case GGUF_METADATA_VALUE_TYPE_UINT8:
		return fmt.Sprintf("%s(%d)", mv.Type, mv.U8)
	case GGUF_METADATA_VALUE_TYPE_INT8:
		return fmt.Sprintf("%s(%d)", mv.Type, mv.I8)
	case GGUF_METADATA_VALUE_TYPE_UINT16:
		return fmt.Sprintf("%s(%d)", mv.Type, mv.U16)
	case GGUF_METADATA_VALUE_TYPE_INT16:
		return fmt.Sprintf("%s(%d)", mv.Type, mv.I16)
	case GGUF_METADATA_VALUE_TYPE_UINT32:
		return fmt.Sprintf("%s(%d)", mv.Type, mv.U32)
	case GGUF_METADATA_VALUE_TYPE_INT32:
		return fmt.Sprintf("%s(%d)", mv.Type, mv.I32)
	case GGUF_METADATA_VALUE_TYPE_FLOAT32:
		return fmt.Sprintf("%s(%g)", mv.Type, mv.F32)
	case GGUF_METADATA_VALUE_TYPE_BOOL:
		return fmt.Sprintf("%s(%t)", mv.Type, mv.B)
	case GGUF_METADATA_VALUE_TYPE_UINT64:
		return fmt.Sprintf("%s(%d)", mv.Type, mv.U64)
	case GGUF_METADATA_VALUE_TYPE_INT64:
		return fmt.Sprintf("%s(%d)", mv.Type, mv.I64)
	case GGUF_METADATA_VALUE_TYPE_FLOAT64:
		return fmt.Sprintf("%s(%g)", mv.Type, mv.F64)
	case GGUF_METADATA_VALUE_TYPE_STRING:
		if mv.StrVal != nil {
			return fmt.Sprintf("STRING(len=%d, data=%q)", mv.StrVal.Len, mv.StrVal.Data)
		}
		return "STRING(nil)"
	case GGUF_METADATA_VALUE_TYPE_ARRAY:
		// For an array, we might show length + element type, or go deeper
		if mv.ArrVal != nil {
			return fmt.Sprintf("ARRAY(elem=%s, len=%d)", mv.ArrVal.ElemType, mv.ArrVal.Len)
			// or if you want to deeply print each element, do a loop
		}
		return "ARRAY(nil)"

	default:
		return fmt.Sprintf("UNKNOWN_VALUE_TYPE(%d)", mv.Type)
	}
}
