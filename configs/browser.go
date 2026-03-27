package configs

var (
	useHeadless = true

	binPath     = ""
	userDataDir = ""
	account     = ""
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

func SetUserDataDir(dir string) {
	userDataDir = dir
}

func GetUserDataDir() string {
	return userDataDir
}

func SetAccount(acc string) {
	account = acc
}

func GetAccount() string {
	return account
}
