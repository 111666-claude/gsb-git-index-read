# git-index-read

读 git 索引文件（`DIRC` 第二版）并列出条目。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/gitidx --file .git/index
go run ./cmd/gitidx --sample=pad
go run ./cmd/gitidx --sample=count
go run ./cmd/gitidx --sample=version
go run ./cmd/gitidx --sample=work
```

## 口径

- **头**：12 字节，前 4 字节是 `DIRC`，第 4 到第 8 字节是版本，只支持第二版；第 8 到第 12 字节是条目数。
- **条目**：固定 62 字节，其中第 24 字节起是模式、第 36 字节起是文件长度、第 40 到第 60 字节是对象号、最后 2 字节的低 12 位是名字长度。
- **对齐**：每条条目（含名字）按 8 字节对齐，不足补零。
- **条数**：只读头部声明的那几条，多出来的字节不算。

## 不变量

- 条目数与头部声明一致。
- 路径名按 UTF-8 解码，不带补齐用的零字节。
- `scanned` 不随条目数乘查询次数放大：两千条条目查两千次的 `scanned` 不超过 9000。

## 输出契约

打印对齐表格，列顺序为 `path`、`mode`、`size`、`sha`（对象号前 8 位）。
`--sample=work` 只打印 `{"entries": N, "scanned": N}`。
