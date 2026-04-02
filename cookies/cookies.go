package cookies

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
)

// AccountsConfig accounts.json 配置结构
type AccountsConfig struct {
	DefaultAccount string                 `json:"default_account"`
	Accounts       map[string]interface{} `json:"accounts"`
}

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

// getDefaultAccount 从 accounts.json 获取默认账号
// 如果指定了 default_account 则使用它；否则如果只有一个账号则使用该账号
func getDefaultAccount() string {
	accountsPath := filepath.Join(configs.GetWorkspace(), "skills", "post-to-xhs", "config", "accounts.json")
	data, err := os.ReadFile(accountsPath)
	if err != nil {
		panic("无法读取 accounts.json: " + err.Error())
	}

	var cfg AccountsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		panic("无法解析 accounts.json: " + err.Error())
	}

	// 如果指定了默认账号，使用它
	if cfg.DefaultAccount != "" {
		return cfg.DefaultAccount
	}

	// 如果只有一个账号，使用它
	if len(cfg.Accounts) == 1 {
		for name := range cfg.Accounts {
			return name
		}
	}

	panic("accounts.json 中未指定 default_account 且账号数量不为 1")
}

// GetCookiesFilePath 获取 cookies 文件路径
// 自动从 accounts.json 获取默认账号
func GetCookiesFilePath() string {
	return GetCookiesFilePathWithAccount("")
}

// GetCookiesFilePathWithAccount 根据账号名获取 cookies 文件路径
// 如果 account 为空，从 accounts.json 获取默认账号
// 返回路径格式：{OPENCLAW_WORKSPACE}/xhs-accounts/{account}/cookies.json
func GetCookiesFilePathWithAccount(account string) string {
	if account == "" {
		account = getDefaultAccount()
	}

	accountDir := filepath.Join(configs.GetWorkspace(), "xhs-accounts", account)
	if err := os.MkdirAll(accountDir, 0755); err != nil {
		panic("无法创建账号目录: " + accountDir + ": " + err.Error())
	}
	return filepath.Join(accountDir, "cookies.json")
}
