package cmd

import (
	"io/ioutil"
	"os"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Converts OpenAPI specifications",
	Long:  "Converts to different versions of OpenAPI specifications",
	Run:   convert,
}

// Flags
var fromVersion int
var toVersion int
var convertFile string

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.PersistentFlags().IntVar(&fromVersion, "from-version", 2, "OpenAPI version the file uses")
	convertCmd.PersistentFlags().IntVar(&toVersion, "to-version", 3, "OpenAPI version to convert to")
	convertCmd.PersistentFlags().StringVar(&convertFile, "file", "swagger.yaml", "the file to convert")
	convertCmd.PersistentFlags().StringVar(&outputFile, "output-file", "swagger-out.yaml", "Location of the output swagger yaml specification file")
	convertCmd.PersistentFlags().BoolVar(&debug, "v", false, "Enable to get verbose output")
}

func convert(cmd *cobra.Command, args []string) {
	log.Verbose = debug

	if fromVersion != 2 || toVersion != 3 {
		log.Warn("Only the following conversions are supported: \n2 ==> 3")
		return
	}

	log.Info("Loading file " + inputFile)
	var doc openapi2.T
	input, err := ioutil.ReadFile(convertFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(input, &doc)
	if err != nil {
		log.Error("Failed to open file: " + err.Error())
		os.Exit(1)
	}

	newDoc, err := openapi2conv.ToV3(&doc)

	log.Info("Writting file " + outputFile)
	json, err := yaml.Marshal(newDoc)
	if err != nil {
		log.Error(err.Error())
	}
	os.WriteFile(outputFile, json, 0666)
}
