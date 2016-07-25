// Logger используется совместно с пакетом log для простой ротации логов
//
//   log.SetOutput(&logger.Logger{
//       Filename: "/var/log/app/foo.log",
//	     MaxSize:  100, // megabytes
//       Everyday: true,
//   })
//   log.Println("Test")
//
// Filename задает местоположение log файла
// MaxSize задает максимальный размер log файла в мегабайтах
// Everyday указадывает на то, что в одном файле не может быть логов за разные дни
// т.е. Logger каждый день будет создавать новый файл, а если тот привысит размер
// указанный в MaxSize, то будет создан ещё один файл для этого дня.
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	backupTimeFormat = "[2006-01-02][15-04-05]"
	defaultMaxSize   = 100
	megabyte         = 1048576 // bytes
)

var _ io.WriteCloser = (*Logger)(nil)

type Logger struct {
	Filename string
	MaxSize  int64
	Everyday bool
	size     int64
	file     *os.File
	time     time.Time
	mu       sync.Mutex
}

func (l *Logger) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	writeLen := int64(len(p))
	if writeLen > l.max() {
		return 0, fmt.Errorf("logger: write length %d bytes exceeds maximum file size %d bytes", writeLen, l.max())
	}

	if l.file == nil {
		if err := l.initFile(writeLen); err != nil {
			return 0, err
		}
	}
	if l.Everyday && !isToday(l.time) {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}
	if l.size+writeLen > l.max() {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := l.file.Write(p)

	l.size += int64(n)
	l.time = time.Now()

	return n, err
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

func (l *Logger) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

func (l *Logger) initFile(writeLen int64) error {
	filename := l.filename()
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return l.newFile()
	}
	if err != nil {
		return fmt.Errorf("logger: error getting logfile: %s", err)
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return fmt.Errorf("logger: error opening logfile: %s", err)
	}

	l.file = file
	l.size = info.Size()
	l.time = info.ModTime()

	return nil
}

func (l *Logger) newFile() error {
	filename := l.filename()

	err := os.MkdirAll(filepath.Dir(filename), 0744)
	if err != nil {
		return fmt.Errorf("logger: can't make directories for new logfile: %s", err)
	}

	if _, err := os.Stat(filename); err == nil {
		newfilename := backupName(filename)
		if err := os.Rename(filename, newfilename); err != nil {
			return fmt.Errorf("logger: can't rename logfile: %s", err)
		}
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0644)

	if err != nil {
		return fmt.Errorf("logger: can't open new logfile: %s", err)
	}

	l.file = file
	l.size = 0

	return nil
}

func (l *Logger) rotate() error {
	if err := l.close(); err != nil {
		return err
	}
	if err := l.newFile(); err != nil {
		return err
	}
	return nil
}

func (l *Logger) filename() string {
	return l.Filename
}

func backupName(name string) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]
	datetime := time.Now().Format(backupTimeFormat)

	return filepath.Join(dir, prefix+datetime+ext)
}

func (l *Logger) max() int64 {
	if l.MaxSize == 0 {
		return defaultMaxSize * megabyte
	}
	return l.MaxSize * megabyte
}

// Быстрый способ проверить дату на принадлежность к сегодняшнему дню
func isToday(t time.Time) bool {
	n := time.Now()
	return t.Year() == n.Year() && t.Month() == n.Month() && t.Day() == n.Day()
}
