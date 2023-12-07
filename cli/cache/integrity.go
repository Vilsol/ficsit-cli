package cache

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/puzpuzpuz/xsync/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/satisfactorymodding/ficsit-cli/utils"
)

type hashInfo struct {
	Modified time.Time
	Hash     string
	Size     int64
}

var hashCache *xsync.MapOf[string, hashInfo]

var integrityFilename = ".integrity"

func getFileHash(file string) (string, error) {
	if hashCache == nil {
		loadHashCache()
	}
	cachedHash, ok := hashCache.Load(file)
	if !ok {
		return cacheFileHash(file)
	}
	downloadCache := filepath.Join(viper.GetString("cache-dir"), "downloadCache")
	stat, err := os.Stat(filepath.Join(downloadCache, file))
	if err != nil {
		return "", errors.Wrap(err, "failed to stat file")
	}
	if stat.Size() != cachedHash.Size || stat.ModTime() != cachedHash.Modified {
		return cacheFileHash(file)
	}
	return cachedHash.Hash, nil
}

func cacheFileHash(file string) (string, error) {
	downloadCache := filepath.Join(viper.GetString("cache-dir"), "downloadCache")
	stat, err := os.Stat(filepath.Join(downloadCache, file))
	if err != nil {
		return "", errors.Wrap(err, "failed to stat file")
	}
	f, err := os.Open(filepath.Join(downloadCache, file))
	if err != nil {
		return "", errors.Wrap(err, "failed to open file")
	}
	defer f.Close()
	hash, err := utils.SHA256Data(f)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash file")
	}
	hashCache.Store(file, hashInfo{
		Hash:     hash,
		Size:     stat.Size(),
		Modified: stat.ModTime(),
	})
	saveHashCache()
	return hash, nil
}

func loadHashCache() {
	hashCache = xsync.NewMapOf[string, hashInfo]()
	cacheFile := filepath.Join(viper.GetString("cache-dir"), "downloadCache", integrityFilename)
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return
	}
	f, err := os.Open(cacheFile)
	if err != nil {
		log.Warn().Err(err).Msg("failed to open hash cache, recreating")
		return
	}
	defer f.Close()

	hashCacheJSON, err := io.ReadAll(f)
	if err != nil {
		log.Warn().Err(err).Msg("failed to read hash cache, recreating")
		return
	}

	if err := json.Unmarshal(hashCacheJSON, &hashCache); err != nil {
		log.Warn().Err(err).Msg("failed to unmarshal hash cache, recreating")
		return
	}
}

func saveHashCache() {
	cacheFile := filepath.Join(viper.GetString("cache-dir"), "downloadCache", integrityFilename)
	plainCache := make(map[string]hashInfo, hashCache.Size())
	hashCache.Range(func(k string, v hashInfo) bool {
		plainCache[k] = v
		return true
	})
	hashCacheJSON, err := json.Marshal(plainCache)
	if err != nil {
		log.Warn().Err(err).Msg("failed to marshal hash cache")
		return
	}

	if err := os.WriteFile(cacheFile, hashCacheJSON, 0o755); err != nil {
		log.Warn().Err(err).Msg("failed to write hash cache")
		return
	}
}
