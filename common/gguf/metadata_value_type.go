package gguf

import "fmt"

// gguf_metadata_value_type is a 32-bit enum
type GgufMetadataValueType uint32

const (
	GGUF_METADATA_VALUE_TYPE_UINT8   GgufMetadataValueType = 0
	GGUF_METADATA_VALUE_TYPE_INT8    GgufMetadataValueType = 1
	GGUF_METADATA_VALUE_TYPE_UINT16  GgufMetadataValueType = 2
	GGUF_METADATA_VALUE_TYPE_INT16   GgufMetadataValueType = 3
	GGUF_METADATA_VALUE_TYPE_UINT32  GgufMetadataValueType = 4
	GGUF_METADATA_VALUE_TYPE_INT32   GgufMetadataValueType = 5
	GGUF_METADATA_VALUE_TYPE_FLOAT32 GgufMetadataValueType = 6
	GGUF_METADATA_VALUE_TYPE_BOOL    GgufMetadataValueType = 7
	GGUF_METADATA_VALUE_TYPE_STRING  GgufMetadataValueType = 8
	GGUF_METADATA_VALUE_TYPE_ARRAY   GgufMetadataValueType = 9
	GGUF_METADATA_VALUE_TYPE_UINT64  GgufMetadataValueType = 10
	GGUF_METADATA_VALUE_TYPE_INT64   GgufMetadataValueType = 11
	GGUF_METADATA_VALUE_TYPE_FLOAT64 GgufMetadataValueType = 12
)

// String implements the fmt.Stringer interface for GgufMetadataValueType.
func (t GgufMetadataValueType) String() string {
	switch t {
	case GGUF_METADATA_VALUE_TYPE_UINT8:
		return "UINT8"
	case GGUF_METADATA_VALUE_TYPE_INT8:
		return "INT8"
	case GGUF_METADATA_VALUE_TYPE_UINT16:
		return "UINT16"
	case GGUF_METADATA_VALUE_TYPE_INT16:
		return "INT16"
	case GGUF_METADATA_VALUE_TYPE_UINT32:
		return "UINT32"
	case GGUF_METADATA_VALUE_TYPE_INT32:
		return "INT32"
	case GGUF_METADATA_VALUE_TYPE_FLOAT32:
		return "FLOAT32"
	case GGUF_METADATA_VALUE_TYPE_BOOL:
		return "BOOL"
	case GGUF_METADATA_VALUE_TYPE_STRING:
		return "STRING"
	case GGUF_METADATA_VALUE_TYPE_ARRAY:
		return "ARRAY"
	case GGUF_METADATA_VALUE_TYPE_UINT64:
		return "UINT64"
	case GGUF_METADATA_VALUE_TYPE_INT64:
		return "INT64"
	case GGUF_METADATA_VALUE_TYPE_FLOAT64:
		return "FLOAT64"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", uint32(t))
	}
}
