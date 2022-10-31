package transport

import (
	"apt-explorer/lib/log"
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"time"
)

const cachePath = "/tmp/apt-explorer/"

func init() {
	err := os.MkdirAll(cachePath, 0755)
	if err != nil {
		log.Fatal(true, "Cannot create cache directory")
		return
	}
}

type FetchCache struct {
	originalFetcher Fetcher
}

//goland:noinspection GoUnusedExportedFunction
func UseCache(f Fetcher) Fetcher {
	return FetchCache{
		originalFetcher: f,
	}
}

//func (cache FetchCache) Fetch(f VerifiedFile) (string, error) {
//	reader, err := cache.Fetch(f)
//	if err != nil {
//		return "", err
//	}
//
//	b, err := io.ReadAll(reader)
//	if err != nil {
//		return "", err
//	}
//
//	return string(b), nil
//}

func (cache FetchCache) Fetch(f VerifiedFile) (io.Reader, error) {
	// First, look in the well-known cache places
	//log.Debug("Looking for %s in cache", f.Path)
	if len(f.Checksums) == 0 {
		log.Debug(" No checksums provided for %s; will need to go by name or path", f.Path)
		return cache.originalFetcher.Fetch(f)
	}
	for _, algo := range PreferredAlgorithms() {
		if hash, found := f.Checksums[algo]; found {
			expectedPath := path.Join(cachePath, "by-hash", AlgoToString(algo), hex.EncodeToString(hash))
			_, err := os.Stat(expectedPath)
			if errors.Is(err, fs.ErrNotExist) {
				//log.Debug(" Not found: %s", expectedPath)
				continue
			} else {
				//log.Debug(" Found cached file: %s", expectedPath)
				cacheReader, err := os.Open(expectedPath)
				// FIXME: do I need to close the file handle with a defer? if not here, then where?
				if err != nil {
					return cacheReader, nil
				}
			}
		}
	}
	//log.Debug("No cached version available for \"%s\" (%d hashes)", f.Path, len(f.Checksums))

	if IsInNegativeCache(f.Path) {
		log.Warning("Skipping " + f.Path + " (negative cache hit)")
		return nil, errors.New("Negative cache hit for " + f.Path)
	}

	reader, err := cache.originalFetcher.Fetch(f)
	if err != nil {
		AddToNegativeCache(f.Path)
		var none io.Reader
		return none, err
	}

	// TODO: make clever use of a TeeReader to avoid delays with the read-through cache
	// TODO: all the hashes, simultaneously!
	// with symlinks

	// TODO: don't create this dir every. single. time.
	_ = os.MkdirAll("/tmp/apt-explorer/download", 0755)

	// Download the file to a temporary location because we don't know the real hashes yet
	file, err := os.CreateTemp("/tmp/apt-explorer/download", "")
	if err != nil {
		log.Error("Failed creating temporary file")
		var none io.Reader
		return none, err
	}
	size, err := io.Copy(file, reader)
	_ = file.Close()
	if err != nil {
		log.Error("Failed downloading into temp file for cache")
		log.Error("%v", err)
		var none io.Reader
		return none, err
	}
	_ = size
	//log.Info("Downloaded file of %d bytes", size)

	vff, err := addFileToCache(file.Name(), f.Checksums, false) // TODO: yes, unlink
	if err != nil {
		log.Error("%v", err)
		log.Error("Failed adding to cache")
		time.Sleep(100 * time.Millisecond) // TODO: do not pause here
		var none io.Reader
		return none, err
	}
	// TODO: open one of the new cache files instead
	_ = vff
	finalName := file.Name()
	finalReader, err := os.Open(finalName)
	if err != nil {
		log.Error("Failed re-opening cache file for streaming")
		var none io.Reader
		return none, err
	}
	return finalReader, nil
}

