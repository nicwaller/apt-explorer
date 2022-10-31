package transport

import (
	"fmt"
)

type VerifiedFile struct {
	Path      string // relative to the archive root
	Size      uint64 // the largest .debs today are around 1 GB
	Checksums map[ChecksumAlgorithm][]byte
}

type ChecksumAlgorithm uint8

//goland:noinspection GoSnakeCaseUsage
const (
	CHECKSUM_ALGO_MD5    ChecksumAlgorithm = iota
	CHECKSUM_ALGO_SHA1   ChecksumAlgorithm = iota
	CHECKSUM_ALGO_SHA256 ChecksumAlgorithm = iota
	CHECKSUM_ALGO_SHA512 ChecksumAlgorithm = iota
)

var ChecksumType = map[string]ChecksumAlgorithm{
	"MD5Sum": CHECKSUM_ALGO_MD5,
	"SHA1":   CHECKSUM_ALGO_SHA1,
	"SHA256": CHECKSUM_ALGO_SHA256,
	"SHA512": CHECKSUM_ALGO_SHA512,
}

func AlgoToString(algo ChecksumAlgorithm) string {
	switch algo {
	case CHECKSUM_ALGO_MD5:
		return "MD5Sum"
	case CHECKSUM_ALGO_SHA1:
		return "SHA1"
	case CHECKSUM_ALGO_SHA256:
		return "SHA256"
	case CHECKSUM_ALGO_SHA512:
		return "SHA512"
	default:
		panic(fmt.Sprintf("Unrecognized checksum identity %d", algo))
	}
}

func PreferredAlgorithms() []ChecksumAlgorithm {
	return []ChecksumAlgorithm{
		CHECKSUM_ALGO_SHA512,
		CHECKSUM_ALGO_SHA256,
		CHECKSUM_ALGO_SHA1,
		CHECKSUM_ALGO_MD5,
	}
}
