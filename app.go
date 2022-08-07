package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/peterh/liner"
)

const (
	Sender     = "send"
	Receiver   = "receive"
	BackSelect = ".."
	FileConfig = "file config"
)

var ErrBackSelect = errors.New("go back")

type app struct {
	opts       *startOpts
	w          io.Writer
	hasService bool
	isSender   bool
	isReceiver bool
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

func (a *app) Start() (err error) {
	if !a.hasService {
		err := a.selectService()
		if err != nil {
			return err
		}
	}

	if a.isSender {
		err = a.sendStart()
	}

	if a.isReceiver {
		err = a.receiveStart()
	}
	return err
}

func (a *app) receiveStart() error {
	r := NewReceive("0.0.0.0:1234")
	r.receive()
	return nil
}

func (a *app) sendStart() error {
	for {
		addr, err := a.selectAddr()
		if err != nil {
			return err
		}
		s := NewSend(addr)
		conn, err := net.Dial("tcp", s.addr)
		if err != nil {
			s.logf("net dial err: %v", err)
			return err
		}
		s.conn = conn

		for {
			paths, err := a.selectFilepath()
			if err != nil {
				if err == ErrBackSelect {
					break
				}
				return err
			}
			s.addPath(paths...)

			ok, err := s.ack()
			if err == nil && ok {
				for _, path := range s.paths {
					s.sendFile(path)
				}
			}
		}
	}
}

func (a *app) selectService() (err error) {
	var target string
	err = survey.AskOne(&survey.Select{
		Message: "Choose a service:",
		Options: []string{Sender, Receiver},
	}, &target, survey.WithValidator(survey.Required), surveyIcons())
	if err != nil {
		return err
	}
	if target == Sender {
		a.isSender = true
	} else if target == Receiver {
		a.isReceiver = true
	}
	return
}

func (a *app) selectAddr() (target string, err error) {
	err = survey.AskOne(&survey.Input{
		Message: "Input a receiver address:",
		Help:    "127.0.0.1:1234",
		Suggest: func(toComplete string) []string {
			return []string{defaultAddr}
		},
	}, &target, survey.WithValidator(survey.Required), surveyIcons())
	if err != nil {
		return "", err
	}
	return
}

func (a *app) selectFilepath() (paths []string, err error) {
	var target string
	err = survey.AskOne(&survey.Input{
		Message: "Input a file path to send",
		Help:    "You can enter the file path and batch write the file path to the configuration file './ftf.conf' separated by commas",
		Suggest: func(toComplete string) []string {
			return []string{BackSelect, FileConfig, "./test.txt"}
		},
	}, &target, survey.WithValidator(survey.Required), surveyIcons())
	if err != nil {
		return nil, err
	}

	if target == FileConfig {
		var f []byte
		f, err = ioutil.ReadFile("./ftf.conf")
		if err != nil {
			return
		}
		paths := strings.Split(string(f), ",")
		return paths, nil
	}

	if target == BackSelect {
		return nil, ErrBackSelect
	}
	return []string{target}, nil
}

func (a *app) predict() {
	line := liner.NewLiner()
	defer line.Close()

	var filePath = make([]string, 0)
	line.SetCtrlCAborts(true)

	f, err := os.Open(HistoryFile)
	if err == nil {
		defer func() {
			line.WriteHistory(f)
			f.Close()
		}()

		buf := new(bytes.Buffer)
		line.ReadHistory(f)
		line.WriteHistory(buf)
		history := strings.Split(buf.String(), "\n")
		filterDuplicate(history)
		filePath = append(filePath, history...)
	}

	line.SetCompleter(func(line string) (c []string) {
		for _, n := range filePath {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	prompt := "input file path> "
	for {
		if target, err := line.Prompt(prompt); err == nil {
			if target == "" {
				continue
			}
			if target == "quit" || target == "exit" {
				PrintRed("bye")
				break
			}
			target = strings.TrimSpace(target)
			PrintCyan(fmt.Sprintf("input: %v", target))
			line.AppendHistory(target)
		} else if err == liner.ErrPromptAborted {
			PrintRed("bye")
			break
		} else {
			PrintRed(fmt.Sprintf("Error reading line: %s", err.Error()))
			break
		}
	}
}

func surveyIcons() survey.AskOpt {
	return survey.WithIcons(func(icons *survey.IconSet) {
		icons.SelectFocus.Text = "→"
	})
}
