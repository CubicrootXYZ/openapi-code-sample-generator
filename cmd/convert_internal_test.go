package cmd

import (
	"testing"
)

func Test_convert(t *testing.T) {
	t.Run("petstore_minimal.yaml", func(t *testing.T) {
		inputFile = "./testdata/petstore_minimal.yaml"
		convertFile = "./testdata/petstore_minimal.yaml"
		fromVersion = 2
		toVersion = 3
		outputFile = t.TempDir() + "out.yaml"

		convert(nil, nil)
	})
}
