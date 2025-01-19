package gguf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
)

// Magic, Version, etc.
const (
	GGUF_MAGIC   = 0x46554747 // 'GGUF' in little-endian
	GGUF_VERSION = 3          // per the updated spec
)

// -----------------------------------------------------------------------------
// GGUFFile: top-level structure for parsing
// -----------------------------------------------------------------------------

type GGUFFile struct {
	Magic           uint32
	Version         uint32
	TensorCount     uint64
	MetadataKVCount uint64

	// array of metadata key-value
	Metadata []GGUFMetadataKV

	// array of tensor descriptors
	Tensors []GGUFTensorInfo
}

// Parse parses the entire GGUF header and metadata (and optionally tensor info),
// but does NOT read actual tensor data. That part is up to you afterwards.
func Parse(r io.Reader) (*GGUFFile, error) {
	var file GGUFFile

	// 1) read magic
	if err := binary.Read(r, binary.LittleEndian, &file.Magic); err != nil {
		return nil, fmt.Errorf("read magic: %w", err)
	}
	if file.Magic != GGUF_MAGIC {
		return nil, fmt.Errorf("not a GGUF file, magic=0x%X", file.Magic)
	}

	// 2) version
	if err := binary.Read(r, binary.LittleEndian, &file.Version); err != nil {
		return nil, fmt.Errorf("read version: %w", err)
	}
	if file.Version != GGUF_VERSION {
		return nil, fmt.Errorf("unsupported GGUF version: %d (expected %d)", file.Version, GGUF_VERSION)
	}

	// 3) tensor_count (uint64)
	if err := binary.Read(r, binary.LittleEndian, &file.TensorCount); err != nil {
		return nil, fmt.Errorf("read tensor_count: %w", err)
	}

	// 4) metadata_kv_count (uint64)
	if err := binary.Read(r, binary.LittleEndian, &file.MetadataKVCount); err != nil {
		return nil, fmt.Errorf("read metadata_kv_count: %w", err)
	}

	magicBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(magicBytes, file.Magic)
	log.Printf("Magic: %s(%x)\n", magicBytes, file.Magic)
	log.Printf("Version: %v\n", file.Version)
	log.Printf("metadataCount: %v\n", file.MetadataKVCount)
	log.Printf("tensorCount: %v\n", file.TensorCount)

	// 5) read all metadata_kv
	file.Metadata = make([]GGUFMetadataKV, 0, file.MetadataKVCount)
	for i := uint64(0); i < file.MetadataKVCount; i++ {
		kv, err := parseMetadataKV(r)
		if err != nil {
			return nil, fmt.Errorf("metadata_kv[%d]: %w", i, err)
		}
		file.Metadata = append(file.Metadata, kv)
	}

	// 6) read all tensor infos
	file.Tensors = make([]GGUFTensorInfo, 0, file.TensorCount)
	for i := uint64(0); i < file.TensorCount; i++ {
		ti, err := parseTensorInfo(r)
		if err != nil {
			return nil, fmt.Errorf("tensor_info[%d]: %w", i, err)
		}
		file.Tensors = append(file.Tensors, ti)
	}

	return &file, nil
}

// parseMetadataKV reads one gguf_metadata_kv_t
//
//	struct gguf_metadata_kv_t {
//	  gguf_string_t key;
//	  gguf_metadata_value_type value_type;
//	  gguf_metadata_value_t value;
//	};
func parseMetadataKV(r io.Reader) (GGUFMetadataKV, error) {
	var kv GGUFMetadataKV
	// parse key
	str, err := parseGGUFString(r)
	if err != nil {
		return kv, fmt.Errorf("parse key: %w", err)
	}
	kv.Key = str

	// parse value_type (uint32)
	var vt uint32
	if err := binary.Read(r, binary.LittleEndian, &vt); err != nil {
		return kv, fmt.Errorf("read value_type: %w", err)
	}
	kv.ValueType = GgufMetadataValueType(vt)

	// parse value union
	val, err := parseMetadataValue(r, kv.ValueType)
	if err != nil {
		return kv, fmt.Errorf("parse value: %w", err)
	}
	kv.Value = val
	return kv, nil
}

// parseTensorInfo corresponds to gguf_tensor_info_t
//
//	struct gguf_tensor_info_t {
//	  gguf_string_t name;
//	  uint32_t n_dimensions;
//	  uint64_t dimensions[n_dimensions];
//	  ggml_type type;
//	  uint64_t offset;
//	};
func parseTensorInfo(r io.Reader) (GGUFTensorInfo, error) {
	var ti GGUFTensorInfo

	// name (gguf_string_t)
	name, err := parseGGUFString(r)
	if err != nil {
		return ti, fmt.Errorf("parse tensor name: %w", err)
	}
	ti.Name = name

	// n_dimensions
	if err := binary.Read(r, binary.LittleEndian, &ti.NDimensions); err != nil {
		return ti, fmt.Errorf("read n_dimensions: %w", err)
	}

	// read dimension array
	if ti.NDimensions > 16 {
		// somewhat arbitrary check, but official doc says "Currently at most 4"
		return ti, fmt.Errorf("n_dimensions too large: %d", ti.NDimensions)
	}
	dims := make([]uint64, ti.NDimensions)
	for i := uint32(0); i < ti.NDimensions; i++ {
		if err := binary.Read(r, binary.LittleEndian, &dims[i]); err != nil {
			return ti, fmt.Errorf("read dimension[%d]: %w", i, err)
		}
	}
	ti.Dimensions = dims

	// type (ggml_type) is a uint32
	var t uint32
	if err := binary.Read(r, binary.LittleEndian, &t); err != nil {
		return ti, fmt.Errorf("read ggml_type: %w", err)
	}
	ti.Type = GgmlType(t)

	// offset (uint64)
	var off uint64
	if err := binary.Read(r, binary.LittleEndian, &off); err != nil {
		return ti, fmt.Errorf("read offset: %w", err)
	}
	ti.Offset = off

	return ti, nil
}

