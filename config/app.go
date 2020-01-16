package config

var mode Mode

type Mode byte

const (
	ModeProd Mode = iota
	ModeDev
	ModeTest
)

func init() {
	SetMode(ModeDev)
}

// SetMode 设置app运行模式
func SetMode(m Mode) {
	mode = m
}

// GetMode 获取app运行模式
func GetMode() Mode {
	return mode
}
