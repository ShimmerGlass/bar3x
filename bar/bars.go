package bar

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	"github.com/shimmerglass/bar3x/ui/module"
	"github.com/shimmerglass/bar3x/x"
)

type Bars struct {
	x    *xgbutil.XUtil
	ctx  ui.Context
	Bars []*Bar

	mk    *markup.Markup
	clock *module.Clock
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

	err := bars.createBars()
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
		widthTray := false
		if !trayCreated {
			for _, output := range s.Outputs {
				if output != trayOutput {
					continue
				}
				widthTray = true
				trayCreated = true
			}
		}

		bar, err := NewBar(b.ctx, b.x, s, b.mk, widthTray)
		if err != nil {
			return err
		}
		b.Bars = append(b.Bars, bar)
	}

	return nil
}

func (b *Bars) onClockTick() {
	for _, bar := range b.Bars {
		if bar.LeftRoot != nil {
			bar.LeftRoot.Paint()
			bar.PaintLeft(bar.LeftRoot.Image())
		}

		if bar.CenterRoot != nil {
			bar.CenterRoot.Paint()
			bar.PaintCenter(bar.CenterRoot.Image())
		}

		if bar.RightRoot != nil {
			bar.RightRoot.Paint()
			bar.PaintRight(bar.RightRoot.Image())
		}
	}
}
