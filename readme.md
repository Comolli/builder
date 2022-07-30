# QuickStart
```golang
type Foo struct {
	Str string `json:"str"`
	Gg  int    `json:"gg"`
}
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
    fmt.Println(jsonBuilder.MarshalPretty())
//output:
// "{
//     "str":"xxx",
//     "gg":1,
// }"
```