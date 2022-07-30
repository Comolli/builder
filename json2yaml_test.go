package training

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
	"training/utils"

	"github.com/stretchr/testify/assert"
)

var jsonBuilder *Helper

func init() {
	foo := Foo{"xxx", 1}
	jsonBuilder = NewJson2Yaml(foo,
		WithErrHandler(func(err error) {
			panic(err)
		}),
		WithContext(context.Background()),
		WithPrecheckFunc((func(self interface{}) error {
			b, err := json.Marshal(self)
			if err != nil {
				return err
			}

			if !json.Valid(b) {
				return errors.New("invalid json")
			}

			return nil
		}))).Do()
}

type Foo struct {
	Str string `json:"str"`
	Gg  int    `json:"gg"`
}

type Bar struct {
	Aa struct {
		Cc int     `json:"cc"`
		Dd float64 `json:"dd"`
	} `json:"aa"`
	Cc  string `json:"cc"`
	Gg  int    `json:"gg"`
	Str string `json:"str"`
	Xxx struct {
		Bb bool `json:"bb"`
	} `json:"xxx"`
}

func TestJsonBuilderType(t *testing.T) {
	jsonBuilder.Start("xxx").Set("bb", true).End().Set("cc", "dd")
	jsonBuilder.Start("aa").Set("cc", 1).Set("dd", 1.1).End()
	jsonBuilder.Start("kk").Set("arr", 1, 2, 3, 4, 5, false).End()

	assert.Equal(t, reflect.Map, reflect.TypeOf(jsonBuilder.Get("xxx")).Kind())
	assert.Equal(t, reflect.Bool, reflect.TypeOf(jsonBuilder.Get("xxx").(map[string]interface{})["bb"]).Kind())
	assert.Equal(t, reflect.Int, reflect.TypeOf(jsonBuilder.Get("aa").(map[string]interface{})["cc"]).Kind())
	assert.Equal(t, reflect.Float64, reflect.TypeOf(jsonBuilder.Get("aa").(map[string]interface{})["dd"]).Kind())
	assert.Equal(t, reflect.String, reflect.TypeOf(jsonBuilder.Get("cc")).Kind())
	assert.Equal(t, reflect.Slice, reflect.TypeOf(jsonBuilder.Get("kk").(map[string]interface{})["arr"]).Kind())

	t.Log(jsonBuilder.MarshalPretty())
}

func TestJsonBuilderUnMarshall(t *testing.T) {
	//test timeout
	ctx, _ := context.WithTimeout(context.Background(), time.Microsecond/10000)
	helper := NewJson2Yaml("[[[[][][][]/***	}}}}-=", WithErrHandler(func(err error) {
		t.Log(err)
	}),
		WithContext(ctx),
		WithPrecheckFunc(func(self interface{}) error {
			return nil
		})).Do()

	assert.Equal(t, true, assert.Nil(t, helper))

	//test wrong parse
	helper = NewJson2Yaml(math.Inf(1),
		WithErrHandler(func(err error) {
			t.Log(err)
		}),
		WithContext(context.Background()),
		WithPrecheckFunc(func(self interface{}) error {
			return errors.New("invalid json")
		})).Do()

	assert.Equal(t, true, assert.Nil(t, helper))

	jsonBuilder.Start("xxx").Set("bb", true).End().Set("cc", "dd")
	jsonBuilder.Start("aa").Set("cc", 1).Set("dd", 1.1).End()

	bar := &Bar{}
	err := json.Unmarshal(utils.UnsafeGetBytes(jsonBuilder.Marshal()), bar)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, nil, err)
}

func TestJsonBuilderOperation(t *testing.T) {
	kk := []Foo{{"xxx", 1}, {"yyy", 2}}
	helper := Helper{
		self: kk,
	}

	err := helper.GenerateJson(helper.self)
	assert.Equal(t, nil, err)
	t.Log(helper.MarshalPretty())

	//test painic
	var errMsg string
	assert.NotPanics(t, func() {
		defer func() {
			if p := recover(); p != nil {
				errMsg = fmt.Sprintf("%v", p)
			}
		}()
		helper.BackRoot()
	})
	assert.Equal(t, "len(h.parents) == 0", errMsg)

	assert.NotPanics(t, func() {
		defer func() {
			if p := recover(); p != nil {
				errMsg = fmt.Sprintf("%v", p)
			}
		}()
		helper.Back()
	})

	//test marshall
	assert.Equal(t, "len(h.parents) == 0", errMsg)
	assert.Equal(t, true, assert.NotNil(t, helper.Marshal()))

	//test objectArray
	helper.Access(1)
	assert.Equal(t, true, assert.NotNil(t, helper.Access(1)))

	var h *Helper
	jsonBuilder.Start("xxx").
		Set("bb", true).
		End().
		Set("cc", "dd")

	jsonBuilder.Start("aa").
		Set("cc", 1).
		Set("dd", 1.1).
		End()

	//set
	h = jsonBuilder.Access("xxx").Set("kkk", "asss").Access("ll").Set("kk", "kkk").Set("sad", 1, "cousin_2").Back().Set("dd", 1)
	t.Log(jsonBuilder.MarshalPretty())
	helper.Set(1, "1", 1)

	//get
	element := h.DeepAccess("ll").Get("sad").([]interface{})[0]
	assert.Equal(t, 1, element)
	assert.Equal(t, "dd", h.BackRoot().Get("cc"))
	assert.Equal(t, true, assert.NotNil(t, helper.Get(0)))

	element = h.DeepAccess("ll").GetArrary(0, "sad")
	assert.Equal(t, 1, element)
	element = h.DeepAccess("ll").GetArrary(1000, "sad")
	assert.Equal(t, nil, element)

	//update
	h = h.DeepAccess("ll", "sad").Set(0, "updateElement").Back()
	element = h.Get("sad").([]interface{})[0]
	assert.Equal(t, "updateElement", element)
	t.Log(jsonBuilder.MarshalPretty())

	//delete
	h.Access("sad").Delete(0)
	assert.Equal(t, 1, len(h.Get("sad").([]interface{})))
	h.Delete("sad")
	assert.Equal(t, nil, h.Get("sad"))
}

func TestJson2Yaml(t *testing.T) {
	jsonBuilder.Start("xxx").
		Set("bb", true).
		End().
		Set("cc", "dd")

	jsonBuilder.Start("aa").
		Set("cc", 1).
		Set("dd", 1.1).
		End()

	t.Log(jsonBuilder.MarshalPretty())
	j := jsonBuilder.Access("xxx").Set("kkk", "asss").Access("ll").Set("kk", "kkk").Set("sad", 1, "cousin_2").Back().Set("dd", 1)
	t.Log(j.GetDepth())

	//json to yaml
	yaml, err := jsonBuilder.Json2Yaml(utils.UnsafeGetBytes(j.Marshal()))
	assert.Equal(t, nil, err)
	t.Log(yaml)

	j = j.BackRoot()
	t.Log(j.GetDepth())
	yaml, err = jsonBuilder.Json2Yaml(utils.UnsafeGetBytes(j.Marshal()))
	assert.Equal(t, nil, err)
	t.Log(yaml)

	//save json to file
	err = jsonBuilder.SaveFile("./test_save_json.json", []byte(jsonBuilder.MarshalPretty()))
	assert.Equal(t, nil, err)

	//save yaml to file
	err = jsonBuilder.SaveFile("./test_save_yaml.yaml", []byte(yaml))
	assert.Equal(t, nil, err)
}
