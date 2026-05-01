package logxext

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type BizWriter struct {
	Out io.Writer
}

func (w *BizWriter) writeLine(level string, v any) {
	msg := fmt.Sprint(v)

	ts := time.Now().Format("2006-01-02 15:04:05")
	caller := findCaller()

	level = strings.ToUpper(level)
	levelColor := colorizeLevel(level)

	line := fmt.Sprintf("%s | %s | %s | %s\n", ts, levelColor, caller, msg)
	_, _ = w.Out.Write([]byte(line))
}

func colorizeLevel(level string) string {
	switch level {
	case "DEBUG":
		return colorCyan + level + colorReset
	case "INFO":
		return colorGreen + level + colorReset
	case "WARN", "WARNING":
		return colorYellow + level + colorReset
	case "ERROR":
		return colorRed + level + colorReset
	default:
		return level
	}
}

var (
	projectRoot string
	rootOnce    sync.Once
)

func findProjectRoot(startFile string) string {
	rootOnce.Do(func() {
		dir := filepath.Dir(startFile)

		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				projectRoot = dir
				return
			}

			parent := filepath.Dir(dir)
			if parent == dir {
				projectRoot = ""
				return
			}
			dir = parent
		}
	})

	return projectRoot
}

func findCaller() string {
	for i := 2; i < 30; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		if strings.Contains(file, "/go-zero/") ||
			strings.Contains(file, "/core/logx/") ||
			strings.Contains(file, "/runtime/") ||
			strings.Contains(file, "/common/logxext/") ||
			strings.Contains(file, "/pkg/mod/") {
			continue
		}

		root := findProjectRoot(file)
		if root != "" {
			if rel, err := filepath.Rel(root, file); err == nil {
				return fmt.Sprintf("%s:%d", filepath.ToSlash(rel), line)
			}
		}

		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return "unknown"
}

func (w *BizWriter) Alert(v any)                     { w.writeLine("ALERT", v) }
func (w *BizWriter) Close() error                    { return nil }
func (w *BizWriter) Debug(v any, _ ...logx.LogField) { w.writeLine("DEBUG", v) }
func (w *BizWriter) Error(v any, _ ...logx.LogField) { w.writeLine("ERROR", v) }
func (w *BizWriter) Warn(v any, _ ...logx.LogField)  { w.writeLine("WARN", v) }
func (w *BizWriter) Info(v any, _ ...logx.LogField)  { w.writeLine("INFO", v) }
func (w *BizWriter) Severe(v any)                    { w.writeLine("SEVERE", v) }
func (w *BizWriter) Slow(v any, _ ...logx.LogField)  { w.writeLine("WARN", v) }
func (w *BizWriter) Stack(v any)                     { w.writeLine("ERROR", v) }
func (w *BizWriter) Stat(_ any, _ ...logx.LogField)  {}
