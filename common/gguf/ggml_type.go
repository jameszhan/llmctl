package gguf

import "fmt"

// ggml_type (for tensor data), mapped to a uint32
type GgmlType uint32

const (
	GGML_TYPE_F32     GgmlType = iota
	GGML_TYPE_F16              = 1
	GGML_TYPE_Q4_0             = 2
	GGML_TYPE_Q4_1             = 3
	GGML_TYPE_Q4_2             = 4 // support has been removed
	GGML_TYPE_Q4_3             = 5 // support has been removed
	GGML_TYPE_Q5_0             = 6
	GGML_TYPE_Q5_1             = 7
	GGML_TYPE_Q8_0             = 8
	GGML_TYPE_Q8_1             = 9
	GGML_TYPE_Q2_K             = 10
	GGML_TYPE_Q3_K             = 11
	GGML_TYPE_Q4_K             = 12
	GGML_TYPE_Q5_K             = 13
	GGML_TYPE_Q6_K             = 14
	GGML_TYPE_Q8_K             = 15
	GGML_TYPE_IQ2_XXS          = 16
	GGML_TYPE_IQ2_XS           = 17
	GGML_TYPE_IQ3_XXS          = 18
	GGML_TYPE_IQ1_S            = 19
	GGML_TYPE_IQ4_NL           = 20
	GGML_TYPE_IQ3_S            = 21
	GGML_TYPE_IQ2_S            = 22
	GGML_TYPE_IQ4_XS           = 23
	GGML_TYPE_I8               = 24
	GGML_TYPE_I16              = 25
	GGML_TYPE_I32              = 26
	GGML_TYPE_I64              = 27
	GGML_TYPE_F64              = 28
	GGML_TYPE_IQ1_M            = 29
	GGML_TYPE_COUNT            = 30
)

func (t GgmlType) String() string {
	switch t {
	case GGML_TYPE_F32:
		return "F32"
	case GGML_TYPE_F16:
		return "F16"
	case GGML_TYPE_Q4_0:
		return "Q4_0"
	case GGML_TYPE_Q4_1:
		return "Q4_1"
	case GGML_TYPE_Q4_2:
		return "Q4_2"
	case GGML_TYPE_Q4_3:
		return "Q4_3"
	case GGML_TYPE_Q5_0:
		return "Q5_0"
	case GGML_TYPE_Q5_1:
		return "Q5_1"
	case GGML_TYPE_Q8_0:
		return "Q8_0"
	case GGML_TYPE_Q8_1:
		return "Q8_1"
	case GGML_TYPE_Q2_K:
		return "Q2_K"
	case GGML_TYPE_Q3_K:
		return "Q3_K"
	case GGML_TYPE_Q4_K:
		return "Q4_K"
	case GGML_TYPE_Q5_K:
		return "Q5_K"
	case GGML_TYPE_Q6_K:
		return "Q6_K"
	case GGML_TYPE_Q8_K:
		return "Q8_K"
	case GGML_TYPE_IQ2_XXS:
		return "IQ2_XXS"
	case GGML_TYPE_IQ2_XS:
		return "IQ2_XS"
	case GGML_TYPE_IQ3_XXS:
		return "IQ3_XXS"
	case GGML_TYPE_IQ1_S:
		return "IQ1_S"
	case GGML_TYPE_IQ4_NL:
		return "IQ4_NL"
	case GGML_TYPE_IQ3_S:
		return "IQ3_S"
	case GGML_TYPE_IQ2_S:
		return "IQ2_S"
	case GGML_TYPE_IQ4_XS:
		return "IQ4_XS"
	case GGML_TYPE_I8:
		return "I8"
	case GGML_TYPE_I16:
		return "I16"
	case GGML_TYPE_I32:
		return "I32"
	case GGML_TYPE_I64:
		return "I64"
	case GGML_TYPE_F64:
		return "F64"
	case GGML_TYPE_IQ1_M:
		return "IQ1_M"
	case GGML_TYPE_COUNT:
		return "COUNT"
	default:
		return fmt.Sprintf("UNKNOWN_TYPE(%d)", uint32(t))
	}
}
