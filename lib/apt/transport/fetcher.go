package transport

import (
	"io"
)

type Fetcher interface {
	Fetch(file VerifiedFile) (io.Reader, error)
}
