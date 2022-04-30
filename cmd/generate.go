package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/encoding"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/extractor"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/languages"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/codesample"

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
var selectedLanguages string
var debug bool

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	generateCmd.PersistentFlags().StringVar(&inputFile, "input-file", "swagger.yaml", "Location of the input swagger yaml specification file")
	generateCmd.PersistentFlags().StringVar(&outputFile, "output-file", "swagger-out.yaml", "Location of the output swagger yaml specification file")
	generateCmd.PersistentFlags().StringVar(&selectedLanguages, "languages", "curl,php,js", "Comma separated list of languages to make examples for")
	generateCmd.PersistentFlags().BoolVar(&debug, "v", false, "Enable to get verbose output")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generate(cmd *cobra.Command, args []string) {
	log.Verbose = debug

	encoders := encoding.Encoders()
	extractor := extractor.NewOpenAPIExtractor()
	generators := languages.Generators(encoders, extractor)

	log.Info("Loading file " + inputFile)
	doc, err := openapi.NewLoader().LoadFromFile(inputFile)
	if err != nil {
		log.Error("Failed to open file: " + err.Error())
		os.Exit(1)
	}

	executor := codesample.NewExecutor(doc, generators)
	executor.AddSamples(languagesFromCSV(selectedLanguages))

	json, err := yaml.Marshal(doc)
	if err != nil {
		log.Error(err.Error())
	}

	log.Info(fmt.Sprintf("Writing to file: %s", outputFile))
	os.WriteFile(outputFile, json, 0666)
}

func languagesFromCSV(commaSeparatedLanguages string) []types.Language {
	languages := make([]types.Language, 0)
	for _, language := range strings.Split(commaSeparatedLanguages, ",") {
		newLang := types.StringToLanguage(language)
		if newLang == types.LanguageEmpty {
			log.Error(fmt.Sprintf("Language %s unknown", language))
			os.Exit(1)
		}
		languages = append(languages, newLang)
	}

	return languages
}
