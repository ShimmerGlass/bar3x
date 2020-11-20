package markup

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"text/scanner"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/pkg/errors"
	"github.com/shimmerglass/bar3x/ui"
	log "github.com/sirupsen/logrus"
)

func (m *Markup) evaluateCtxAttr(expr string) (gval.Evaluable, error) {
	if len(expr) == 0 {
		return func(c context.Context, parameter interface{}) (interface{}, error) {
			return "", nil
		}, nil
	}

	// literal
	if expr[0] != '{' || expr[len(expr)-1] != '}' {
		eval, err := literalExpr(reflect.TypeOf(""), expr)
		if err != nil {
			return nil, err
		}
		return eval, nil
	}

	expr = expr[1 : len(expr)-1]

	return gval.NewLanguage(gval.Full()).NewEvaluable(expr)
}

func (m *Markup) evaluateAttr(elem ui.Drawable, targetFieldName, state string, expr string, vars interface{}, refs reflect.Value, refsReady map[string]*sync.WaitGroup) (*property, error) {
	targetField, err := newField(reflect.ValueOf(elem), targetFieldName)
	if err != nil {
		return nil, err
	}

	var stateMask elementState
	switch state {
	case "Hover":
		stateMask = statePointerHover
	default:
		stateMask = stateNone
	}

	if len(expr) == 0 {
		return &property{
			state: stateMask,
			field: targetField,
			expr: func(c context.Context, parameter interface{}) (interface{}, error) {
				return "", nil
			},
		}, nil
	}

	// literal
	if expr[0] != '{' || expr[len(expr)-1] != '}' {
		eval, err := literalExpr(targetField.Type, expr)
		if err != nil {
			return nil, errors.Wrapf(err, "field %s", targetFieldName)
		}
		return &property{
			state: stateMask,
			field: targetField,
			expr:  eval,
		}, nil
	}

	expr = expr[1 : len(expr)-1]

	var refsStateLock sync.Mutex
	refsState := map[string]interface{}{}

	// expression
	var eval gval.Evaluable
	var evalReady sync.WaitGroup
	evalReady.Add(1)
	lang := gval.NewLanguage(gval.Full(),
		gval.PrefixExtension('$', func(c context.Context, p *gval.Parser) (gval.Evaluable, error) {
			r := p.Scan()
			if r != scanner.Ident {
				return nil, fmt.Errorf("expected ref name")
			}
			refName := p.TokenText()
			r = p.Scan()
			if r != '.' {
				return nil, fmt.Errorf("expected dot accessor after ref name")
			}
			r = p.Scan()
			if r != scanner.Ident {
				return nil, fmt.Errorf("expected ref field name")
			}
			fieldName := p.TokenText()

			ref := refs.FieldByName(refName)
			field, err := newField(ref, fieldName)
			if err != nil {
				return nil, err
			}

			refStateName := fmt.Sprintf("%s.%s", refName, fieldName)

			go func() {
				evalReady.Wait()
				refsReady[refName].Wait()

				err = field.Watch(func(v interface{}) {
					refsStateLock.Lock()
					refsState[refStateName] = v
					refsStateLock.Unlock()
					res, err := eval(c, vars)
					if err != nil {
						log.Fatal(err)
					}

					err = targetField.Set(res)
					if err != nil {
						log.Fatal(err)
					}
				})
				if err != nil {
					log.Fatal(err)
				}
			}()

			return func(c context.Context, parameter interface{}) (interface{}, error) {
				refsStateLock.Lock()
				defer refsStateLock.Unlock()

				v, ok := refsState[refStateName]
				if ok {
					return v, nil
				}
				return reflect.Zero(field.Type).Interface(), nil
			}, nil
		}),
	)

	eval, err = lang.NewEvaluable(expr)
	if err != nil {
		return nil, err
	}
	evalReady.Done()

	return &property{
		state: stateMask,
		field: targetField,
		expr:  eval,
	}, nil
}

func literalExpr(typ reflect.Type, str string) (gval.Evaluable, error) {
	var value interface{}
	switch typ.Kind() {
	case reflect.String:
		value = str

	case reflect.Int:
		v, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		value = v

	case reflect.Int64:
		if typ.String() == "time.Duration" {
			v, err := time.ParseDuration(str)
			if err != nil {
				return nil, err
			}
			value = v
		} else {
			v, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return nil, err
			}
			value = v
		}

	case reflect.Float64:
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		value = v

	case reflect.Bool:
		v, err := strconv.ParseBool(str)
		if err != nil {
			return nil, err
		}
		value = v

	case reflect.Slice:
		elType := typ.Elem()
		switch elType.Kind() {
		case reflect.String:
			value = strings.Split(str, ",")

		default:
			return nil, fmt.Errorf("slice of type %s not handled", elType)
		}

	case reflect.Interface:
		switch typ.String() {
		case "color.Color":
			v, err := ui.ParseColor(str)
			if err != nil {
				return nil, err
			}
			value = v

		default:
			return nil, fmt.Errorf("type %s not handled", typ)
		}

	default:
		return nil, fmt.Errorf("type %s not handled", typ)
	}

	return func(c context.Context, parameter interface{}) (interface{}, error) {
		return value, nil
	}, nil
}