// TODO: create symlinks or hardlinks?
func addFileToCache(fullpath string, knownHashes map[ChecksumAlgorithm][]byte, unlinkOriginal bool) (VerifiedFile, error) {
	//log.Debug("trying to add file to cache: %s", fullpath)
	x := hex.EncodeToString(knownHashes[CHECKSUM_ALGO_SHA256])
	log.Debug("%s", x)
	rdr, err := os.Open(fullpath)
	if err != nil {
		return VerifiedFile{}, err
	}
	defer func(rdr *os.File) {
		_ = rdr.Close()
	}(rdr)
	hashes := CalculateBasicHashes(rdr)
	for _, algo := range PreferredAlgorithms() {
		if hash, found := hashes.Checksums[algo]; found {
			// FIXME: !!i really need to handle all the hashes, otherwise it will re-download always!!
			// TODO: also compare file size as a quick check for faster failure

			// Compare against known hashes
			if known, ok := knownHashes[algo]; ok {
				if 0 != bytes.Compare(hash, known) {
					log.Debug("Hash mismatch (type=%v) on downloaded file. This commonly happens with 404 pages.\n Expected: %v\n Calculated: %v", AlgoToString(algo), known, hash)
					return VerifiedFile{}, errors.New("hash Mismatch")
				}
			}

			// TODO: refactor this function for uniformity
			dirPath := path.Join(cachePath, "by-hash", AlgoToString(algo))
			err2 := os.MkdirAll(dirPath, 0755)
			if err2 != nil {
				log.Error("%v", err2)
			}
			expectedPath := path.Join(cachePath, "by-hash", AlgoToString(algo), hex.EncodeToString(hash))
			_, err := os.Stat(expectedPath)
			if err == nil {
				log.Debug("Seems like this file already exists. OK. %v", expectedPath)
				// TODO: am I supposed to get the hashes now? ugh.
				return VerifiedFile{}, nil
			}
			linkErr := os.Link(fullpath, expectedPath)
			if linkErr != nil {
				// FIXME: the file may exist here
				// TODO:
				return hashes, linkErr
			}
			//log.Debug("Added cache file %s", expectedPath)
		}
	}

	if unlinkOriginal {
		log.Debug("Will attempt to delete temp file: %s", fullpath)
		_ = os.Remove(fullpath)
	}

	return hashes, nil
}

// This function shamelessly copied from
// http://marcio.io/2015/07/calculating-multiple-file-hashes-in-a-single-pass/
// TODO: add unit tests here too. In particular, compare no data result to some data result
func CalculateBasicHashes(rd io.Reader) VerifiedFile {
	md5hash := md5.New()
	sha1hash := sha1.New()
	sha256hash := sha256.New()
	sha512hash := sha512.New()

	// For optimum speed, Getpagesize returns the underlying system's memory page size.
	pagesize := os.Getpagesize()

	// wraps the Reader object into a new buffered reader to read the files in chunks
	// and buffering them for performance.
	reader := bufio.NewReaderSize(rd, pagesize)

	// creates a multiplexer Writer object that will duplicate all write
	// operations when copying data from source into all different hashing algorithms
	// at the same time
	multiWriter := io.MultiWriter(md5hash, sha1hash, sha256hash, sha512hash)

	// Using a buffered reader, this will write to the writer multiplexer
	// so we only traverse through the file once, and can calculate all hashes
	// in a single byte buffered scan pass.
	//
	_, err := io.Copy(multiWriter, reader)
	if err != nil {
		panic(err.Error())
	}

	return VerifiedFile{
		Checksums: map[ChecksumAlgorithm][]byte{
			//CHECKSUM_ALGO_MD5: md5.Sum(nil)[0:16], // TODO: re-enable MD5Sum
			CHECKSUM_ALGO_SHA1:   sha1hash.Sum(nil),
			CHECKSUM_ALGO_SHA256: sha256hash.Sum(nil),
			CHECKSUM_ALGO_SHA512: sha512hash.Sum(nil),
		},
	}
}
