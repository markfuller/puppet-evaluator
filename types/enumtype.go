package types

import (
	"io"

	"github.com/lyraproj/puppet-evaluator/eval"
	"github.com/lyraproj/puppet-evaluator/utils"
	"reflect"
	"strings"
)

type EnumType struct {
	caseInsensitive bool
	values          []string
}

var EnumMetaType eval.ObjectType

func init() {
	EnumMetaType = newObjectType(`Pcore::EnumType`,
		`Pcore::ScalarDataType {
	attributes => {
		values => Array[String[1]],
		case_insensitive => {
			type => Boolean,
			value => false
		}
	}
}`, func(ctx eval.Context, args []eval.Value) eval.Value {
			return NewEnumType2(args...)
		})
}

func DefaultEnumType() *EnumType {
	return enumType_DEFAULT
}

func NewEnumType(enums []string, caseInsensitive bool) *EnumType {
	if caseInsensitive {
		top := len(enums)
		if top > 0 {
			lce := make([]string, top)
			for i, v := range enums {
				lce[i] = strings.ToLower(v)
			}
			enums = lce
		}
	}
	return &EnumType{caseInsensitive, enums}
}

func NewEnumType2(args ...eval.Value) *EnumType {
	return NewEnumType3(WrapValues(args))
}

func NewEnumType3(args eval.List) *EnumType {
	if args.Len() == 0 {
		return DefaultEnumType()
	}
	var enums []string
	top := args.Len()
	caseInsensitive := false
	first := args.At(0)
	if top == 1 {
		switch first.(type) {
		case stringValue:
			enums = []string{first.String()}
		case *ArrayValue:
			return NewEnumType3(first.(*ArrayValue))
		default:
			panic(NewIllegalArgumentType2(`Enum[]`, 0, `String or Array[String]`, args.At(0)))
		}
	} else {
		if ar, ok := first.(*ArrayValue); ok {
			enumArgs := ar.AppendTo(make([]eval.Value, 0, ar.Len()+top-1))
			for i := 1; i < top; i++ {
				enumArgs = append(enumArgs, args.At(i))
			}
			if len(enumArgs) == 0 {
				return DefaultEnumType()
			}
			args = WrapValues(enumArgs)
			top = args.Len()
		}

		enums = make([]string, top)
		args.EachWithIndex(func(arg eval.Value, idx int) {
			str, ok := arg.(stringValue)
			if !ok {
				if ci, ok := arg.(booleanValue); ok && idx == top-1 {
					caseInsensitive = ci.Bool()
					return
				}
				panic(NewIllegalArgumentType2(`Enum[]`, idx, `String`, arg))
			}
			enums[idx] = string(str)
		})
	}
	return NewEnumType(enums, caseInsensitive)
}

func (t *EnumType) Accept(v eval.Visitor, g eval.Guard) {
	v(t)
}

func (t *EnumType) Default() eval.Type {
	return enumType_DEFAULT
}

func (t *EnumType) Equals(o interface{}, g eval.Guard) bool {
	if ot, ok := o.(*EnumType); ok {
		return t.caseInsensitive == ot.caseInsensitive && len(t.values) == len(ot.values) && utils.ContainsAllStrings(t.values, ot.values)
	}
	return false
}

func (t *EnumType) Generic() eval.Type {
	return enumType_DEFAULT
}

func (t *EnumType) Get(key string) (eval.Value, bool) {
	switch key {
	case `values`:
		return WrapValues(t.pvalues()), true
	case `case_insensitive`:
		return booleanValue(t.caseInsensitive), true
	default:
		return nil, false
	}
}

func (t *EnumType) IsAssignable(o eval.Type, g eval.Guard) bool {
	if len(t.values) == 0 {
		switch o.(type) {
		case *stringType, *vcStringType, *scStringType, *EnumType, *PatternType:
			return true
		}
		return false
	}

	if st, ok := o.(*vcStringType); ok {
		return eval.IsInstance(t, stringValue(st.value))
	}

	if en, ok := o.(*EnumType); ok {
		oEnums := en.values
		if len(oEnums) > 0 && (t.caseInsensitive || !en.caseInsensitive) {
			for _, v := range en.values {
				if !eval.IsInstance(t, stringValue(v)) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (t *EnumType) IsInstance(o eval.Value, g eval.Guard) bool {
	if str, ok := o.(stringValue); ok {
		if len(t.values) == 0 {
			return true
		}
		s := string(str)
		if t.caseInsensitive {
			s = strings.ToLower(s)
		}
		for _, v := range t.values {
			if v == s {
				return true
			}
		}
	}
	return false
}

func (t *EnumType) MetaType() eval.ObjectType {
	return EnumMetaType
}

func (t *EnumType) Name() string {
	return `Enum`
}

func (t *EnumType) ReflectType(c eval.Context) (reflect.Type, bool) {
	return reflect.TypeOf(`x`), true
}

func (t *EnumType) String() string {
	return eval.ToString2(t, NONE)
}

func (t *EnumType) Parameters() []eval.Value {
	result := t.pvalues()
	if t.caseInsensitive {
		result = append(result, BooleanTrue)
	}
	return result
}

func (t *EnumType) CanSerializeAsString() bool {
	return true
}

func (t *EnumType) SerializationString() string {
	return t.String()
}

func (t *EnumType) ToString(b io.Writer, f eval.FormatContext, g eval.RDetect) {
	TypeToString(t, b, f, g)
}

func (t *EnumType) PType() eval.Type {
	return &TypeType{t}
}

func (t *EnumType) pvalues() []eval.Value {
	top := len(t.values)
	if top == 0 {
		return eval.EMPTY_VALUES
	}
	v := make([]eval.Value, top)
	for idx, e := range t.values {
		v[idx] = stringValue(e)
	}
	return v
}

var enumType_DEFAULT = &EnumType{false, []string{}}
