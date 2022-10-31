package transport

import (
	"apt-explorer/lib/log"
	"encoding/json"
	"io"
	"os"
	"time"
)

const negativeCacheFilepath = "/tmp/apt-explorer/negative-cache"

// FIXME: negative cache needs to work with fully-qualified URLs, maybe
var negativeCache = make(map[string]int64)

func init() {
	LoadNegativeCache()
}

func LoadNegativeCache() {
	jsonFile, err := os.Open(negativeCacheFilepath)
	if err != nil {
		log.Error("Failed to open %s", negativeCacheFilepath)
		return
	}
	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)

	// read our opened xmlFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &negativeCache)
	if err != nil {
		log.Error("Failed to unmarshal json of negative cache")
		return
	}
}

func SaveNegativeCache() {
	jsonFile, err := os.Create(negativeCacheFilepath)
	if err != nil {
		log.Error("Failed to open %s", negativeCacheFilepath)
		log.Error("negative cache may have been destroyed")
		return
	}
	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	bytes, err := json.Marshal(negativeCache)
	if err != nil {
		log.Error("Failed to unmarshal json of negative cache")
		return
	}

	written, err := jsonFile.Write(bytes)
	_ = written
	if err != nil {
		log.Error("failed writing bytes to negative cache file")
		return
	}
	//log.Debug("Wrote %d bytes to negative cache file", written)
}

//goland:noinspection GoUnusedExportedFunction
func IsInNegativeCache(key string) bool {
	_, found := negativeCache[key]
	return found
}

//goland:noinspection GoUnusedExportedFunction
func AddToNegativeCache(entry string) {
	negativeCache[entry] = time.Now().Unix()
	SaveNegativeCache()
}
