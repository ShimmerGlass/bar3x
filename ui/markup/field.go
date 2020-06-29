package markup

import (
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/shimmerglass/bar3x/ui"
)

type field struct {
	StructName     string
	Name           string
	Type           reflect.Type
	getter         reflect.Value
	setter         reflect.Value
	onChange       reflect.Value
	onChangeCbType reflect.Type
}

func newField(base reflect.Value, name string) (*field, error) {
	f := &field{
		Name:       name,
		StructName: base.Type().String(),
	}

	getter := base.MethodByName(name)
	if !getter.IsValid() {
		return nil, fmt.Errorf("type %q had no field %q", f.StructName, name)
	}

	if getter.Type().NumOut() != 1 {
		return nil, fmt.Errorf("type %q field %q getter should return 1 value", f.StructName, name)
	}
	getterType := getter.Type().Out(0)
	f.Type = getterType
	f.getter = getter

	setter := base.MethodByName("Set" + name)
	if setter.IsValid() {
		if setter.Type().NumIn() != 1 {
			return nil, fmt.Errorf("%s: field %q setter should accept 1 argument", f.StructName, name)
		}
		setterType := setter.Type().In(0)
		if setterType != getterType {
			return nil, fmt.Errorf("%s: field %q setter should accept a %s", f.StructName, name, getterType)
		}

		f.setter = setter
	}

	onChange := base.MethodByName("On" + name + "Change")
	if onChange.IsValid() {
		if onChange.Type().NumIn() != 1 {
			return nil, fmt.Errorf("%s: field %q onChange should accept 1 argument", f.StructName, name)
		}
		onChangeCb := onChange.Type().In(0)
		if onChangeCb.Kind() != reflect.Func {
			return nil, fmt.Errorf("%s: field %q onChange should accept a function", f.StructName, name)
		}
		if onChangeCb.NumIn() != 1 {
			return nil, fmt.Errorf("%s: field %q onChange callback should accept 1 argument", f.StructName, name)
		}
		if onChangeCb.In(0) != getterType {
			return nil, fmt.Errorf("%s: field %q onChange callback should accept a %s", f.StructName, name, getterType)
		}
		f.onChange = onChange
		f.onChangeCbType = onChangeCb
	}

	return f, nil
}

func (f *field) Get() interface{} {
	res := f.getter.Call([]reflect.Value{})
	return res[0].Interface()
}

func (f *field) Set(v interface{}) error {
	if !f.setter.IsValid() {
		return fmt.Errorf("%s: field %q cannot be set: it is readonly", f.StructName, f.Name)
	}
	val := reflect.ValueOf(v)

	// color special case
	if f.Type.String() == "color.Color" && val.Kind() == reflect.String {
		c, err := ui.ParseColor(v.(string))
		if err != nil {
			return fmt.Errorf("%s: cannot set field %q: %s", f.StructName, f.Name, err)
		}
		val = reflect.ValueOf(c)
	}

	// duration special case
	if f.Type.String() == "time.Duration" && val.Kind() == reflect.String {
		c, err := time.ParseDuration(v.(string))
		if err != nil {
			return fmt.Errorf("%s: cannot set field %q: %s", f.StructName, f.Name, err)
		}
		val = reflect.ValueOf(c)
	}

	// float special case
	if f.Type.Kind() == reflect.Int && val.Kind() == reflect.Float64 {
		val = reflect.ValueOf(int(math.Round(v.(float64))))
	}

	if !val.Type().AssignableTo(f.Type) {
		return fmt.Errorf(
			"%s: cannot set field %q (%s) with %s, expected %s",
			f.StructName,
			f.Name,
			f.Type,
			val.Type().String(),
			f.Type.String(),
		)
	}
	f.setter.Call([]reflect.Value{val})
	return nil
}

func (f *field) Watch(cb func(interface{})) error {
	if !f.onChange.IsValid() {
		return fmt.Errorf("%s: type %q field %q cannot be watched (no onChange found)", f.StructName, f.Type, f.Name)
	}

	fn := reflect.MakeFunc(f.onChangeCbType, func(args []reflect.Value) []reflect.Value {
		cb(args[0].Interface())
		return nil
	})
	f.onChange.Call([]reflect.Value{fn})

	return nil
}
