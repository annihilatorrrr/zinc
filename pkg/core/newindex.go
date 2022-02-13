package core

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/directory"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

// NewIndex creates an instance of a physical zinc index that can be used to store and retrieve data.
func NewIndex(name string, storageType string) (*Index, error) {
	var config bluge.Config
	if storageType == "s3" {
		ZINC_S3_BUCKET := zutils.GetEnv("ZINC_S3_BUCKET", "")
		config = directory.GetS3Config(ZINC_S3_BUCKET, name)
	} else if storageType == "minio" {
		MINIO_BUCKET := zutils.GetEnv("ZINC_MINIO_BUCKET", "")
		config = directory.GetMinIOConfig(MINIO_BUCKET, name)
	} else { // Default storage type is disk
		ZINC_DATA_PATH := zutils.GetEnv("ZINC_DATA_PATH", "./data")
		config = bluge.DefaultConfig(ZINC_DATA_PATH + "/" + name)
	}

	writer, err := bluge.OpenWriter(config)

	if err != nil {
		return nil, err
	}

	index := &Index{
		Name:        name,
		Writer:      writer,
		StorageType: storageType,
	}

	mapping, err := index.GetStoredMapping()
	if err != nil {
		return nil, err
	}

	index.CachedMapping = mapping

	return index, nil
}

func GetIndex(indexName string) (*Index, bool) {
	index, ok := ZINC_INDEX_LIST[indexName]
	return index, ok
}

func FormatMapping(mappings *Mappings) (map[string]string, error) {
	newMappings := make(map[string]string)
	for field, prop := range mappings.Properties {
		ptype := strings.ToLower(prop.Type)
		switch ptype {
		case "text", "keyword", "numeric", "bool", "time":
		// ptype can be used as is. no need to do anything
		case "boolean":
			ptype = "bool"
		case "date", "datetime":
			ptype = "time"
		default:
			return nil, fmt.Errorf("mappings unsupport type: [%s] for field [%s]", prop.Type, field)
		}
		newMappings[field] = ptype
	}

	return newMappings, nil
}
