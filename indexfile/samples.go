package indexfile

// Sample 是固定场景，用来复现题面里的锚点。
type Sample struct {
	Name    string
	Payload []byte
	Paths   []string
}

const zeroSHA = "0000000000000000000000000000000000000000"

func Samples() []Sample {
	padded := File(2, [][]byte{
		Record("a.txt", 0o100644, 11, zeroSHA),
		Record("bb.txt", 0o100755, 22, zeroSHA),
	}, -1)
	extra := File(2, [][]byte{
		Record("a.txt", 0o100644, 11, zeroSHA),
		Record("b.txt", 0o100644, 22, zeroSHA),
	}, 1)
	v4 := File(4, [][]byte{Record("a.txt", 0o100644, 11, zeroSHA)}, -1)
	return []Sample{
		{Name: "pad", Payload: padded, Paths: []string{"a.txt", "bb.txt"}},
		{Name: "count", Payload: extra, Paths: []string{"a.txt", "b.txt"}},
		{Name: "version", Payload: v4, Paths: []string{"a.txt"}},
	}
}
