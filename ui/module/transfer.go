package module

import (
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
)

type transfer struct {
	moduleBase

	RxTxt   *TextUnit
	TxTxt   *TextUnit
	RxGraph *base.Graph
	TxGraph *base.Graph
	TxtRow  *base.Row

	mk *markup.Markup

	maxRX, maxTX   int
	rxData, txData []int
}

func newTransfer(p ui.ParentDrawable, mk *markup.Markup) *transfer {
	return &transfer{
		mk:         mk,
		moduleBase: newBase(p),
	}
}

func (b *transfer) Init() error {
	_, err := b.mk.Parse(b, b, `
		<Layers ref="Root">
			<Col>
				<Graph
					ref="RxGraph"
					Color="{inactive_color}"
					Height="{height / 2}"
					Width="{$TxtRow.Width}"
					Direction="up"
				/>
				<Graph
					ref="TxGraph"
					Color="{inactive_color}"
					Height="{height / 2}"
					Width="{$TxtRow.Width}"
					Direction="down"
				/>
			</Col>
			<Row ref="TxtRow">
				<Sizer Width="50" HAlign="left">
					<TxtUnit ref="RxTxt" />
				</Sizer>
				<Icon>{icons["transfer"]}</Icon>
				<Sizer Width="50" HAlign="right">
					<TxtUnit ref="TxTxt" />
				</Sizer>
			</Row>
		</Layers>
	`)
	return err
}

func (b *transfer) Set(rx, tx int) {
	b.RxTxt.Set(humanateBytes(uint64(rx)))
	b.TxTxt.Set(humanateBytes(uint64(tx)))

	if rx > b.maxRX {
		b.maxRX = rx
	}
	if tx > b.maxTX {
		b.maxTX = tx
	}

	b.rxData = append(b.rxData, rx)
	b.txData = append(b.txData, tx)

	w := b.RxGraph.Width()
	if len(b.rxData) > w {
		b.rxData = b.rxData[len(b.rxData)-w:]
	}
	if len(b.txData) > w {
		b.txData = b.txData[len(b.txData)-w:]
	}

	rxdr := make([]float64, len(b.rxData))
	for i, v := range b.rxData {
		if b.maxRX > 0 {
			rxdr[i] = float64(v) / float64(b.maxRX)
		} else {
			rxdr[i] = 0
		}
	}

	txdr := make([]float64, len(b.txData))
	for i, v := range b.txData {
		if b.maxTX > 0 {
			txdr[i] = float64(v) / float64(b.maxTX)
		} else {
			txdr[i] = 0
		}
	}

	b.RxGraph.SetData(rxdr)
	b.TxGraph.SetData(txdr)
}
