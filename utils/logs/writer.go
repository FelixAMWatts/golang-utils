/*
 * Copyright (C) 2020-2021 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package logs

import (
	"io"
	"sync"
	"time"

	"github.com/rs/zerolog/diode"
)

type WriterWithSource interface {
	io.WriteCloser
	SetSource(source string) error
}

type MultipleWritersWithSource struct {
	mu      sync.RWMutex
	writers []WriterWithSource
}

func (w *MultipleWritersWithSource) GetWriters() ([]WriterWithSource, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.writers, nil
}

func (w *MultipleWritersWithSource) AddWriters(writers ...WriterWithSource) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.writers = append(w.writers, writers...)
	return nil
}
func (w *MultipleWritersWithSource) Write(p []byte) (n int, err error) {
	writers, err := w.GetWriters()
	if err != nil {
		return
	}
	for _, writer := range writers {
		n, _ = writer.Write(p)
	}
	return
}

func (w *MultipleWritersWithSource) SetSource(source string) (err error) {
	writers, err := w.GetWriters()
	if err != nil {
		return
	}
	for _, writer := range writers {
		err = writer.SetSource(source)
	}
	return
}

func (w *MultipleWritersWithSource) Close() (err error) {
	writers, err := w.GetWriters()
	if err != nil {
		return
	}
	for _, writer := range writers {
		err1 := writer.Close()
		if err1 != nil {
			err = err1
		}
	}
	return
}

func CreateMultipleWritersWithSource(writers ...WriterWithSource) (writer *MultipleWritersWithSource, err error) {
	writer = &MultipleWritersWithSource{}
	err = writer.AddWriters(writers...)
	return
}

type DiodeWriter struct {
	WriterWithSource
	diodeWriter io.Writer
	slowWriter  WriterWithSource
}

func (w *DiodeWriter) Write(p []byte) (n int, err error) {
	return w.diodeWriter.Write(p)
}

func (w *DiodeWriter) Close() error {
	err := w.slowWriter.Close()
	if err != nil {
		return err
	}
	if diodeCloser, ok := w.diodeWriter.(io.Closer); ok {
		return diodeCloser.Close()
	}
	return err
}

func (w *DiodeWriter) SetSource(source string) error {
	return w.slowWriter.SetSource(source)
}

func NewDiodeWriterForSlowWriter(slowWriter WriterWithSource, ringBufferSize int, pollInterval time.Duration, droppedMessagesLogger Loggers) WriterWithSource {
	return &DiodeWriter{diodeWriter: diode.NewWriter(slowWriter, ringBufferSize, pollInterval, func(missed int) {
		if droppedMessagesLogger != nil {
			droppedMessagesLogger.LogError("Logger Dropped %d messages", missed)
		}
	}),
		slowWriter: slowWriter,
	}
}
