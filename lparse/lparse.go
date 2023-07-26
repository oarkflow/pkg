package lparse

import (
	"bufio"
	json "encoding/json"
	"github.com/oarkflow/pkg/str"
	"github.com/sourcegraph/conc/iter"
	"github.com/sourcegraph/conc/pool"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

func read(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return ln, err
}

func mapLines(lines []string) []map[string]any {
	return iter.Map[string, map[string]any](lines, func(s *string) map[string]any {
		var entry map[string]any
		json.Unmarshal([]byte(*s), &entry)
		return entry
	})
}

func ParseLog(logFile string) []map[string]any {
	// start := time.Now()
	file, err := os.Open(logFile)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewReader(file)
	var logEntries []map[string]any

	linesChunkLen := 64 * 1024
	linesChunkPoolAllocated := int64(0)
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]string, 0, linesChunkLen)
		atomic.AddInt64(&linesChunkPoolAllocated, 1)
		return lines
	}}
	lines := linesPool.Get().([]string)[:0]

	entriesPoolAllocated := int64(0)
	entriesPool := sync.Pool{New: func() interface{} {
		entries := make([]map[string]any, 0, linesChunkLen)
		atomic.AddInt64(&entriesPoolAllocated, 1)
		return entries
	}}
	wg := pool.New().WithMaxGoroutines(runtime.NumCPU())
	mutex := &sync.Mutex{}
	for {
		line, err := read(scanner)
		if err != nil {
			if err == io.EOF && len(lines) > 0 {
				linesToProcess := lines
				wg.Go(func() {
					entries := entriesPool.Get().([]map[string]any)[:0]
					entries = mapLines(linesToProcess)
					linesPool.Put(linesToProcess)
					mutex.Lock()
					for _, entry := range entries {
						logEntries = append(logEntries, entry)
					}
					mutex.Unlock()
					entriesPool.Put(entries)
				})
				lines = linesPool.Get().([]string)[:0]
			}
			break
		}
		lines = append(lines, str.FromByte(line))
		if len(lines) == linesChunkLen {
			linesToProcess := lines
			wg.Go(func() {
				entries := entriesPool.Get().([]map[string]any)[:0]
				entries = mapLines(linesToProcess)
				linesPool.Put(linesToProcess)
				mutex.Lock()
				for _, entry := range entries {
					logEntries = append(logEntries, entry)
				}
				mutex.Unlock()
				entriesPool.Put(entries)
			})
			lines = linesPool.Get().([]string)[:0]
		}
	}
	wg.Wait()
	// fmt.Printf("time: %v, lines: %d\n", time.Since(start), len(logEntries))
	return logEntries
}
