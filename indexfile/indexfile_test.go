package indexfile

import "testing"

func TestSingleEntry(t *testing.T) {
	data := File(2, [][]byte{Record("a.txt", 0o100644, 11, zeroSHA)}, -1)
	got, err := Read(data, &Counter{})
	if err != nil {
		t.Fatalf("解析失败：%v", err)
	}
	if len(got) != 1 || got[0].Path != "a.txt" || got[0].Mode != 0o100644 {
		t.Fatalf("条目解析不对：%+v", got)
	}
}

func TestBadHeader(t *testing.T) {
	if _, err := Read([]byte("XXXX"), &Counter{}); err == nil {
		t.Fatal("头不对应报错")
	}
}

func TestEmptyIndex(t *testing.T) {
	got, err := Read(File(2, nil, -1), &Counter{})
	if err != nil {
		t.Fatalf("解析失败：%v", err)
	}
	if len(got) != 0 {
		t.Fatalf("空索引不该有条目：%+v", got)
	}
}

func TestLookupFindsPath(t *testing.T) {
	data := File(2, [][]byte{Record("a.txt", 0o100644, 11, zeroSHA)}, -1)
	got, err := Lookup(data, []string{"a.txt"}, &Counter{})
	if err != nil {
		t.Fatalf("查找失败：%v", err)
	}
	if _, ok := got["a.txt"]; !ok {
		t.Fatalf("应找到 a.txt：%+v", got)
	}
}

func TestSizeIsRead(t *testing.T) {
	data := File(2, [][]byte{Record("a.txt", 0o100644, 11, zeroSHA)}, -1)
	got, _ := Read(data, &Counter{})
	if got[0].Size != 11 {
		t.Fatalf("文件长度应是 11，得到 %d", got[0].Size)
	}
}
