//go:generate go-bindata resources/...
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"github.com/shimmerglass/bar3x/ui"
	"gopkg.in/yaml.v2"
)

var defaultCtx ui.Context

func mustB64Asset(path string) string {
	data, err := Asset(path)
	if err != nil {
		log.Fatalf("could not open embeded asset %s: %s", path, err)
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	return "base64:" + b64
}

func init() {
	defaultCtx = ui.Context{
		"tray_output":       "DVI-I-1", // TODO: fix default
		"tray_icon_size":    20,
		"tray_icon_padding": 2,

		"h_padding":  5,
		"v_padding":  5,
		"bar_height": 30,

		"text_font_size": 13.0,
		"icon_font_size": 13.0,
		"text_font":      mustB64Asset("resources/fonts/noto-sans.ttf"),
		"icon_font":      mustB64Asset("resources/fonts/nerdfont-noto-mono.ttf"),

		"bg_color":            "#17191e",
		"text_color":          "#d4e5f7",
		"accent_color":        "#1ebce8",
		"neutral_color":       "#37393e",
		"neutral_light_color": "#90949d",

		"icons": map[string]interface{}{
			"error":    "\uf071",
			"dot":      "\uf444",
			"transfer": "\ufa4e",
			"chip":     "\uf85a",
			"chip2":    "\uf2db",
			"lock":     "\uf023",
			"calendar": "\uf073",
			"disk":     "\uf0a0",
		},

		"bar_left": `
			<ModuleRow>
				<Volume />
			</ModuleRow>
		`,

		"bar_center": `
			<ModuleRow>
				<DateTime />
			</ModuleRow>
		`,

		"bar_right": `
			<ModuleRow>
				<VPN />
				<CPU />
				<RAM />
				<DiskUsage />
			</ModuleRow>
		`,

		"bar_background": `
			<Rect
				Width="{bar_width}"
				Height="{bar_height}"
				Color="{bg_color}"
			/>
		`,

		"module": `
			<Row ctx:mfirst="{is_first_visible}">
				<Sizer
					Visible="{!mfirst}"
					PaddingLeft="10"
					PaddingRight="10"
				>
				<Rect
					Width="5"
					Height="5"
					Color="{neutral_color}"
				/>
				</Sizer>
				<Sizer ref="Content" />
			</Row>
		`,
	}
}

func getConfig(cfgPath, themePath string) (ui.Context, error) {
	ctx := defaultCtx
	if themePath != "" {
		tctx, err := loadConfigFile(themePath)
		if err != nil {
			return nil, err
		}
		ctx = ctx.New(tctx)
	}
	if cfgPath != "" {
		tctx, err := loadConfigFile(cfgPath)
		if err != nil {
			return nil, err
		}
		ctx = ctx.New(tctx)
	}

	return ctx, nil
}

func loadConfigFile(path string) (ui.Context, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	ctx := ui.Context{}
	err = yaml.Unmarshal(contents, &ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}

	return ctx, nil
}
