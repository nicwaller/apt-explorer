package apt

import (
	"apt-explorer/lib/apt/transport"
)

type Component struct {
	RelativePath string // relative to archive root
	Files        []transport.VerifiedFile
}
