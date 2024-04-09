package outformats

// OutWriter format writer interface
type OutWriter interface {
	WriteAll(records [][]string) error // 写入
	Flush()
}
