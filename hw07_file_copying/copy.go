package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return fmt.Errorf("offset is less than zero")
	}
	if limit < 0 {
		return fmt.Errorf("limit is less than zero")
	}

	fromInfo, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("from file '%s': %w", from, errors.Unwrap(err))
	}
	if !fromInfo.Mode().IsRegular() {
		return fmt.Errorf("from file '%s': %w", from, ErrUnsupportedFile)
	}
	fromSize := fromInfo.Size()
	if offset > fromSize {
		return ErrOffsetExceedsFileSize
	}

	fromFile, err := openFrom(fromPath, offset)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := openTo(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	realLimit := realLimit(offset, limit, fromSize)

	bar := pb.Full.Start64(realLimit)
	barReader := bar.NewProxyReader(fromFile)

	_, err = io.CopyN(toFile, barReader, realLimit)

	bar.Finish()

	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

// Открываем файл для чтения.
func openFrom(name string, offset int64) (*os.File, error) {
	fromFile, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("from file '%s': %w", from, errors.Unwrap(err))
	}

	if offset > 0 {
		_, err = fromFile.Seek(offset, io.SeekStart)
		if err != nil {
			_ = fromFile.Close()
			return nil, fmt.Errorf("from file '%s': %w", from, err)
		}
	}

	return fromFile, nil
}

// Открываем файл для записи.
func openTo(name string) (*os.File, error) {
	// Удаляем файл вместо перезаписи (так быстрее).
	if fileExists(name) {
		err := os.Remove(name)
		if err != nil {
			return nil, fmt.Errorf("to file '%s': %w", to, errors.Unwrap(err))
		}
	}

	toFile, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("to file '%s': %w", to, errors.Unwrap(err))
	}

	return toFile, nil
}

// Проверка наличия файла.
func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// Вычисление количества копируемых байт если limit=0
// или offset+limit больше реального размера исходного файла.
func realLimit(offset, limit, fileSize int64) int64 {
	if limit == 0 || offset+limit > fileSize {
		return fileSize - offset
	}
	return limit
}
