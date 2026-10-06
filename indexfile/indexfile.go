// Package indexfile 读 git 索引文件（DIRC v2）。
package indexfile

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
)

// 索引文件常量。
const (
	Sig     = "DIRC"
	EntSize = 62
)

// Entry 是一条索引记录。
type Entry struct {
	Path string `json:"path"`
	Mode uint32 `json:"mode"`
	Size uint32 `json:"size"`
	SHA  string `json:"sha"`
}

// Counter 记录本轮读过的记录条数。
type Counter struct {
	Scanned int
}

// Read 解析索引文件。
func Read(data []byte, c *Counter) ([]Entry, error) {
	if len(data) < 12 || string(data[0:4]) != Sig {
		return nil, errors.New("bad-header")
	}
	if binary.BigEndian.Uint32(data[4:8]) != 2 {
		return nil, errors.New("unsupported-version")
	}
	count := int(binary.BigEndian.Uint32(data[8:12]))
	off := 12
	out := make([]Entry, 0, count)
	for i := 0; i < count; i++ {
		if off+EntSize > len(data) {
			return nil, errors.New("truncated-entry")
		}
		c.Scanned++
		mode := binary.BigEndian.Uint32(data[off+24 : off+28])
		size := binary.BigEndian.Uint32(data[off+36 : off+40])
		sha := hex.EncodeToString(data[off+40 : off+60])
		flags := binary.BigEndian.Uint16(data[off+60 : off+62])
		nameLen := int(flags & 0x0FFF)
		// 每条条目（含名字）按 8 字节对齐，不足补零。
		entryLen := (EntSize + nameLen + 7) / 8 * 8
		if off+entryLen > len(data) {
			return nil, errors.New("truncated-entry")
		}
		name := string(data[off+EntSize : off+EntSize+nameLen])
		out = append(out, Entry{Path: name, Mode: mode, Size: size, SHA: sha[:8]})
		off += entryLen
	}
	return out, nil
}

// Lookup 按路径取记录。
func Lookup(data []byte, paths []string, c *Counter) (map[string]Entry, error) {
	entries, err := Read(data, c)
	if err != nil {
		return nil, err
	}
	byPath := make(map[string]Entry, len(entries))
	for _, entry := range entries {
		byPath[entry.Path] = entry
	}
	found := map[string]Entry{}
	for _, want := range paths {
		if entry, ok := byPath[want]; ok {
			found[want] = entry
		}
	}
	return found, nil
}