// parseGGUFString corresponds to gguf_string_t:
// struct gguf_string_t { uint64_t len; char string[len]; }
func parseGGUFString(r io.Reader) (GGUFString, error) {
	var s GGUFString
	if err := binary.Read(r, binary.LittleEndian, &s.Len); err != nil {
		return s, err
	}
	if s.Len > 0xFFFFFFF { // just a sanity check
		return s, errors.New("string length too large")
	}
	buf := make([]byte, s.Len)
	if _, err := io.ReadFull(r, buf); err != nil {
		return s, err
	}
	s.Data = string(buf)
	return s, nil
}

// parseGGUFArray corresponds to the struct inside union for arrays:
//
//	struct {
//	  gguf_metadata_value_type type;  (32-bit)
//	  uint64_t len;                  (number of elements)
//	  gguf_metadata_value_t array[len]; (recursive!)
//	}
func parseGGUFArray(r io.Reader) (GGUFArray, error) {
	var arr GGUFArray
	// read sub-type
	var st uint32
	if err := binary.Read(r, binary.LittleEndian, &st); err != nil {
		return arr, err
	}
	arr.ElemType = GgufMetadataValueType(st)

	// read length
	if err := binary.Read(r, binary.LittleEndian, &arr.Len); err != nil {
		return arr, err
	}

	// parse array[len] elements
	if arr.Len > 0xFFFFFFF {
		return arr, fmt.Errorf("array length too large: %d", arr.Len)
	}

	arr.Elements = make([]GGUFMetadataValue, 0, arr.Len)
	for i := uint64(0); i < arr.Len; i++ {
		// each element is union gguf_metadata_value_t, so parse recursively
		val, err := parseMetadataValue(r, arr.ElemType)
		if err != nil {
			return arr, fmt.Errorf("parse array elem[%d]: %w", i, err)
		}
		arr.Elements = append(arr.Elements, val)
	}

	return arr, nil
}

// parseMetadataValue parses union gguf_metadata_value_t given a known ValueType.
func parseMetadataValue(r io.Reader, vt GgufMetadataValueType) (GGUFMetadataValue, error) {
	var mv GGUFMetadataValue
	mv.Type = vt

	switch vt {
	case GGUF_METADATA_VALUE_TYPE_UINT8:
		var tmp uint8
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return mv, err
		}
		mv.U8 = tmp

	case GGUF_METADATA_VALUE_TYPE_INT8:
		var tmp int8
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return mv, err
		}
		mv.I8 = tmp

	case GGUF_METADATA_VALUE_TYPE_UINT16:
		var tmp uint16
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return mv, err
		}
		mv.U16 = tmp

	case GGUF_METADATA_VALUE_TYPE_INT16:
		var tmp int16
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return mv, err
		}
		mv.I16 = tmp

	case GGUF_METADATA_VALUE_TYPE_UINT32:
		var tmp uint32
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return mv, err
		}
		mv.U32 = tmp

	case GGUF_METADATA_VALUE_TYPE_INT32:
		var tmp int32
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return mv, err
		}
		mv.I32 = tmp

	case GGUF_METADATA_VALUE_TYPE_FLOAT32:
		var tmp float32
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return mv, err
		}
		mv.F32 = tmp

	case GGUF_METADATA_VALUE_TYPE_BOOL:
		// 1-byte bool, 0/1
		var tmp uint8
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return mv, err
		}
		if tmp != 0 && tmp != 1 {
			return mv, fmt.Errorf("invalid bool value: %d", tmp)
		}
		mv.B = tmp == 1

	case GGUF_METADATA_VALUE_TYPE_STRING:
		s, err := parseGGUFString(r)
		if err != nil {
			return mv, err
		}
		mv.StrVal = &s

	case GGUF_METADATA_VALUE_TYPE_ARRAY:
		arr, err := parseGGUFArray(r)
		if err != nil {
			return mv, err
		}
		mv.ArrVal = &arr

	case GGUF_METADATA_VALUE_TYPE_UINT64:
		var tmp uint64
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return mv, err
		}
		mv.U64 = tmp

	case GGUF_METADATA_VALUE_TYPE_INT64:
		var tmp int64
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return mv, err
		}
		mv.I64 = tmp

	case GGUF_METADATA_VALUE_TYPE_FLOAT64:
		var tmp float64
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return mv, err
		}
		mv.F64 = tmp

	default:
		return mv, fmt.Errorf("unknown metadata value type: %d", vt)
	}

	return mv, nil
}
