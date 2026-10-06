// 手工拼索引文件，供测试与场景使用。
package indexfile

import (
	"encoding/binary"
	"encoding/hex"
)

func Record(name string, mode, size uint32, sha string) []byte {
	buf := make([]byte, EntSize)
	binary.BigEndian.PutUint32(buf[24:28], mode)
	binary.BigEndian.PutUint32(buf[36:40], size)
	raw, err := hex.DecodeString(sha)
	if err != nil || len(raw) != 20 {
		panic("sha 必须是 40 位十六进制")
	}
	copy(buf[40:60], raw)
	binary.BigEndian.PutUint16(buf[60:62], uint16(len(name)&0x0FFF))
	full := append(buf, []byte(name)...)
	for len(full)%8 != 0 {
		full = append(full, 0)
	}
	return full
}

func File(version uint32, entries [][]byte, count int) []byte {
	head := make([]byte, 12)
	copy(head, Sig)
	binary.BigEndian.PutUint32(head[4:8], version)
	if count < 0 {
		count = len(entries)
	}
	binary.BigEndian.PutUint32(head[8:12], uint32(count))
	out := head
	for _, e := range entries {
		out = append(out, e...)
	}
	return out
}
