package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

var chunkSize int64 = 512

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fr, err := os.OpenFile(fromPath, os.O_RDONLY, 0444)
	if err != nil {
		return err
	}
	defer fr.Close()

	fileStat, err := fr.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}
	fileSize := fileStat.Size()
	fileMode := fileStat.Mode()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if offset != 0 {
		if _, err := fr.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}

	fw, err := os.Create(toPath)
	if err != nil {
		return nil
	}
	defer fw.Close()

	if err := fw.Chmod(fileMode); err != nil {
		return err
	}

	// calculate bytes length to read
	var realSizeToCopy = fileSize - offset
	if limit < realSizeToCopy && limit != 0 {
		realSizeToCopy = limit
	}

	bar := pb.Full.Start64(realSizeToCopy)
	for {
		if realSizeToCopy < chunkSize {
			chunkSize = realSizeToCopy
		}
		writtenBytes, err := io.CopyN(fw, fr, chunkSize)
		bar.Add64(writtenBytes)
		realSizeToCopy -= writtenBytes
		if realSizeToCopy == 0 {
			break
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	bar.Finish()

	return nil
}
