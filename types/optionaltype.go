package types

import (
	"io"

	"github.com/lyraproj/puppet-evaluator/errors"
	"github.com/lyraproj/puppet-evaluator/eval"
	"reflect"
)

type OptionalType struct {
	typ eval.Type
}

var OptionalMetaType eval.ObjectType

func init() {
	OptionalMetaType = newObjectType(`Pcore::OptionalType`,
		`Pcore::AnyType {
	attributes => {
		type => {
			type => Optional[Type],
			value => Any
		},
	}
}`, func(ctx eval.Context, args []eval.Value) eval.Value {
			return NewOptionalType2(args...)
		})
}

func DefaultOptionalType() *OptionalType {
	return optionalType_DEFAULT
}

func NewOptionalType(containedType eval.Type) *OptionalType {
	if containedType == nil || containedType == anyType_DEFAULT {
		return DefaultOptionalType()
	}
	return &OptionalType{containedType}
}

func NewOptionalType2(args ...eval.Value) *OptionalType {
	switch len(args) {
	case 0:
		return DefaultOptionalType()
	case 1:
		if containedType, ok := args[0].(eval.Type); ok {
			return NewOptionalType(containedType)
		}
		if containedType, ok := args[0].(stringValue); ok {
			return NewOptionalType3(string(containedType))
		}
		panic(NewIllegalArgumentType2(`Optional[]`, 0, `Variant[Type,String]`, args[0]))
	default:
		panic(errors.NewIllegalArgumentCount(`Optional[]`, `0 - 1`, len(args)))
	}
}

func NewOptionalType3(str string) *OptionalType {
	return &OptionalType{NewStringType(nil, str)}
}

func (t *OptionalType) Accept(v eval.Visitor, g eval.Guard) {
	v(t)
	t.typ.Accept(v, g)
}

func (t *OptionalType) ContainedType() eval.Type {
	return t.typ
}

func (t *OptionalType) Default() eval.Type {
	return optionalType_DEFAULT
}

func (t *OptionalType) Equals(o interface{}, g eval.Guard) bool {
	if ot, ok := o.(*OptionalType); ok {
		return t.typ.Equals(ot.typ, g)
	}
	return false
}

func (t *OptionalType) Generic() eval.Type {
	return NewOptionalType(eval.GenericType(t.typ))
}

func (t *OptionalType) Get(key string) (value eval.Value, ok bool) {
	switch key {
	case `type`:
		return t.typ, true
	}
	return nil, false
}

func (t *OptionalType) IsAssignable(o eval.Type, g eval.Guard) bool {
	return GuardedIsAssignable(o, undefType_DEFAULT, g) || GuardedIsAssignable(t.typ, o, g)
}

func (t *OptionalType) IsInstance(o eval.Value, g eval.Guard) bool {
	return o == _UNDEF || GuardedIsInstance(t.typ, o, g)
}

func (t *OptionalType) MetaType() eval.ObjectType {
	return OptionalMetaType
}

func (t *OptionalType) Name() string {
	return `Optional`
}

func (t *OptionalType) Parameters() []eval.Value {
	if t.typ == DefaultAnyType() {
		return eval.EMPTY_VALUES
	}
	if str, ok := t.typ.(*vcStringType); ok && str.value != `` {
		return []eval.Value{stringValue(str.value)}
	}
	return []eval.Value{t.typ}
}

func (t *OptionalType) ReflectType(c eval.Context) (reflect.Type, bool) {
	return ReflectType(c, t.typ)
}

func (t *OptionalType) Resolve(c eval.Context) eval.Type {
	t.typ = resolve(c, t.typ)
	return t
}

func (t *OptionalType) CanSerializeAsString() bool {
	return canSerializeAsString(t.typ)
}

func (t *OptionalType) SerializationString() string {
	return t.String()
}

func (t *OptionalType) String() string {
	return eval.ToString2(t, NONE)
}

func (t *OptionalType) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	TypeToString(t, b, s, g)
}

func (t *OptionalType) PType() eval.Type {
	return &TypeType{t}
}

var optionalType_DEFAULT = &OptionalType{typ: anyType_DEFAULT}
