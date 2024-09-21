package gzipcompressor

import (
	"compress/gzip"
	"errors"
	"net/http"
	"sync"
)

type gzipWriterPool struct {
	p sync.Pool
}

func (gp *gzipWriterPool) get(
	w http.ResponseWriter,
	compressorLevel int,
) (*GzipWriter, error) {
	gzipWriter := gp.p.Get()
	if gzipWriter == nil {
		gzipWr, err := gzip.NewWriterLevel(w, compressorLevel)
		if err != nil {
			return nil, err
		}
		return &GzipWriter{ResWriter: w, Writer: gzipWr}, nil
	}
	gzipWr, ok := gzipWriter.(*gzip.Writer)
	if !ok {
		return nil, errors.New("wrong type of gzipWriter")
	}
	// need to reset old writer state
	gzipWr.Reset(w)
	return &GzipWriter{ResWriter: w, Writer: gzipWr}, nil
}

func (gp *gzipWriterPool) put(gzipWriter *GzipWriter) error {
	err := gzipWriter.Close()
	if err != nil {
		return err
	}
	gp.p.Put(gzipWriter.Writer)
	return nil
}

func (gp *gzipWriterPool) putNoFlush(gzipWriter *GzipWriter) error {
	gp.p.Put(gzipWriter.Writer)
	return nil
}
