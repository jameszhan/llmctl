package main

import (
	"fmt"
	"os"

	"github.com/jameszhan/llmctl/common/gguf"
	"github.com/spf13/cobra"
)

var mainCommand = &cobra.Command{
	Use:    "llmctl",
	PreRun: preRun,
	Run:    run,
}

func init() {
	mainCommand.Flags().StringP("file", "f", "", "GGUF file path")
	mainCommand.Flags().BoolP("verbose", "V", false, "Verbose mode")
}

func main() {
	if err := mainCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func preRun(cmd *cobra.Command, args []string) {

}

func run(cmd *cobra.Command, args []string) {
	// verbose, _ := cmd.Flags().GetBool("verbose")
	filepath, _ := cmd.Flags().GetString("file")
	if filepath == "" {
		fmt.Println("invalid GGUF file:")
		return
	}

	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer f.Close()

	gguf, err := gguf.Parse(f)
	if err != nil {
		fmt.Println("Parse GGUF error:", err)
		return
	}

	fmt.Printf("Parsed GGUF: version=%d, tensor_count=%d, metadata_kv_count=%d\n",
		gguf.Version, gguf.TensorCount, gguf.MetadataKVCount)

	// Print some of the metadata
	for i, kv := range gguf.Metadata {
		//fmt.Printf("Metadata[%d] Key=%q, Type=%d\n", i, kv.Key.Data, kv.ValueType)
		fmt.Printf("Metadata[%d] Key=%q => Value: %v\n", i, kv.Key.Data, kv.Value)
	}

	// Print some of the tensor info
	for i, t := range gguf.Tensors {
		fmt.Printf("Tensor[%d]: name=%q, dims=%v, type=%v, offset=%d\n",
			i, t.Name.Data, t.Dimensions, t.Type, t.Offset)
	}
}
