package gguf

// -----------------------------------------------------------------------------
// GGUFTensorInfo: corresponds to gguf_tensor_info_t
// -----------------------------------------------------------------------------

type GGUFTensorInfo struct {
	Name        GGUFString // must be <= 64 bytes
	NDimensions uint32
	Dimensions  []uint64
	Type        GgmlType // e.g. Q4_0 -> 12, etc.
	Offset      uint64   // offset in the tensor_data region
}
