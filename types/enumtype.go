package types

import (
	. "io"

	. "github.com/puppetlabs/go-evaluator/utils"
	. "github.com/puppetlabs/go-evaluator/evaluator"
)

type EnumType struct {
	values []string
}

func DefaultEnumType() *EnumType {
	return enumType_DEFAULT
}

func NewEnumType(enums []string) *EnumType {
	return &EnumType{enums}
}

func NewEnumType2(args ...PValue) *EnumType {
	if len(args) == 0 {
		return DefaultEnumType()
	}
	var enums []string
	top := len(args)
	if top == 1 {
		first := args[0]
		switch first.(type) {
		case *StringValue:
			enums = []string{first.String()}
		case *ArrayValue:
			return NewEnumType2(first.(*ArrayValue).Elements()...)
		default:
			panic(NewIllegalArgumentType2(`Enum[]`, 0, `String or Array[String]`, args[0]))
		}
	} else {
		enums = make([]string, top)
		for idx, arg := range args {
			str, ok := arg.(*StringValue)
			if !ok {
				panic(NewIllegalArgumentType2(`Enum[]`, idx, `String`, arg))
			}
			enums[idx] = str.String()
		}
	}
	return NewEnumType(enums)
}

func (t *EnumType) Equals(o interface{}, g Guard) bool {
	if ot, ok := o.(*EnumType); ok {
		return len(t.values) == len(ot.values) && ContainsAllStrings(t.values, ot.values)
	}
	return false
}

func (t *EnumType) Generic() PType {
	return enumType_DEFAULT
}

func (t *EnumType) IsAssignable(o PType, g Guard) bool {
	if len(t.values) == 0 {
		switch o.(type) {
		case *StringType, *EnumType, *PatternType:
			return true
		}
		return false
	}

	if st, ok := o.(*StringType); ok {
		return ContainsString(t.values, st.value)
	}

	if en, ok := o.(*EnumType); ok {
		oEnums := en.values
		return len(oEnums) > 0 && ContainsAllStrings(t.values, oEnums)
	}
	return false
}

func (t *EnumType) IsInstance(o PValue, g Guard) bool {
	str, ok := o.(*StringValue)
	return ok && (len(t.values) == 0 || ContainsString(t.values, str.String()))
}

func (t *EnumType) Name() string {
	return `Enum`
}

func (t *EnumType) String() string {
	return ToString2(t, NONE)
}

func (t *EnumType) Parameters() []PValue {
	top := len(t.values)
	if top == 0 {
		return EMPTY_VALUES
	}
	v := make([]PValue, top)
	for idx, e := range t.values {
		v[idx] = WrapString(e)
	}
	return v
}

func (t *EnumType) ToString(b Writer, f FormatContext, g RDetect) {
	TypeToString(t, b, f, g)
}

func (t *EnumType) Type() PType {
	return &TypeType{t}
}

var enumType_DEFAULT = &EnumType{[]string{}}