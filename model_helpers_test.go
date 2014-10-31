package gohm

import(
	"testing"
)

type validModel struct {
	ID    string `ohm:"id"`
	Name  string `ohm:"name"`
	Email string `ohm:"email index"`
	UUID  string `ohm:"name unique"`
}

type unexporedFieldModel struct {
	ID   string `ohm:"id"`
	name string `ohm:"name"`
}
type noIDModel struct {
	Name string `ohm:"name"`
}

func TestValidateModel(t *testing.T) {
	var err error
	if err = ValidateModel(&validModel{}); err != nil {
		t.Error(err)
	}

}
