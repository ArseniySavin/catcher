// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/ArseniySavin/catcher/pkg/internal"
)

const (
	ErrorKey = "Error"
)

type CatcherHandler struct {
	opt   *slog.HandlerOptions
	l     *log.Logger
	mu    *sync.Mutex
	attr  []slog.Attr
	group string
}

func NewCatcherHandler(opts *slog.HandlerOptions) slog.Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{Level: slog.LevelInfo}
	}
	return &CatcherHandler{
		opt: opts,
		l:   log.New(os.Stderr, "", log.LstdFlags),
		mu:  &sync.Mutex{}}
}

func (h *CatcherHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opt.Level.Level()
}

func (h *CatcherHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	for i := 0; i < len(attrs); i++ {
		h.attr = append(h.attr, attrs[i])
	}
	return h
}
func (h *CatcherHandler) WithGroup(name string) slog.Handler {
	h.group = name
	return h
}

func (h *CatcherHandler) Handle(_ context.Context, r slog.Record) error {
	msg := internal.LogMsg{
		Level:   r.Level.String(),
		Host:    internal.GetHost(),
		Message: r.Message,
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.opt.AddSource {
		msg.Payload = append(msg.Payload, fmt.Sprintf("Source:%s", internal.CallSource(r.PC).String()))
	}

	if !h.opt.AddSource && r.Level == slog.LevelError {
		msg.Payload = append(msg.Payload, fmt.Sprintf("Source:%s", internal.CallSource(r.PC).String()))
	}

	msg = attr(msg, h.attr, h.group)
	msg = payload(msg, r, h.group)

	h.l.Printf("%s", internal.MarshalStruct(msg))

	return nil
}

func payload(payload internal.LogMsg, r slog.Record, group string) internal.LogMsg {
	if r.NumAttrs() > 0 {
		r.Attrs(func(attr slog.Attr) bool {
			if group != "" {
				payload.Payload = append(payload.Payload, fmt.Sprintf("%s.%s:%s", group, attr.Key, attr.Value.String()))
				return true
			}
			payload.Payload = append(payload.Payload, fmt.Sprintf("%s:%s", attr.Key, attr.Value.String()))
			return true
		})
	}
	return payload
}

func attr(payload internal.LogMsg, attr []slog.Attr, group string) internal.LogMsg {
	if len(attr) > 0 {
		for _, v := range attr {
			if group != "" {
				payload.Payload = append(payload.Payload, fmt.Sprintf("%s.%s:%s", group, v.Key, v.Value.String()))
				continue
			}
			payload.Payload = append(payload.Payload, fmt.Sprintf("%s:%s", v.Key, v.Value.String()))
		}
	}
	return payload
}
