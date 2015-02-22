package gohm

import(
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/pote/go-msgpack"
	"reflect"
)


func (c *Connection) Save(model interface{}) (error) {
	if err := validateModel(model); err != nil {
		return err
	}

	modelData := reflect.ValueOf(model).Elem()
	modelType := modelData.Type()

	// Prepare Ohm-scripts `features` parameter.
	features := map[string]string{
		"name": modelType.Name(),
	}
	if modelID(model) != "" {
		features["id"] = modelID(model)
	}
	ohmFeatures, err := msgpack.Marshal(features)
	if err != nil {
		return err
	}

	// Prepare Ohm-scripts `attributes` parameter.
	attrs := []string{}
	attrIndexMap := modelAttrIndexMap(model)
	for attr, index := range attrIndexMap {
		attrs = append(attrs, attr)
		attrs = append(attrs, modelData.Field(index).String())
	}
	ohmAttrs, err := msgpack.Marshal(attrs)
	if err != nil {
		return err
	}

	// Prepare Ohm-scripts `indices` parameter.
	ohmIndices, err := msgpack.Marshal(modelIndices(model))
	if err != nil {
		return err
	}

	// TODO
	// Prepare Ohm-scripts `uniques` parameter.
	ohmUniques, err := msgpack.Marshal(&map[string]string{})
	if err != nil {
		return err
	}

	conn := c.RedisPool.Get()
	defer conn.Close()
	id, err :=  redis.String(c.luaSave.Do(conn, ohmFeatures, ohmAttrs, ohmIndices, ohmUniques))
	if err != nil {
		return err
	}
	modelSetID(id, model)

	return nil
}

func (c *Connection) Load(model interface{}) (err error) {
	if err := validateModel(model); err != nil {
		return err
	}

	conn := c.RedisPool.Get()
	defer conn.Close()

	attrs, err := redis.Strings(conn.Do("HGETALL", modelKey(model)))
	if err != nil {
		return
	}
	if len(attrs) == 0 {
		err = errors.New(`Couldn't find "` + modelKey(model) + `" in redis`)
		return
	}
	modelLoadAttrs(attrs, model)

	return
}
