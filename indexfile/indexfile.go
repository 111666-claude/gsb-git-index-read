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
	off := 12
	out := []Entry{}
	for off+EntSize <= len(data) {
		c.Scanned++
		mode := binary.BigEndian.Uint32(data[off+24 : off+28])
		size := binary.BigEndian.Uint32(data[off+36 : off+40])
		sha := hex.EncodeToString(data[off+40 : off+60])
		flags := binary.BigEndian.Uint16(data[off+60 : off+62])
		nameLen := int(flags & 0x0FFF)
		if off+EntSize+nameLen > len(data) {
			break
		}
		name := string(data[off+EntSize : off+EntSize+nameLen])
		out = append(out, Entry{Path: name, Mode: mode, Size: size, SHA: sha[:8]})
		off += EntSize + nameLen
	}
	return out, nil
}

// Lookup 按路径取记录。
func Lookup(data []byte, paths []string, c *Counter) (map[string]Entry, error) {
	found := map[string]Entry{}
	for _, want := range paths {
		entries, err := Read(data, c)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.Path == want {
				found[want] = entry
			}
		}
	}
	return found, nil
}
