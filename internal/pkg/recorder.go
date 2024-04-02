// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package pkg

import (
	"encoding/json"
	"os"
	"path"
	"time"
)

type Env struct {
	Shell string `json:"SHELL"`
	Term  string `json:"TERM"`
}

type Header struct {
	Title     string `json:"title"`
	Version   int    `json:"version"`
	Height    int    `json:"height"`
	Width     int    `json:"width"`
	Env       Env    `json:"env"`
	Timestamp int    `json:"Timestamp"`
}

type Recorder struct {
	File      *os.File
	Timestamp int
}

func (recorder *Recorder) Close() {
	if recorder.File != nil {
		_ = recorder.File.Close()
	}
}

func (recorder *Recorder) WriteHeader(header *Header) (err error) {
	var p []byte

	if p, err = json.Marshal(header); err != nil {
		return
	}

	if _, err := recorder.File.Write(p); err != nil {
		return err
	}
	if _, err := recorder.File.Write([]byte("\n")); err != nil {
		return err
	}

	recorder.Timestamp = header.Timestamp

	return
}

func (recorder *Recorder) WriteData(data string) (err error) {
	now := int(time.Now().UnixNano())

	delta := float64(now-recorder.Timestamp*1000*1000*1000) / 1000 / 1000 / 1000

	row := make([]interface{}, 0)
	row = append(row, delta)
	row = append(row, "o")
	row = append(row, data)

	var s []byte
	if s, err = json.Marshal(row); err != nil {
		return
	}
	if _, err := recorder.File.Write(s); err != nil {
		return err
	}
	if _, err := recorder.File.Write([]byte("\n")); err != nil {
		return err
	}
	return
}

func NewRecorder(recordingPath, term string, h int, w int) (recorder *Recorder, err error) {
	recorder = &Recorder{}

	if _, err := os.Stat(path.Dir(recordingPath)); err != nil {
		if err := os.MkdirAll(path.Dir(recordingPath), os.ModePerm); err != nil {
			return recorder, err
		}
	}

	file, err := os.Create(recordingPath)
	if err != nil {
		return nil, err
	}

	recorder.File = file

	header := &Header{
		Title:     "",
		Version:   2,
		Height:    h,
		Width:     w,
		Env:       Env{Shell: "/bin/bash", Term: term},
		Timestamp: int(time.Now().Unix()),
	}

	if err := recorder.WriteHeader(header); err != nil {
		return nil, err
	}

	return recorder, nil
}
