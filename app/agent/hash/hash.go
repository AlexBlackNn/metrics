package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"hash"
)

type MetricHash struct {
	hashCalculator hash.Hash
}

func New(cfg *configagent.Config) *MetricHash {
	return &MetricHash{
		hashCalculator: hmac.New(sha256.New, []byte(cfg.HashKey)),
	}
}

func (mh *MetricHash) MetricHash(body string) string {
	mh.hashCalculator.Write([]byte(body))
	metricHash := mh.hashCalculator.Sum(nil)
	base64Result := make([]byte, base64.StdEncoding.EncodedLen(len(metricHash)))
	base64.StdEncoding.Encode(base64Result, metricHash)
	return string(base64Result)
}
