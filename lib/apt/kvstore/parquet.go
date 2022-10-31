package kvstore

import (
	"compress/bzip2"
	"compress/gzip"
	"github.com/segmentio/parquet-go"
	"github.com/segmentio/parquet-go/compress/zstd"
	"io"
	"os"
	"path"
	"strings"
)

// performance testing
//
// 3171/3172 Package records...
//
// Apt format, plaintext: 4882 KiB
// Apt format, gzip: 1303 KiB
// Parquet:
//  Snappy 2128 KiB
//  Gzip 493 KiB
//  Brotli 352 KiB
//  Zstd 247 KiB (fastest)
//  Zstd 201 KiB (default) <-- winner
//  Zstd 241 KiB (better)
//  Zstd 234 KiB (best)

//goland:noinspection GoUnusedExportedFunction
func ConvertFileToParquet(inFile string, outFile string) error {
	var reader io.Reader
	reader, err := os.Open(inFile)
	if err != nil {
		return err
	}
	// TODO: also support .bz2 files
	switch path.Ext(inFile) {
	case "gz":
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return err
		}
	case "bz2":
		reader = bzip2.NewReader(reader)
	}
	return ConvertToParquet(reader, outFile)
}

// ConvertToParquet inFile MUST be plaintext
func ConvertToParquet(inReader io.Reader, outFile string) error {
	rows := make([]Block, 0)

	// PERF: this uses a lot of memory, but it's probably fine.
	// I assume these file conversions will be rare. -NW
	err := Parse(inReader, func(block Block) {
		rows = append(rows, block)
	})
	if err != nil {
		return err
	}

	compressZstd := parquet.Compression(&zstd.Codec{
		Level: zstd.SpeedDefault, // best size AND speed in my testing -NW
	})
	err = parquet.WriteFile(outFile, rows, compressZstd)
	if err != nil {
		return err
	}

	return nil
}

func QueryParquet(pqFilename string, needle string) ([]Block, error) {
	//rdr, e := os.Open(pqFilename)
	//if e != nil {
	//	return []Block{}, e
	//}
	//defer func(rdr *os.File) {
	//	_ = rdr.Close()
	//}(rdr)

	rows, e := parquet.ReadFile[Block](pqFilename)
	if e != nil {
		return []Block{}, e
	}

	matches := make([]Block, 0)
	for _, r := range rows {
		if haystack, ok := r.SingleValues["Package"]; ok {
			if strings.Contains(haystack, needle) {
				matches = append(matches, r)
			}
		}
	}

	return matches, nil
}

//func streamingRead(pqFile string) {
//	r2, _ := os.Open(pqFile)
//	f, err := parquet.OpenFile(r2, fi.Size())
//	if err != nil {
//		fmt.Println("failed opening parquet file")
//	}
//	fmt.Println(pqFile)
//
//	for _, rowGroup := range f.RowGroups() {
//		fmt.Printf(" RowGroup with %d rows\n", rowGroup.NumRows())
//		for _, columnChunk := range rowGroup.ColumnChunks() {
//			fmt.Printf("  ColumnChunk with %d values of type %v\n", columnChunk.NumValues(), columnChunk.Type())
//
//			pgs := columnChunk.Pages()
//			for {
//				p, err := pgs.ReadPage()
//				if err != nil {
//					break
//				}
//				fmt.Printf("   Page type=%v numRows=%d numValues=%d \n", p.Type(), p.NumRows(), p.NumValues())
//				d := p.Values()
//				bb := make([]parquet.Value, 1024)
//				for {
//					nn, ee := d.ReadValues(bb)
//					if nn == 0 { // improper
//						break
//					}
//					if ee == io.EOF { // proper
//						break
//					}
//					fmt.Printf("    %v\n", bb[0])
//					fmt.Printf("    %v\n", bb[1])
//					fmt.Printf("    %v\n", bb[2])
//					break
//				}
//
//			}
//		}
//	}
//}
