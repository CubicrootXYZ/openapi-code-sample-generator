package cmd

import (
	"openapi-code-sample-generator/internal/codesample"
	"openapi-code-sample-generator/internal/log"
	"os"

	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates the code samples",
	Long:  "Adds code samples to the inputed OpenApi spec file",
	Run:   generate,
}

// Flags
var inputFile string
var outputFile string
var debug bool

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	generateCmd.PersistentFlags().StringVar(&inputFile, "input-file", "swagger.yaml", "Location of the input swagger yaml specification file")
	generateCmd.PersistentFlags().StringVar(&outputFile, "output-file", "swagger-out.yaml", "Location of the output swagger yaml specification file")
	generateCmd.PersistentFlags().BoolVar(&debug, "v", false, "Enable to get verbose output")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generate(cmd *cobra.Command, args []string) {
	log.Info("Loading file " + inputFile)

	doc, err := openapi.NewLoader().LoadFromFile(inputFile)
	if err != nil {
		log.Error("Failed to open file: " + err.Error())
		os.Exit(1)
	}

	constructor := codesample.NewConstructor(doc, debug)
	constructor.AddSamples([]codesample.Language{codesample.LanguageCurl})

	json, err := yaml.Marshal(doc)
	os.WriteFile(outputFile, json, 0666)
}
