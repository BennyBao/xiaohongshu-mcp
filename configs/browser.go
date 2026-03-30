package configs

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
