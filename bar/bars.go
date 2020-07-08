package bar

import (
	"fmt"

	"github.com/BurntSushi/xgbutil"
	"github.com/shimmerglass/bar3x/tray"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	"github.com/shimmerglass/bar3x/ui/module"
	"github.com/shimmerglass/bar3x/x"
)

type Bars struct {
	x    *xgbutil.XUtil
	ctx  ui.Context
	bars []*Bar

	mk    *markup.Markup
	clock *module.Clock

	LeftRoot   *ui.Root
	CenterRoot *ui.Root
	RightRoot  *ui.Root
}

func CreateBars(ctx ui.Context, x *xgbutil.XUtil) (*Bars, error) {
	bars := &Bars{
		x:   x,
		ctx: ctx,
	}

	clock := module.NewClock(bars.onClockTick)
	mk := markup.New()
	base.RegisterMarkup(mk)
	module.RegisterMarkup(mk, clock)

	bars.mk = mk
	bars.clock = clock

	err := bars.createRoots()
	if err != nil {
		return nil, err
	}

	err = bars.createBars()
	if err != nil {
		return nil, err
	}

	go clock.Run()

	return bars, nil
}

func (b *Bars) createBars() error {
	screens, err := x.Screens(b.x.Conn())
	if err != nil {
		return err
	}

	var trayOutput string
	if b.ctx.Has("tray_output") {
		trayOutput = b.ctx.MustString("tray_output")
	} else {
		trayOutput = screens[0].Outputs[0]
	}
	trayCreated := false

	for _, s := range screens {
		bar, err := NewBar(b.ctx, b.x, s, b.LeftRoot, b.CenterRoot, b.RightRoot)
		if err != nil {
			return err
		}
		b.bars = append(b.bars, bar)

		if trayCreated {
			continue
		}

		for _, output := range s.Outputs {
			if output != trayOutput {
				continue
			}
			err = b.createTray(bar, s)
			if err != nil {
				return err
			}
			trayCreated = true
		}
	}

	return nil
}

func (b *Bars) createRoots() error {
	if b.ctx.Has("bar_left") {
		ctx := b.ctx.New(ui.Context{"bar_align": "left"})
		b.LeftRoot = ui.NewRoot(ctx, func() {
			b.LeftRoot.Paint()
			for _, bar := range b.bars {
				bar.PaintLeft(b.LeftRoot.Image())
			}
		})
	}

	if b.ctx.Has("bar_center") {
		ctx := b.ctx.New(ui.Context{"bar_align": "center"})
		b.CenterRoot = ui.NewRoot(ctx, func() {
			b.CenterRoot.Paint()
			for _, bar := range b.bars {
				bar.PaintCenter(b.CenterRoot.Image())
			}
		})
	}

	if b.ctx.Has("bar_right") {
		ctx := b.ctx.New(ui.Context{"bar_align": "right"})
		b.RightRoot = ui.NewRoot(ctx, func() {
			b.RightRoot.Paint()
			for _, bar := range b.bars {
				bar.PaintRight(b.RightRoot.Image())
			}
		})
	}

	modules, err := b.mk.Parse(b.LeftRoot, nil, b.ctx.MustString("bar_left"))
	if err != nil {
		return fmt.Errorf("config: bar_left: %w", err)
	}
	b.LeftRoot.Inner = modules

	modules, err = b.mk.Parse(b.CenterRoot, nil, b.ctx.MustString("bar_center"))
	if err != nil {
		return fmt.Errorf("config: bar_center: %w", err)
	}
	b.CenterRoot.Inner = modules

	modules, err = b.mk.Parse(b.RightRoot, nil, b.ctx.MustString("bar_right"))
	if err != nil {
		return fmt.Errorf("config: bar_right: %w", err)
	}
	b.RightRoot.Inner = modules

	return nil
}

func (b *Bars) createTray(bar *Bar, screen x.Screen) error {
	tr := tray.New(
		b.x,
		b.bars[0].Win,
		tray.Config{
			BarWidth:        screen.Width,
			BarHeight:       b.ctx.MustInt("bar_height"),
			IconSize:        b.ctx.MustInt("tray_icon_size"),
			IconPadding:     b.ctx.MustInt("tray_icon_padding"),
			BackgroundColor: b.ctx.MustColor("bg_color"),
		},
		func(s tray.State) {
			b.bars[0].SetTrayWidth(s.Width)
			b.RightRoot.Notify()
		},
	)
	return tr.Init()
}

func (b *Bars) onClockTick() {
	b.LeftRoot.Paint()
	b.CenterRoot.Paint()
	b.RightRoot.Paint()

	for _, bar := range b.bars {
		bar.PaintLeft(b.LeftRoot.Image())
		bar.PaintCenter(b.CenterRoot.Image())
		bar.PaintRight(b.RightRoot.Image())
	}
}
