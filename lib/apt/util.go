package apt

type SourcePackage struct {
	Package            *string
	Priority           string
	Section            string
	InstalledSize      uint64
	Maintainer         string // TODO: deduplicate this
	OriginalMaintainer string // TODO: deduplicate this
	Architecture       string
	Source             string
	Version            string
	Replaces           string
	Depends            []Depdendency
	Breaks             string
	Filename           string
	Size               uint64
	Checksums          map[uint32][]byte // Parquet doesn't understand uint8
	Description        string
	Homepage           string
	DescriptionMD5     []byte
	BugsURL            string
	Origin             string
	Supported          string
	Task               []string
	// this is a substantial change
	//MD5sum             string
	//SHA1               string
	//SHA256             string
	//SHA512             string
}

type Depdendency struct {
	PackageName string
	Version     string
	Operator    string
}
