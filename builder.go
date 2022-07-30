package training

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"training/utils"

	"sigs.k8s.io/yaml"
)

// type Iterator interface {
// 	First()
// 	IsDone() bool
// 	Next() interface{}
// }

type Helper struct {
	parents    []*Helper
	Objects    map[string]interface{}
	ObjectsArr []interface{}
	self       interface{}
}

func (h *Helper) GenerateJson(s interface{}) error {
	h.parents = make([]*Helper, 0)

	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_type := reflect.ValueOf(s).Kind()
	switch _type {
	case reflect.Array, reflect.Slice, reflect.Pointer:
		json.Unmarshal(buf, &h.ObjectsArr)
		h.Objects = nil
	case reflect.Struct:
		json.Unmarshal(buf, &h.Objects)
		h.ObjectsArr = nil
	}

	return nil
}

func (h *Helper) GetArrary(index int, name interface{}) interface{} {
	if h.Objects != nil {
		if len(h.Objects[name.(string)].([]interface{}))-1 < index {
			return nil
		}
		return h.Objects[name.(string)].([]interface{})[index]
	}

	if len(h.ObjectsArr)-1 < index {
		return nil
	}

	return h.ObjectsArr[name.(int)].([]interface{})[index]
}

func (h *Helper) Get(name interface{}) interface{} {
	if h.Objects != nil {
		return h.Objects[name.(string)]
	}

	return h.ObjectsArr[name.(int)]
}

func (h *Helper) Set(name interface{}, values ...interface{}) *Helper {
	var x []interface{}
	var ObjectsIsExited bool

	switch h.Objects {
	case nil:
		var index int = name.(int)
		if len(values) > 1 {
			h.ObjectsArr[index] = make([]interface{}, len(values))
			x = h.ObjectsArr[index].([]interface{})
		}

	default:
		ObjectsIsExited = true
		h.ObjectsArr = nil

		if len(values) > 1 {
			h.Objects[name.(string)] = make([]interface{}, len(values))
			//todo:
			x = h.Objects[name.(string)].([]interface{})
		}
	}

	for i, v := range values {
		var value interface{}
		if _, ok := v.(*Helper); ok {
			if v.(*Helper).Objects == nil {
				value = v.(*Helper).ObjectsArr
			}
		} else {
			value = v
		}

		if len(values) == 1 {
			switch ObjectsIsExited {
			case true:
				h.Objects[name.(string)] = value
				continue
			case false:
				h.ObjectsArr[name.(int)] = value
				continue
			}
		}

		x[i] = value
	}
	return h
}

func (h *Helper) Delete(name interface{}) *Helper {
	if h.Objects != nil {
		delete(h.Objects, name.(string))
		return h
	}

	index := name.(int)
	if index < len(h.ObjectsArr) {
		foo := append(h.ObjectsArr[:index], h.ObjectsArr[index+1:]...)

		switch h.self.(type) {
		case int:
			h.parents[len(h.parents)-1].ObjectsArr[h.self.(int)] = foo
		case string:
			h.parents[len(h.parents)-1].Objects[h.self.(string)] = foo
		}
	}

	return h
}

func (h *Helper) BackRoot() *Helper {
	if len(h.parents) == 0 {
		panic("len(h.parents) == 0")
	}

	root := h.parents[0]
	return &Helper{
		parents:    nil,
		self:       root.self,
		Objects:    root.Objects,
		ObjectsArr: root.ObjectsArr,
	}
}

func (h *Helper) Back() *Helper {
	if len(h.parents) == 0 {
		panic("len(h.parents) == 0")
	}

	lp := h.parents[len(h.parents)-1]
	np := h.parents[:len(h.parents)-1]

	return &Helper{
		parents:    np,
		self:       lp.self,
		Objects:    lp.Objects,
		ObjectsArr: lp.ObjectsArr,
	}
}

func (h *Helper) MarshalPretty() string {
	if h.ObjectsArr == nil {
		buf, _ := json.MarshalIndent(h.Objects, "", "    ")
		return string(buf)
	}

	buf, _ := json.MarshalIndent(h.ObjectsArr, "", "    ")
	return string(buf)
}

func (h *Helper) Marshal() string {
	if h.ObjectsArr == nil {
		buf, _ := json.Marshal(h.Objects)
		return string(buf)
	}

	buf, _ := json.Marshal(h.ObjectsArr)
	return string(buf)
}

func (h *Helper) Start(n interface{}) *Helper {
	return h.Access(n)
}

func (h *Helper) End() *Helper {
	return h.Back()
}

//access deep
func (h *Helper) DeepAccess(n ...interface{}) *Helper {
	p := h
	for _, v := range n {
		p = p.Access(v)
	}
	return p
}

func (h *Helper) Access(n interface{}) *Helper {
	if h.Objects != nil {
		name := n.(string)
		if _, e := h.Objects[name]; !e {
			h.Objects[name] = make(map[string]interface{})
		}

		switch h.Objects[name].(type) {
		case map[string]interface{}:
			return &Helper{
				parents:    append(h.parents, h),
				self:       n,
				Objects:    h.Objects[name].(map[string]interface{}),
				ObjectsArr: nil,
			}

		case []interface{}:
			return &Helper{
				parents:    append(h.parents, h),
				self:       n,
				Objects:    nil,
				ObjectsArr: h.Objects[name].([]interface{}),
			}
		}
	}

	index := n.(int)
	switch h.ObjectsArr[index].(type) {
	case map[string]interface{}:
		return &Helper{
			parents:    append(h.parents, h),
			self:       n,
			Objects:    h.ObjectsArr[index].(map[string]interface{}),
			ObjectsArr: nil,
		}

	case []interface{}:
		return &Helper{
			parents:    append(h.parents, h),
			self:       n,
			Objects:    nil,
			ObjectsArr: h.ObjectsArr[index].([]interface{}),
		}
	}

	return nil
}

func (h *Helper) Json2Yaml(json []byte) (str string, err error) {
	bytes, err := yaml.JSONToYAML(json)
	return utils.Bytes2String(bytes), err
}

func (h *Helper) GetDepth() int {
	return len(h.parents)
}

func (h *Helper) SaveFile(fileName string, data []byte) error {
	return ioutil.WriteFile(fileName, data, 0644)
}

func (h *Helper) GetElementType(v interface{}) reflect.Kind {
	if h.Objects != nil {
		return reflect.TypeOf(h.Objects[v.(string)]).Kind()
	}
	return reflect.TypeOf(h.ObjectsArr[v.(int)]).Kind()
}
