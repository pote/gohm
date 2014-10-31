package gohm

import(
	"github.com/garyburd/redigo/redis"
	"github.com/pote/go-msgpack"
	"github.com/pote/redisurl"
	"reflect"
)

type gohm struct {
	RedisPool *redis.Pool
	LuaSave   *redis.Script
}

func NewDefaultGohm() (*gohm, error) {
	pool, err := redisurl.NewPool(3, 200, "240s")
	if err != nil {
		return nil, err
	}

	return NewGohm(pool), nil
}

func NewGohm(r *redis.Pool) (*gohm) {
	g := &gohm{
		RedisPool: r,
	}

	g.LuaSave = redis.NewScript(0, LUA_SAVE)

	return g
}

func (g *gohm) Save(model interface{}, indices... string) (error) {
	if err := ValidateModel(model); err != nil {
		return err
	}

	var hasID bool
	var idFieldIndex int
	modelData := reflect.ValueOf(model).Elem()
	modelType := modelData.Type()

	// Prepare Ohm-scripts `attributes` parameter.
	attrs := []string{}
	for i := 0; i < modelData.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get("ohm")
		if !hasID && tag == "id" {
			hasID = true
			idFieldIndex = i
			continue
		}
		if tag == "" || tag == "-" {
			continue
		}

		attrs = append(attrs, tag)
		attrs = append(attrs, modelData.Field(i).String())
	}
	ohmAttrs, err := msgpack.Marshal(attrs)
	if err != nil {
		return err
	}

	// Prepare Ohm-scripts `features` parameter.
	features := map[string]string{
		"name": modelType.Name(),
	}
	if modelData.Field(idFieldIndex).String() != "" {
		features["id"] = modelData.FieldByName("ID").String()
	}
	ohmFeatures, err := msgpack.Marshal(features)
	if err != nil {
		return err
	}

	// TODO
	// Prepare Ohm-scripts `indices` parameter.
	ohmIndices, err := msgpack.Marshal(&map[string]string{})
	if err != nil {
		return err
	}

	// TODO
	// Prepare Ohm-scripts `uniques` parameter.
	ohmUniques, err := msgpack.Marshal(&map[string]string{})
	if err != nil {
		return err
	}

	conn := g.RedisPool.Get()
	defer conn.Close()
	id, err :=  redis.String(g.LuaSave.Do(conn, ohmFeatures, ohmAttrs, ohmIndices, ohmUniques))
	if err != nil {
		return err
	}

	modelData.FieldByName(ModelIDFieldName(model)).SetString(id)

	return nil
}

func (g *gohm) Find(id string, model interface{}) (err error) {
	if err := ValidateModel(model); err != nil {
		return err
	}

	conn := g.RedisPool.Get()
	defer conn.Close()
	modelData := reflect.ValueOf(model).Elem()
	modelType := modelData.Type()

	idFieldName := ModelIDFieldName(model)
	modelData.FieldByName(idFieldName).SetString(id)


	attrs, err := redis.Strings(conn.Do("HGETALL", ModelKey(model)))
	if err != nil {
		return
	}

	attrIndexMap := ModelAttrIndexMap(model)
	for i := 0; i < len(attrs); i = i + 2 {
		attrName := attrs[i]
		attrValue := attrs[i + 1]
		attrIndex := attrIndexMap[attrName]

		if ModelHasAttribute(model, attrName) {
			attrValueValue := reflect.ValueOf(attrValue)
			typedAttrValue := attrValueValue.Convert(modelType.Field(attrIndex).Type)
			modelData.Field(attrIndex).Set(typedAttrValue)
		}
	}

	return
}
