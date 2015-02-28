package gohm

import(
	`reflect`
	`testing`
)

type validModel struct {
	ID    string `ohm:"id"`
	Name  string `ohm:"name"`
	Email string `ohm:"email index"`
	UUID  string `ohm:"uuid  unique"`
}

type unexportedFieldModel struct {
	ID   string `ohm:"id"`
	name string `ohm:"name"`
}

type noIDModel struct {
	Name string `ohm:"name"`
}

type nonStringIDModel struct {
	Name int `ohm:"id"`
}

type modelWithReference struct {
	ID                string `ohm:"id"`
	Email             string `ohm:"email index"`
	ExternalModelID   string `ohm:"external_model_id reference ExternalModel"`
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

func TestModelIndices(t *testing.T) {
  u := validModel{Email: "pote@tardis.com.uy"}
  indices := modelIndices(&u)

  if len(indices) == 0 {
    t.Error("there were indicices in the model but function didn't return them")
    return
  }

  if indices["email"] != "pote@tardis.com.uy" {
    t.Error("incorrect value for index")
  }
}

func TestModelReferences(t *testing.T) {
  u := modelWithReference{ExternalModelID: "42"}
  references := modelReferences(&u)

  if len(references) != 1 {
    t.Error("there were indicices in the model but function didn't return them")
    return
  }

  if references[0]["value"] != "42" {
    t.Errorf("incorrect value for reference: %v", references[0]["value"])
  }

  if references[0]["tag"] != "external_model_id" {
    t.Errorf("incorrect tag for reference: %v", references[0]["tag"])
  }

  if references[0]["referencee"] != "ExternalModel" {
    t.Errorf("incorrect referencee name for reference: %v", references[0]["referencee"])
  }
}

func TestRefernecesAreIndicesToo(t *testing.T) {
  u := modelWithReference{
    Email:            "pote@tardis.com.uy",
    ExternalModelID:  "42",
  }

  indices := modelIndexMap(&u)

  if len(indices) != 2 {
    t.Error("not enough indices returned by function")
    return
  }

  if indices["email"][0] != "pote@tardis.com.uy" {
    t.Error("regular indices are not included in modelIndexMap")
  }

  if indices["external_model_id"][0] != "42" {
    t.Error("reference is not treated as index")
  }
}
