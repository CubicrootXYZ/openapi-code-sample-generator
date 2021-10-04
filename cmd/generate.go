package cmd

import (
	"openapi-sample-generator/internal/codesample"
	"openapi-sample-generator/internal/log"
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

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	generateCmd.PersistentFlags().StringVar(&inputFile, "input-file", "swagger.yaml", "Location of the input swagger yaml specification file")
	generateCmd.PersistentFlags().StringVar(&outputFile, "output-file", "swagger-out.yaml", "Location of the output swagger yaml specification file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generate(cmd *cobra.Command, args []string) {
	log.Info("Loading file " + inputFile)

	doc, err := openapi.NewLoader().LoadFromFile(inputFile)
	if err != nil {
		log.Error("Failed to open file: " + err.Error())
		return
	}

	for _, path := range doc.Paths {
		if path.Post != nil {
			path.Post.ExtensionProps.Extensions = make(map[string]interface{})
			path.Post.ExtensionProps.Extensions["x-codeSamples"] = codesample.GetSamples(codesample.LanguageCurl, path.Post, path)
		}
	}

	json, err := yaml.Marshal(doc)
	os.WriteFile(outputFile, json, 0666)
}
