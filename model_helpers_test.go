package gohm

import(
	`reflect`
	`testing`
)

type validModel struct {
	ID    string `ohm:"id"`
	Name  string `ohm:"name"`
	Email string `ohm:"email index"`
	UUID  string `ohm:"uuid unique"`
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

func TestModelAttrIndexMap(t *testing.T) {
	attrMap := ModelAttrIndexMap(&validModel{})

	expectedMap := map[string]int{
		`name`:  1,
		`email`: 2,
		`uuid`:  3,
	}

	if !reflect.DeepEqual(expectedMap, attrMap) {
		t.Errorf(`expected %v, got %v`, expectedMap, attrMap)
	}
}

func TestModelID(t *testing.T) {
	u := &validModel{}
	u2 := &validModel{ID: `2`}

	if ModelID(u) != `` {
		t.Errorf(`expected model ID to be empty, but its set to "%v"`, ModelID(u))
	}

	if ModelID(u2) != `2` {
		t.Errorf(`model ID should be 2, but its "%v"`, ModelID(u))
	}
}

func TestModelHasAttribute(t *testing.T) {
	if !ModelHasAttribute(&validModel{}, `email`) {
		t.Error(`model has attribute "email", but the function return false`)
	}

	if ModelHasAttribute(&validModel{}, `palangana`) {
		t.Error(`model doesnt have the attribute "palangana", but the function return true`)
	}
}
