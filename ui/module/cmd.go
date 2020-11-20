package module

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
)

type Cmd struct {
	moduleBase

	clock *Clock
	mk    *markup.Markup

	interval time.Duration
	command  string
	format   string

	width  base.WatchInt
	height base.WatchInt

	errorRoot ui.Drawable
	ErrorTxt  *base.Text
}

func NewCmd(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *Cmd {
	return &Cmd{
		mk:         mk,
		clock:      clock,
		moduleBase: newBase(p),

		format:   "markup",
		interval: time.Second,
	}
}

func (c *Cmd) Interval() time.Duration {
	return c.interval
}

func (c *Cmd) SetInterval(v time.Duration) {
	c.interval = v
}

func (c *Cmd) Text() string {
	return c.command
}

func (c *Cmd) SetText(v string) {
	c.command = v
}

func (c *Cmd) Format() string {
	return c.format
}

func (c *Cmd) SetFormat(v string) {
	c.format = v
}

func (c *Cmd) Width() int {
	return c.width.V
}
func (c *Cmd) OnWidthChange(cb func(int)) {
	c.width.Add(cb)
}
func (c *Cmd) Height() int {
	return c.height.V
}
func (c *Cmd) OnHeightChange(cb func(int)) {
	c.height.Add(cb)
}

func (c *Cmd) Init() error {
	root, err := c.mk.Parse(c, c, `
		<Row>
			<Sizer PaddingRight="{h_padding}">
				<Icon>{icons.error}</Icon>
			</Sizer>
			<Text ref="ErrorTxt" />
		</Row>
	`)
	if err != nil {
		return err
	}
	c.errorRoot = root

	c.clock.Add(c, c.interval)
	return nil
}

func (c *Cmd) Update(ctx context.Context) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", c.command)
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	errStr := stderr.String()

	err := cmd.Run()
	if err != nil {
		var msg string
		if errStr != "" {
			msg = errStr
		} else {
			msg = err.Error()
		}
		c.showError(msg)
		return
	}

	var root ui.Drawable
	switch c.format {
	case "markup":
		root, err = c.mk.Parse(c, nil, stdout.String())
		if err != nil {
			c.showError(err.Error())
			return
		}
	case "plain":
		txt := base.NewText(c)
		txt.SetContext(c.Context())
		txt.SetText(strings.TrimSpace(stdout.String()))
		root = txt
	default:
		log.Fatalf("cmd: unknown format %q", c.format)
	}

	c.Root = root
	c.width.Set(root.Width())
	c.height.Set(root.Height())
}

func (c *Cmd) showError(msg string) {
	msg = strings.ReplaceAll(msg, "\n", " ")
	if len(msg) > 30 {
		msg = msg[:30] + "…"
	}

	c.ErrorTxt.SetText(msg)
	c.Root = c.errorRoot
	c.width.Set(c.errorRoot.Width())
	c.height.Set(c.errorRoot.Height())
}
