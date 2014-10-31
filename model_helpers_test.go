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
	if err = validateModel(&validModel{}); err != nil {
		t.Error(err)
	}

	if err = validateModel(&unexportedFieldModel{}); err != NonExportedAttrError {
		t.Error(`unexported fields with ohm tags should make the model invalid`)
	}

	if err = validateModel(&noIDModel{}); err != NoIDError {
		t.Error(`models with no ohm:"id" tag should be invalid`)
	}

	if err = validateModel(&nonStringIDModel{}); err != NonStringIDError {
		t.Error(`models should be invalid when their ohm:"id" field is not a string`)
	}
}

func TestModelAttrIndexMap(t *testing.T) {
	attrMap := modelAttrIndexMap(&validModel{})

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

	if modelID(u) != `` {
		t.Errorf(`expected model ID to be empty, but its set to "%v"`, modelID(u))
	}

	if modelID(u2) != `2` {
		t.Errorf(`model ID should be 2, but its "%v"`, modelID(u))
	}
}

func TestModelHasAttribute(t *testing.T) {
	if !modelHasAttribute(&validModel{}, `email`) {
		t.Error(`model has attribute "email", but the function return false`)
	}

	if modelHasAttribute(&validModel{}, `palangana`) {
		t.Error(`model doesnt have the attribute "palangana", but the function return true`)
	}
}

func TestModelIDFieldName(t *testing.T) {
	if modelIDFieldName(&validModel{}) != `ID` {
		t.Error(`function is not correctly reporting the ID field name`)
	}
}

func TestModelType(t *testing.T) {
	if modelType(&validModel{}) != `validModel` {
		t.Error(`function does not return correct model name`)
	}
}
