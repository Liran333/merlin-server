package domain

import (
	"context"
)

type SeverSentStream struct {
	Parameter   StreamParameter
	Ctx         context.Context
	StreamWrite func(doOnce func() ([]byte, error))
}

type StreamParameter struct {
	Token     string `json:"token"`
	StreamUrl string `json:"stream_url" required:"true"`
}

type SeverSentEvent interface {
	Request(*SeverSentStream) error
}
