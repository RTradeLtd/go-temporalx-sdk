package client

import "github.com/schollz/progressbar/v2"

type progressTracker struct {
	bar *progressbar.ProgressBar
}

func newPT(maxBytes int64) *progressTracker {
	return &progressTracker{
		bar: progressbar.NewOptions64(
			maxBytes,
			progressbar.OptionSetRenderBlankState(true),
		),
	}
}

func (p *progressTracker) Update(sent int) {
	p.bar.Add(sent)
}
