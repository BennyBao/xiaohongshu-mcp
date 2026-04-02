package configs

var (
	useHeadless = true

	binPath     = ""
	account     = ""
	userDataDir = ""
	workspace   = ""
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

// GetWorkspace 获取工作目录。
// 若未设置则 panic，程序应在启动时尽早调用以快速失败。
func GetWorkspace() string {
	if workspace == "" {
		panic("workDir 未设置，程序无法启动")
	}
	return workspace
}

// SetWorkspace 设置工作目录
func SetWorkspace(dir string) {
	if dir == "" {
		panic("workDir 不能为空")
	}
	workspace = dir
}

// SetUserDataDir 设置浏览器 userData 目录
func SetUserDataDir(dir string) {
	userDataDir = dir
}

// GetUserDataDir 获取浏览器 userData 目录
func GetUserDataDir() string {
	return userDataDir
}
