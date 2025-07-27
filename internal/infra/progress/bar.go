package progress

import (
    "github.com/schollz/progressbar/v3"
    "downloader/internal/domain"
)

type TerminalProgressBar struct {
    bar *progressbar.ProgressBar
}

func NewTerminalProgressBar() *TerminalProgressBar {
    return &TerminalProgressBar{}
}

func (tp *TerminalProgressBar) Start(total int64) {
    tp.bar = progressbar.NewOptions64(total,
        progressbar.OptionSetDescription("Downloading"),
        progressbar.OptionShowBytes(true),
        progressbar.OptionSetWidth(40),
    )
}

func (tp *TerminalProgressBar) Update(current int64) {
    tp.bar.Set64(current)
}

func (tp *TerminalProgressBar) Finish() {
    tp.bar.Finish()
}
