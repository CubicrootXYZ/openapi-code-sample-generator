package encoding

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/errors"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
)

// FormDataEncode groups multipart/form-data features
type FormDataEncode struct {
}

// EnocdeParameter encodes the given parameter and its value to multipart/form-data
func (f *FormDataEncode) EnocdeParameter(name string, value interface{}) (string, error) {
	return f.EnocdeValue("", map[string]interface{}{name: value}, nil)
}

// EnocdeValue encodes a single value to amultipart/form-data
func (f *FormDataEncode) EnocdeValue(ref string, value interface{}, meta *types.FormattingMeta) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if meta != nil {
		// New endpoint, reset boundaries to use the same ones across endpoints.
		resetBoundaries()
	}
	err := writer.SetBoundary(getNextBoundary())
	if err != nil {
		return "", err
	}

	if newValue, ok := value.(map[interface{}]interface{}); ok {
		for key, value := range newValue {
			f.writePart(fmt.Sprint(key), value, writer)
		}
		err := writer.Close()
		if err != nil {
			return "", err
		}

		return f.prefix(writer.Boundary(), meta) + body.String(), nil
	}

	if newValue, ok := value.(map[string]interface{}); ok {
		for key, value := range newValue {
			f.writePart(key, value, writer)
		}
		err := writer.Close()
		if err != nil {
			return "", err
		}

		return f.prefix(writer.Boundary(), meta) + body.String(), nil
	}

	log.Error(fmt.Sprint(value))

	log.Warn(fmt.Sprintf("Data type %T is not supported for multipart/form-data", value))
	return "", errors.ErrUnsupportedDataType
}

func (f *FormDataEncode) prefix(boundary string, meta *types.FormattingMeta) string {
	if meta != nil {
		meta.FormData.OuterBoundary = &boundary
	}
	return "Content-Type: multipart/form-data; boundary=" + boundary + "\r\n\r\n"
}

func (f *FormDataEncode) writePart(key string, value interface{}, writer *multipart.Writer) {
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "text/plain")
	metadataHeader.Set("Content-ID", key)
	part, _ := writer.CreatePart(metadataHeader)

	_, ok1 := value.(map[interface{}]interface{})
	_, ok2 := value.(map[string]interface{})

	if ok1 || ok2 {
		newVal, err := f.EnocdeValue("", value, nil)
		if err == nil {
			_, _ = part.Write([]byte(newVal))
		} else {
			log.Warn(err.Error())
		}
	} else {
		_, _ = fmt.Fprint(part, value)
	}
}

var currentBoundaryIndex = -1
var boundaries = []string{
	"ueng3raexexiFah9xae6keejie1oac7aid8thoe2aineiw4ohcip6aef4kee",
	"oviongienesei6eel0uqueipungaixaecoh4wiezighahPh9aiShiokahy6h",
	"iesh2FoonaloPhic0wohs8aey7thie0tuofahy8ahsh1unaeQueifohmoh5i",
	"ul9ni9is2aiPhie7comejaW7eeb2sheekaaP9na2UXaPuh4pez9kein0aiw1",
	"eniiBeedaefiu2eeh9aiG6ool8cheequ8the0raipeeniuP9ja8ohPhooSu8",
	"ootaigiephei8avupuik5theida9Za8ezoB9ung7bee9Rieyauf6Ech5Chuo",
	"oongie9phuyooTh8aeyaiC8eequaigh3ge1paeth9vataok7Vidaimeiquoo",
	"inoo6op7itee0iuch8uecoipoChai4iexu5ex9kieshaithohae0jozaipee",
	"Aehieju8aiziTh8iigh7rai4ainahchae3ocheehud2kohreexo6ooxei8Sa",
	"Aek3eij4cei8Eehee9eemee7roPh0ahX6aim9Shaitho6vi7vee9phah5ahN",
	"aidi1iudeichie8au7aizeowaiYail8ugh7aino8ohtoishai7eizailahGh",
	"Shailahvie2ok9aigh8Huus7siet6weeshein4Ii5puoThoh8kaek4no5rie",
	"aPhooveucheireiQuooxohpha3Oshohthae8ka4Oosie4egh0goos9cisooH",
	"eeQuaphieb4eis6OhngooBil0le8Auw4gaeyooth3oomi1aiheiphaemae5h",
	"shooNah7iejiChahth1koon6soqu8revopo7ooch6utayedouN7aeshoo3oa",
	"daeng2koo0ma2Keisuy7peirecaiGhee6oSho6oojaeJaechoo1bielohvah",
	"ahGee5ca9maej8jiesahph6eiGa9aPh2Eequuvoh2xukuu7Faed0aePh0ieZ",
	"Ahnoo9Ixoo1id6Ohth7ief3neiChoh1webokahv4eeM3Bieshi2LeeW3Roka",
	"Eetu2Eephie7EiJ0du0aeNohZ7ai3aiwi7gue2xahhodaJaht9xuquushieh",
	"ia7aiGh1ieL5jei8aex4pe9zee0quairo6te3eijaighiek3phehaith0oth",
}

// getNextBoundary is a workaround to use the same boundaries across examples.
func getNextBoundary() string {
	currentBoundaryIndex++
	if currentBoundaryIndex >= len(boundaries) {
		currentBoundaryIndex = 0
	}

	return boundaries[currentBoundaryIndex]
}

func resetBoundaries() {
	currentBoundaryIndex = -1
}
