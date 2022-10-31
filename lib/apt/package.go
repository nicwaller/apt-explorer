package apt

type Package struct {
	Filename            string // Required. The path of the package archive relative to the base directory of the repository
	Checksum            string // Required.
	ChecksumAlgorithm   string // Required. May be MD5, SHA1, SHA256, or SHA512
	Description         string // Required (unless DescriptionMD5 is provided)
	DescriptionMD5      string // Optional. An MD5 checksum (in hex representation) of the complete English language description.
	PhasedUpdatePercent uint8  // Optional. An integer (0-100%) defining the percentage of machines that should get this update.
}
