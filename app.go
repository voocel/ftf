package main

import (
	"io"
	"os"
)

type app struct {
	opts *startOpts
	w    io.Writer
}

type startOpts struct {
	service string
	addr    string
	path    string
	w       io.Writer
}

func newApp(opts *startOpts) *app {
	a := &app{
		opts: opts,
	}
	a.w = opts.w
	if a.w == nil {
		a.w = os.Stdout
	}

	return a
}
