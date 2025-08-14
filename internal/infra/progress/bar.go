package progress

import (
	progressbar "github.com/schollz/progressbar/v3"
)

// ProgressBarClient interface to abstract progressbar.ProgressBar operations
type ProgressBarClient interface {
	NewOptions64(total int64, options ...progressbar.Option) *progressbar.ProgressBar
	Set64(value int64)
	Finish()
}

// DefaultProgressBarClient implements ProgressBarClient using progressbar.ProgressBar
type DefaultProgressBarClient struct {
	bar *progressbar.ProgressBar
}

func (d *DefaultProgressBarClient) NewOptions64(total int64, options ...progressbar.Option) *progressbar.ProgressBar {
	d.bar = progressbar.NewOptions64(total, options...)
	return d.bar
}

func (d *DefaultProgressBarClient) Set64(value int64) {
	d.bar.Set64(value)
}

func (d *DefaultProgressBarClient) Finish() {
	d.bar.Finish()
}

// TerminalProgressBar struct with injected dependencies
type TerminalProgressBar struct {
	barClient ProgressBarClient
}

func NewTerminalProgressBar(barClient ProgressBarClient) *TerminalProgressBar {
	return &TerminalProgressBar{barClient: barClient}
}

func (tp *TerminalProgressBar) Start(total int64) {
	tp.barClient.NewOptions64(total,
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
	)
}

func (tp *TerminalProgressBar) Update(current int64) {
	tp.barClient.Set64(current)
}

func (tp *TerminalProgressBar) Finish() {
	tp.barClient.Finish()
}
