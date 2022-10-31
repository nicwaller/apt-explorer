package apt

import (
	"apt-explorer/lib/apt/kvstore"
	"apt-explorer/lib/apt/transport"
	"apt-explorer/lib/log"
	"encoding/hex"
	"fmt"
	"io"
	path2 "path"
	"strconv"
	"strings"
)

// https://wiki.debian.org/DebianRepository/Format#A.22Release.22_files

type Distribution struct {
	Fetcher transport.Fetcher
	Name    string
	Indexes map[string]transport.VerifiedFile
	DistributionMeta
}

// DistributionMeta is information we get from the `Release` file
type DistributionMeta struct {
	// TODO: maybe store these in a map for future-proofing?
	Origin        string   // Optional. Original publisher of the repository
	Label         string   // Optional. Typically used for repos split over multiple CDs
	Suite         string   // Required. Examples: "bionic" (Ubuntu) or "oldstable" (Debian)
	Codename      string   // Required. Examples "bionic" (Ubuntu) or "sid" (Debian)
	Version       string   // Required. Examples "18.04" (Ubuntu)
	Date          string   // Required. Time at which the release file was created.
	ValidUntil    string   // Optional. Clients should consider Release file expired after this date.
	Architectures []string // Required. List of processor architectures supported by binary packages.
	Components    []string // Required. List of components in the distribution. Examples: main, universe
	Description   string   // Required.
	NotAutomatic  bool     // Optional. Suggestion to clients to require human approval before updating.
	AcquireByHash bool     // Optional. Does the repo support content-based addressing of packages?
	SignedBy      []string // Optional. List of OpenPGP key fingerprints to be used for validating the next Release file.
	RequireAuth   bool     // Deprecated. Hint that downloading packages requires authorization.
}

// UseDistribution loads the index of indexes immediately
func UseDistribution(f transport.Fetcher, distName string) (Distribution, error) {
	d := Distribution{
		Fetcher:          f,
		Name:             distName,
		Indexes:          make(map[string]transport.VerifiedFile, 0),
		DistributionMeta: DistributionMeta{},
	}
	_, err := d.Refresh()
	return d, err
}

// Refresh returns true if the Release file changed
func (dist *Distribution) Refresh() (bool, error) {
	log.Debug("Requested Refresh() of distribution \"%v\"", dist.Name)
	// FIXME: make sure to do a global debounce (don't request /Release too often from remote servers)
	distPrefix := strings.Join([]string{"dists", dist.Name}, "/")
	relPath := strings.Join([]string{distPrefix, "Release"}, "/")
	reader, err := dist.Fetcher.Fetch(transport.VerifiedFile{Path: relPath})
	if err != nil {
		return false, err
	}
	oldDate := dist.DistributionMeta.Date
	dist.DistributionMeta, dist.Indexes = ParseReleaseFile(reader, distPrefix)
	// TODO: more rigorous check for changes, with checksum
	if dist.DistributionMeta.Date != oldDate {
		log.Info("Found updated Release file for repo \"%s\"", dist.Description)
		return true, nil
	}
	return false, nil
}

// TODO: write tests for this (I tend to introduce errors in this fn)
func ParseReleaseFile(reader io.Reader, distPrefix string) (DistributionMeta, map[string]transport.VerifiedFile) {
	meta := DistributionMeta{}
	files := make(map[string]transport.VerifiedFile, 0)

	records, err := kvstore.ParseImmediate(reader)
	if err != nil {
		log.Error("Failed parsing release file")
		// TODO: don't panic
		panic(1)
	}

	if len(records) != 1 {
		log.Error("Expected 1 record but got %d", len(records))
		// TODO: don't panic
		panic(2)
	}

	rel := records[0]

	meta.Origin = rel.SingleValues["Origin"]
	meta.Label = rel.SingleValues["Label"]
	meta.Suite = rel.SingleValues["Suite"]
	meta.Version = rel.SingleValues["Version"]
	meta.Codename = rel.SingleValues["Codename"]
	meta.Date = rel.SingleValues["Date"]
	meta.Description = rel.SingleValues["Description"]
	meta.ValidUntil = rel.SingleValues["Valid-Until"]
	meta.Architectures = strings.Fields(rel.SingleValues["Architectures"])
	meta.Components = strings.Fields(rel.SingleValues["Components"])
	b, _ := strconv.ParseBool(rel.SingleValues["Acquire-By-Hash"])
	meta.AcquireByHash = b

	// TODO: handle all the hash types, MD5sum etc
	if hashlines, ok := rel.MultiValues[transport.AlgoToString(transport.CHECKSUM_ALGO_SHA256)]; ok {
		for _, line := range hashlines {
			parts := strings.Fields(line)
			if len(parts) != 3 {
				fmt.Println(parts)
				fmt.Printf("expected 3 parts but got %d\n", len(parts))
				continue
			}

			// Checksum
			checksum, err := hex.DecodeString(parts[0])
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Size
			sizeStr := parts[1]
			size, err := strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}

			path := path2.Join(distPrefix, parts[2])

			vf := transport.VerifiedFile{
				Path: path,
				Size: uint64(size),
				Checksums: map[transport.ChecksumAlgorithm][]byte{
					transport.CHECKSUM_ALGO_SHA256: checksum,
				},
			}
			files[path] = vf
		}
	} else {
		log.Fatal(true, "Without SHA256 hashes I give up")
		// TODO: don't give up, just work with what we have
	}

	return meta, files
}

//func (dist *Distribution) Releases() (DistributionMeta, iter.Iterator[Release]) {
//	//http://archive.ubuntu.com/ubuntu/dists/bionic-updates/Release
//	resp, err := http.Get(dist.Url.String() + "/Release")
//	if err != nil {
//		return DistributionMeta{}, iter.Empty[Release]()
//	}
//	defer func(Body io.ReadCloser) {
//		_ = Body.Close()
//	}(resp.Body)
//
//	return ParseReleaseFile(resp.Body)
//}

//func (dist *Distribution) Component(componentName string) Component {
//	return Component{
//		ArchiveRoot:  dist.Archive,
//		RelativePath: componentName,
//	})
//}

//func (dist *Distribution) Components() []Component {
//	cmps := make([]Component, len(dist.DistributionMeta.Components))
//	for _, componentName := range dist.DistributionMeta.Components {
//		cmps = append(cmps, dist.Component(componentName))
//	}
//	return cmps
//}
