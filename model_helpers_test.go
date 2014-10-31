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

type unexportedFieldModel struct {
	ID   string `ohm:"id"`
	name string `ohm:"name"`
}

type noIDModel struct {
	Name string `ohm:"name"`
}

type nonStringIDModel struct {
	Name int `ohm:"name"`
}

func TestValidateModel(t *testing.T) {
	var err error
	if err = ValidateModel(&validModel{}); err != nil {
		t.Error(err)
	}

	if err = ValidateModel(&unexportedFieldModel{}); err != NonExportedAttrError {
		t.Error(`unexported fields with ohm tags should make the model invalid`)
	}

	if err = ValidateModel(&noIDModel{}); err != NoIDError {
		t.Error(`models with no ohm:"id" tag should be invalid`)
	}

	if err = ValidateModel(&nonStringIDModel{}); err != NonStringIDError {
		t.Error(`models should be invalid when their ohm:"id" field is not a string`)
	}
}
