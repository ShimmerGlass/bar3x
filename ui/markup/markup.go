package markup

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/shimmerglass/bar3x/ui"
)

var DefaultMarkup = New()

type DrawableFn func(ui.ParentDrawable) ui.Drawable

type Markup struct {
	els map[string]DrawableFn
}

func New() *Markup {
	return &Markup{
		els: map[string]DrawableFn{},
	}
}

func (m *Markup) Register(name string, ctor DrawableFn) {
	m.els[name] = ctor
}

func (m *Markup) Parse(
	parent ui.ParentDrawable,
	refs interface{},
	xmlStr string,
) (ui.Drawable, error) {
	buf := strings.NewReader(xmlStr)

	rootNode := node{}
	err := xml.NewDecoder(buf).Decode(&rootNode)
	if err != nil {
		return nil, err
	}

	refsVal := reflect.ValueOf(refs)
	if refsVal.Kind() == reflect.Ptr {
		refsVal = refsVal.Elem()
	}
	refsReady := map[string]*sync.WaitGroup{}
	m.setupRefs(refsReady, rootNode)
	return m.parseNode(0, parent, refsVal, refsReady, rootNode)
}

func (m *Markup) MustParse(
	parent ui.ParentDrawable,
	refs interface{},
	xmlStr string,
) ui.Drawable {
	d, err := m.Parse(parent, refs, xmlStr)
	if err != nil {
		panic(err)
	}

	return d
}

func (m *Markup) parseNode(idx int, parent ui.ParentDrawable, refs reflect.Value, refsReady map[string]*sync.WaitGroup, node node) (ui.Drawable, error) {
	name := node.XMLName.Local
	drawableCtor, ok := m.els[name]
	if !ok {
		return nil, fmt.Errorf("unknown element %s", name)
	}
	drawable := drawableCtor(parent)
	mkDrawable := newMarkupDrawable(parent, drawable)
	ctx := parent.ChildContext(idx)

	for _, a := range node.Attrs {
		if a.Name.Local == "ref" {
			field := refs.FieldByName(a.Value)
			if !field.IsValid() {
				return nil, fmt.Errorf("%s: ref to %q but property does not exists", node.XMLName.Local, a.Value)
			}
			field.Set(reflect.ValueOf(drawable))
			refsReady[a.Value].Done()
			continue
		}
		if a.Name.Space == "ctx" {
			expr, err := m.evaluateCtxAttr(a.Value)
			if err != nil {
				return nil, err
			}

			mkDrawable.ctxProp = append(mkDrawable.ctxProp, &ctxProp{
				name: a.Name.Local,
				expr: expr,
			})
			continue
		} else {
			prop, err := m.evaluateAttr(drawable, a.Name.Local, a.Name.Space, a.Value, ctx, refs, refsReady)
			if err != nil {
				return nil, errors.Wrap(err, name)
			}

			mkDrawable.properties = append(mkDrawable.properties, prop)
		}
	}

	body := strings.TrimSpace(node.Body)
	if len(body) > 0 {
		prop, err := m.evaluateAttr(drawable, "Text", "", body, ctx, refs, refsReady)
		if err != nil {
			return nil, errors.Wrap(err, name)
		}

		mkDrawable.properties = append(mkDrawable.properties, prop)
	}

	mkDrawable.SetContext(ctx)
	err := mkDrawable.Init()
	if err != nil {
		return nil, errors.Wrap(err, name)
	}

	if len(node.Nodes) > 0 {
		parentD, ok := drawable.(ui.ParentDrawable)
		if !ok {
			return nil, fmt.Errorf("%s: children found but drawable is not a parent", name)
		}

		for i, n := range node.Nodes {
			child, err := m.parseNode(i, parentD, refs, refsReady, n)
			if err != nil {
				return nil, err
			}
			parentD.Add(child)
		}
	}

	return mkDrawable, nil
}

func (m *Markup) setupRefs(refsReady map[string]*sync.WaitGroup, node node) {
	for _, a := range node.Attrs {
		if a.Name.Local == "ref" {
			var wg sync.WaitGroup
			wg.Add(1)
			refsReady[a.Value] = &wg
		}
	}

	for _, n := range node.Nodes {
		m.setupRefs(refsReady, n)
	}
}
