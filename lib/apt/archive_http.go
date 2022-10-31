package apt

import (
	"apt-explorer/lib/apt/transport"
	"apt-explorer/lib/log"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ArchiveHttp struct {
	root url.URL
	//CacheEnabled bool
}

//goland:noinspection GoUnusedExportedFunction
func UseArchiveHttp(archiveRoot string) (ArchiveHttp, error) {
	if !strings.HasSuffix(archiveRoot, "/") {
		archiveRoot = archiveRoot + "/"
	}
	parsedUrl, err := url.Parse(archiveRoot)
	if err != nil {
		return ArchiveHttp{}, err
	}
	return ArchiveHttp{
		root: *parsedUrl,
	}, nil
}

//func (archive ArchiveHttp) Distribution(distName string) (Distribution, error) {
//	return FetchDistribution(archive, distName)
//}

//func (archive ArchiveHttp) Fetch(relativePath string) (string, error) {
//	reader, err := archive.Fetch(relativePath)
//	if err != nil {
//		return "", err
//	}
//	b, err := io.ReadAll(reader)
//	if err != nil {
//		return "", err
//	}
//	return string(b), nil
//}

func (archive ArchiveHttp) Fetch(vf transport.VerifiedFile) (io.Reader, error) {
	// FIXME: return error if requesting the same URL too often
	// TODO: be more careful about joining slashes here
	concreteUrl := archive.root.String() + vf.Path
	//log.Debug("HTTP Fetch: %v", concreteUrl)
	resp, err := http.Get(concreteUrl)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		log.Info("[200 OK] %s", concreteUrl)
	case http.StatusNotFound:
		// Fuck. It's pretty common to have 404 results. How long should we wait between retrying 404? We need a negative cache with expiry. I should almost be using Redis at this point.
		log.Warning("404 Not Found: %s", concreteUrl)
		return nil, errors.New(resp.Status)
	default:
		log.Error("Failed HTTP fetch: %v (%s)", resp.Status, concreteUrl)
		time.Sleep(10 * time.Second)
		return nil, errors.New(resp.Status)
	}
	return io.Reader(resp.Body), err
}

//func (archive ArchiveHttp) URL() *url.URL {
//	return &archive.root
//}
