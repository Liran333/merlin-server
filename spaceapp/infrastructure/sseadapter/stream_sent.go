package sseadapter

import (
	"bufio"
	"io"
	"net/http"

	"github.com/openmerlin/merlin-server/spaceapp/domain"
	"github.com/openmerlin/merlin-server/utils"
)

func StreamSentAdapter() *streamSentAdapter {
	return &streamSentAdapter{utils.NewHttpClient(3, 3600)}
}

type streamSentAdapter struct {
	cli utils.HttpClient
}

func (sse *streamSentAdapter) Request(q *domain.SeverSentStream) error {
	accessToken := q.Parameter.Token

	req, err := http.NewRequestWithContext(q.Ctx, http.MethodGet, q.Parameter.StreamUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("TOKEN", accessToken)

	return sse.cli.SendAndHandle(req, func(h http.Header, respBody io.Reader) error {
		st := streamTransfer{
			input: *bufio.NewReader(respBody),
		}

		q.StreamWrite(st.readAndWriteOnce)

		return nil
	})
}
