package configs

import "os"

var (
	useHeadless = true

	binPath = ""
	account = ""
)

func InitHeadless(h bool) {
	useHeadless = h
}

// IsHeadless 是否无头模式。
func IsHeadless() bool {
	return useHeadless
}

func SetBinPath(b string) {
	binPath = b
}

func GetBinPath() string {
	return binPath
}

func SetAccount(acc string) {
	account = acc
}

func GetAccount() string {
	return account
}

// GetWorkspace 从 OPENCLAW_WORKSPACE 环境变量获取工作目录。
// 若未设置则 panic，程序应在启动时尽早调用以快速失败。
func GetWorkspace() string {
	ws := os.Getenv("OPENCLAW_WORKSPACE")
	if ws == "" {
		panic("环境变量 OPENCLAW_WORKSPACE 未设置，程序无法启动")
	}
	return ws
}
