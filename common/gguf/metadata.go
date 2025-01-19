package gguf

// -----------------------------------------------------------------------------
// GGUFMetadataKV: corresponds to struct gguf_metadata_kv_t
// -----------------------------------------------------------------------------

type GGUFMetadataKV struct {
	Key       GGUFString            // gguf_string_t
	ValueType GgufMetadataValueType // enum
	Value     GGUFMetadataValue     // union
}
