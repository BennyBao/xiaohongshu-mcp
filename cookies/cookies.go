package cookies

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
)

type Cookier interface {
	LoadCookies() ([]byte, error)
	SaveCookies(data []byte) error
	DeleteCookies() error
}

type localCookie struct {
	path string
}

func NewLoadCookie(path string) Cookier {
	if path == "" {
		panic("path is required")
	}

	return &localCookie{
		path: path,
	}
}

// LoadCookies 从文件中加载 cookies。
func (c *localCookie) LoadCookies() ([]byte, error) {

	data, err := os.ReadFile(c.path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read cookies from tmp file")
	}

	return data, nil
}

// SaveCookies 保存 cookies 到文件中。
func (c *localCookie) SaveCookies(data []byte) error {
	return os.WriteFile(c.path, data, 0644)
}

// DeleteCookies 删除 cookies 文件。
func (c *localCookie) DeleteCookies() error {
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		// 文件不存在，返回 nil（认为已经删除）
		return nil
	}
	return os.Remove(c.path)
}

// GetCookiesFilePath 获取 cookies 文件路径。
// 为了向后兼容，如果旧路径 /tmp/cookies.json 存在，则继续使用；
// 否则使用当前目录下的 cookies.json
func GetCookiesFilePath() string {
	return GetCookiesFilePathWithAccount("")
}

// GetCookiesFilePathWithAccount 根据账号名获取 cookies 文件路径
// 如果 account 为空，使用默认路径；否则使用 {OPENCLAW_WORKSPACE}/{account}/cookies.json
func GetCookiesFilePathWithAccount(account string) string {
	// 旧路径：/tmp/cookies.json（仅在无账号参数时检查）
	if account == "" {
		tmpDir := os.TempDir()
		oldPath := filepath.Join(tmpDir, "cookies.json")

		// 检查旧路径文件是否存在
		if _, err := os.Stat(oldPath); err == nil {
			return oldPath
		}

		path := os.Getenv("COOKIES_PATH")
		if path == "" {
			path = "cookies.json"
		}
		return path
	}

	// 多账号模式：使用 {OPENCLAW_WORKSPACE}/{account}/cookies.json
	accountDir := filepath.Join(configs.GetWorkspace(), account)
	if err := os.MkdirAll(accountDir, 0755); err != nil {
		panic("无法创建账号目录: " + accountDir + ": " + err.Error())
	}
	return filepath.Join(accountDir, "cookies.json")
}
