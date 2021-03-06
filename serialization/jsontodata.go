package serialization

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/lyraproj/issue/issue"
	"github.com/lyraproj/puppet-evaluator/eval"
	"github.com/lyraproj/puppet-evaluator/types"
)

const firstInArray = 0
const firstInObject = 1
const afterElement = 2
const afterValue = 3
const afterKey = 4

// JsonToData reads JSON from the given reader and streams the values to the
// given ValueConsumer
func JsonToData(path string, in io.Reader, consumer ValueConsumer) {
	defer func() {
		if r := recover(); r != nil {
			panic(eval.Error(eval.EVAL_TASK_BAD_JSON, issue.H{`path`: path, `detail`: r}))
		}
	}()
	d := json.NewDecoder(in)
	d.UseNumber()
	jsonValues(consumer, d)
}

func jsonValues(c ValueConsumer, d *json.Decoder) {
	for {
		t, err := d.Token()
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		if dl, ok := t.(json.Delim); ok {
			ds := dl.String()
			if ds == `}` || ds == `]` {
				return
			}
			if ds == `{` {
				t = nil
				if d.More() {
					t, err = d.Token()
					if err != nil {
						panic(err)
					}
					if ds, ok = t.(string); ok && ds == PCORE_REF_KEY && d.More() {
						t, err = d.Token()
						if err != nil {
							panic(err)
						}
						var n int64
						n, err = t.(json.Number).Int64()
						if err != nil {
							panic(err)
						}
						// Consume end delimiter
						t, err = d.Token()
						if err != nil {
							panic(err)
						}
						if dl, ok = t.(json.Delim); ok && dl.String() == `}` {
							c.AddRef(int(n))
						} else {
							panic(fmt.Errorf("invalid token %T %v", t, t))
						}
						continue
					}
					c.AddHash(8, func() {
						addValue(c, t)
						jsonValues(c, d)
					})
				} else {
					c.AddHash(8, func() {
						jsonValues(c, d)
					})
				}
			} else {
				c.AddArray(8, func() {
					jsonValues(c, d)
				})
			}
		} else {
			addValue(c, t)
		}
	}
}

func addValue(c ValueConsumer, t json.Token) {
	switch t.(type) {
	case bool:
		c.Add(types.WrapBoolean(t.(bool)))
	case float64:
		c.Add(types.WrapFloat(t.(float64)))
	case json.Number:
		if i, err := t.(json.Number).Int64(); err == nil {
			c.Add(types.WrapInteger(i))
		} else {
			f, _ := t.(json.Number).Float64()
			c.Add(types.WrapFloat(f))
		}
	case string:
		c.Add(types.WrapString(t.(string)))
	case nil:
		c.Add(eval.UNDEF)
	}
}
