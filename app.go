package main

import (
	"io"
	"os"

	"github.com/AlecAivazis/survey/v2"
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

func (a *app) Start() error {
	for {
		_, err := a.selectAddr()
		if err != nil {
			return err
		}
	}
}

func (a *app) selectAddr() (addr string, err error) {
	var target string
	addr = a.opts.addr
	err = survey.AskOne(&survey.Input{
		Message:  "Choose an address:",
		Suggest: func(toComplete string) []string {
			return []string{defaultAddr}
		},
	}, &target, survey.WithValidator(survey.Required), surveyIcons())
	if err != nil {
		return "", err
	}
	return
}

func surveyIcons() survey.AskOpt {
	return survey.WithIcons(func(icons *survey.IconSet) {
		icons.SelectFocus.Text = "→"
	})
}