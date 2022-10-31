package apt

import (
	"apt-explorer/lib/apt/transport"
	"apt-explorer/lib/log"
	"errors"
	"io"
	"os"
	"path"
)

type ArchiveFilesystem struct {
	root string
}

func UseArchiveFilesystem(rootPath string) (ArchiveFilesystem, error) {
	if _, err := os.Stat(rootPath); errors.Is(err, os.ErrNotExist) {
		return ArchiveFilesystem{}, err
	}
	return ArchiveFilesystem{
		root: rootPath,
	}, nil
}

//func (archive ArchiveFilesystem) Distribution(distName string) (Distribution, error) {
//	return FetchDistribution(archive, distName)
//}

//func (archive ArchiveFilesystem) Fetch(vf transport.VerifiedFile) (string, error) {
//	reader, err := archive.Fetch(vf)
//	if err != nil {
//		return "", err
//	}
//	b, err := io.ReadAll(reader)
//	if err != nil {
//		return "", err
//	}
//	return string(b), nil
//}

func (archive ArchiveFilesystem) Fetch(vf transport.VerifiedFile) (io.Reader, error) {
	relativePath := vf.Path
	concretePath := path.Join(archive.root, relativePath)
	f, err := os.Open(concretePath)
	if err != nil {
		log.Error("Cannot retrieve %s from filesystem", concretePath)
		return nil, err
	}
	return f, nil
}

//func (archive ArchiveFilesystem) URL() *url.URL {
//	aPath, _ := filepath.Abs(archive.root)
//	urlStr := "file://" + aPath
//	urlGood, _ := url.Parse(urlStr)
//	return urlGood
//}
