package debug

import (
	"fmt"
	"html"
	"reflect"
	"strings"

	"github.com/shimmerglass/bar3x/bar"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/markup"
)

const page = `
<html>
<head>
	<script
		src="https://code.jquery.com/jquery-3.5.1.min.js"
		integrity="sha256-9/aliU8dGd2tb6OSsuzixeV4y/faTqgFtohetphbbj0="
		crossorigin="anonymous"></script>
	<script>
		$(function() {
			$(".children-toggle").click(function() {
				var id = $(this).attr("data-id");
				$("#" + id).toggle();
			});

			$(".details-toggle").click(function() {
				$('.details').hide();
				var id = $(this).attr("data-id");
				$("#" + id).show();
			});
		})
	</script>
	<style>
		html, body {
			margin: 0;
			padding: 0;
			font-family: Arial,"Helvetica Neue",Helvetica,sans-serif;
		}

		.children { display: none; }
		.details { display: none; }

		.element {
			margin-left: 20px;
		}

		body {
			display: flex;
			flex-basis: 100%%;
			flex: 1;
		}

		#tree {
			padding: 15px;
			flex-grow: 1;
		}

		#details {
			max-width: 500px;
			overflow: auto;
		}
	</style>
</head>
<body>
	<div id="tree">%s</div>
	<div id="details">%s</div>
</body>
</html>
`

type element struct {
	name     string
	attrs    map[string]interface{}
	ctx      map[string]interface{}
	children []element
}

func (e element) OpenHTML() string {
	b := strings.Builder{}
	b.WriteByte('<')
	b.WriteString(e.name)
	b.WriteByte('>')

	return html.EscapeString(b.String())
}

func (e element) CloseHTML() string {
	return html.EscapeString(fmt.Sprintf("</%s>", e.name))
}

const indentOne = "  "

func createView(b *bar.Bars) string {
	root := element{
		name: "Root",
	}

	for i, bar := range b.Bars {
		root.children = append(root.children, element{
			name: fmt.Sprint(i),
			children: []element{
				{name: "Left", children: createDrawableEl(bar.LeftRoot).children},
				{name: "Center", children: createDrawableEl(bar.CenterRoot).children},
				{name: "Right", children: createDrawableEl(bar.RightRoot).children},
			},
		})
	}

	treeB := &strings.Builder{}
	propsB := &strings.Builder{}
	createDrawableView(treeB, propsB, 0, root)

	return fmt.Sprintf(page, treeB.String(), propsB.String())
}

func createDrawableView(treeB, propsB *strings.Builder, id int, el element) int {
	id++

	propsB.WriteString(fmt.Sprintf(`<div class="details" id="details-%d">`, id))
	propsB.WriteString("<h3>Attributes</h3><table>")
	for k, v := range el.attrs {
		propsB.WriteString(fmt.Sprintf(`
			<tr>
				<td><b>%s</b></td>
				<td>%s</td>
			</tr>
		`, html.EscapeString(k), html.EscapeString(propValue(fmt.Sprint(v)))))
	}
	propsB.WriteString("</table>")
	propsB.WriteString("<h3>Context</h3><table>")
	for k, v := range el.ctx {
		propsB.WriteString(fmt.Sprintf(`
			<tr>
				<td><b>%s</b></td>
				<td>%s</td>
			</tr>
		`, html.EscapeString(k), html.EscapeString(propValue(fmt.Sprint(v)))))
	}
	propsB.WriteString("</table></div>")

	treeB.WriteString(`<div class="element">`)
	if len(el.children) > 0 {
		treeB.WriteString(fmt.Sprintf(`<button class="children-toggle" data-id="el-%d">▸</button>`, id))
	}
	treeB.WriteString(fmt.Sprintf(`<span class="details-toggle" data-id="details-%d">`, id))
	treeB.WriteString(el.OpenHTML())
	treeB.WriteString("</span>")

	treeB.WriteString(fmt.Sprintf(`<div id="el-%d" class="children">`, id))
	for _, c := range el.children {
		id = createDrawableView(treeB, propsB, id, c)
	}
	treeB.WriteString("</div>")

	treeB.WriteString(el.CloseHTML())
	treeB.WriteString(`</div>`)

	return id
}

func createDrawableEl(d ui.Drawable) element {
	if d == nil {
		return element{name: "nil"}
	}
	if m, ok := d.(*markup.MarkupDrawable); ok {
		return createDrawableEl(m.Children()[0])
	}

	el := element{
		name:  strings.TrimPrefix(fmt.Sprintf("%T", d), "*"),
		attrs: map[string]interface{}{},
		ctx:   d.Context(),
	}

	val := reflect.ValueOf(d)
	typ := val.Type()

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if !strings.HasPrefix(method.Name, "Set") {
			continue
		}

		name := strings.TrimPrefix(method.Name, "Set")
		if name == "Context" {
			continue
		}

		getter := val.MethodByName(name)
		if !getter.IsValid() {
			continue
		}
		v := getter.Call(nil)[0]
		if v.IsZero() {
			continue
		}

		el.attrs[name] = v
	}

	for _, c := range d.Children() {
		el.children = append(el.children, createDrawableEl(c))
	}

	return el
}

func propValue(i interface{}) string {
	v := fmt.Sprint(i)
	if len(v) > 128 {
		v = v[:128] + "..."
	}

	return v
}
