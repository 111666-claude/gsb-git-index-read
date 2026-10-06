// Command gitidx 是索引解析的场景入口。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"example.com/git-index-read/indexfile"
)

func emit(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func render(entries []indexfile.Entry) string {
	rows := [][4]string{{"path", "mode", "size", "sha"}}
	for _, e := range entries {
		rows = append(rows, [4]string{e.Path, fmt.Sprintf("%o", e.Mode),
			fmt.Sprintf("%d", e.Size), e.SHA})
	}
	widths := [4]int{}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	var b strings.Builder
	for _, row := range rows {
		for i, cell := range row {
			b.WriteString(cell)
			if i < 3 {
				b.WriteString(strings.Repeat(" ", widths[i]-len(cell)+2))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func main() {
	sample := flag.String("sample", "", "pad|count|version|work")
	file := flag.String("file", "", "索引文件路径")
	flag.Parse()

	if *file != "" {
		data, err := os.ReadFile(*file)
		if err != nil {
			panic(err)
		}
		entries, err := indexfile.Read(data, &indexfile.Counter{})
		if err != nil {
			emit(map[string]string{"error": err.Error()})
			return
		}
		fmt.Print(render(entries))
		return
	}

	for _, s := range indexfile.Samples() {
		if s.Name != *sample {
			continue
		}
		entries, err := indexfile.Read(s.Payload, &indexfile.Counter{})
		if err != nil {
			emit(map[string]string{"error": err.Error()})
			return
		}
		switch *sample {
		case "version":
			emit(map[string]any{"entries": len(entries)})
		default:
			paths := []string{}
			for _, e := range entries {
				paths = append(paths, e.Path)
			}
			emit(map[string]any{"paths": paths})
		}
		return
	}
	if *sample == "work" {
		var items [][]byte
		var paths []string
		for i := 0; i < 2000; i++ {
			name := fmt.Sprintf("f%04d.txt", i)
			items = append(items, indexfile.Record(name, 0o100644, uint32(i), strings.Repeat("0", 40)))
			paths = append(paths, name)
		}
		data := indexfile.File(2, items, -1)
		c := &indexfile.Counter{}
		indexfile.Lookup(data, paths, c)
		emit(map[string]int{"entries": 2000, "scanned": c.Scanned})
		return
	}
	fmt.Fprintln(os.Stderr, "未知场景："+*sample)
	os.Exit(2)
}
