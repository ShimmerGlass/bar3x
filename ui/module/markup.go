package module

import (
	"github.com/shimmerglass/bar3x/lib/mirror"
	"github.com/shimmerglass/bar3x/lib/process"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/markup"
)

func RegisterMarkup(ctx ui.Context, mk *markup.Markup, clock *Clock) {
	var pr *process.ProcessWatcher
	var mirrorServer *mirror.Server

	mk.Register("TxtUnit", func(p ui.ParentDrawable) ui.Drawable {
		return NewTextUnit(p, mk)
	})

	mk.Register("Transfer", func(p ui.ParentDrawable) ui.Drawable {
		return newTransfer(p, mk)
	})

	mk.Register("Interface", func(p ui.ParentDrawable) ui.Drawable {
		return NewInterface(p, mk, clock)
	})

	mk.Register("CPU", func(p ui.ParentDrawable) ui.Drawable {
		if pr == nil {
			pr = process.New()
		}
		return NewCPU(p, mk, clock, pr)
	})

	mk.Register("DateTime", func(p ui.ParentDrawable) ui.Drawable {
		return NewDateTime(p, mk, clock)
	})

	mk.Register("DiskBandwidth", func(p ui.ParentDrawable) ui.Drawable {
		return NewDiskBW(p, mk, clock)
	})

	mk.Register("DiskUsage", func(p ui.ParentDrawable) ui.Drawable {
		return NewDiskUsage(p, mk, clock)
	})

	mk.Register("Music", func(p ui.ParentDrawable) ui.Drawable {
		return NewMusic(p, mk)
	})

	mk.Register("RAM", func(p ui.ParentDrawable) ui.Drawable {
		if pr == nil {
			pr = process.New()
		}
		return NewRAM(p, mk, clock, pr)
	})

	mk.Register("ModuleRow", func(p ui.ParentDrawable) ui.Drawable {
		return NewSepRow(p, mk)
	})

	mk.Register("Volume", func(p ui.ParentDrawable) ui.Drawable {
		return NewVolume(p, mk)
	})

	mk.Register("Connections", func(p ui.ParentDrawable) ui.Drawable {
		return NewConnections(p, mk, clock)
	})

	mk.Register("Weather", func(p ui.ParentDrawable) ui.Drawable {
		return NewWeather(p, mk, clock)
	})

	mk.Register("Cmd", func(p ui.ParentDrawable) ui.Drawable {
		return NewCmd(p, mk, clock)
	})

	mk.Register("Workspaces", func(p ui.ParentDrawable) ui.Drawable {
		return NewWorkspaces(p, mk)
	})

	mk.Register("Battery", func(p ui.ParentDrawable) ui.Drawable {
		return NewBattery(p, mk, clock)
	})

	mk.Register("MirrorServer", func(p ui.ParentDrawable) ui.Drawable {
		if mirrorServer == nil {
			mirrorServer = mirror.NewServer(ctx.MustString("mirror_server_addr"))
		}
		return NewMirrorServer(p, mirrorServer)
	})

	mk.Register("MirrorClient", func(p ui.ParentDrawable) ui.Drawable {
		return NewMirrorClient(p)
	})
}
