package outformats

type OutWriter interface {
	WriteAll(records [][]string) error // 写入
}
