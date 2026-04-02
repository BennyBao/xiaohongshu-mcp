package browser

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
	"github.com/xpzouying/xiaohongshu-mcp/cookies"
)

// Browser 封装 rod.Browser 和 launcher
type Browser struct {
	browser  *rod.Browser
	launcher *launcher.Launcher
}

// Close 关闭浏览器并清理资源
func (b *Browser) Close() {
	b.browser.MustClose()
	b.launcher.Cleanup()
}

// NewPage 创建新页面（启用 stealth 模式）
func (b *Browser) NewPage() *rod.Page {
	return stealth.MustPage(b.browser)
}

type browserConfig struct {
	binPath     string
	account     string
	userDataDir string
}

type Option func(*browserConfig)

func WithBinPath(binPath string) Option {
	return func(c *browserConfig) {
		c.binPath = binPath
	}
}

func WithAccount(account string) Option {
	return func(c *browserConfig) {
		c.account = account
	}
}

func WithUserDataDir(dir string) Option {
	return func(c *browserConfig) {
		c.userDataDir = dir
	}
}

// maskProxyCredentials masks username and password in proxy URL for safe logging.
func maskProxyCredentials(proxyURL string) string {
	u, err := url.Parse(proxyURL)
	if err != nil || u.User == nil {
		return proxyURL
	}
	if _, hasPassword := u.User.Password(); hasPassword {
		u.User = url.UserPassword("***", "***")
	} else {
		u.User = url.User("***")
	}
	return u.String()
}

func NewBrowser(headless bool, options ...Option) *Browser {
	cfg := &browserConfig{}
	for _, opt := range options {
		opt(cfg)
	}

	// 未指定 userDataDir 时，按账号自动设置
	if cfg.userDataDir == "" && cfg.account != "" {
		cfg.userDataDir = filepath.Join(
			configs.GetWorkspace(),
			"xhs-accounts",
			cfg.account,
			"browser-data",
		)
	}

	l := launcher.New().
		Headless(headless).
		Set("--no-sandbox").
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	if cfg.binPath != "" {
		l = l.Bin(cfg.binPath)
	}

	if proxy := os.Getenv("XHS_PROXY"); proxy != "" {
		l = l.Proxy(proxy)
		logrus.Infof("Using proxy: %s", maskProxyCredentials(proxy))
	}

	if cfg.userDataDir != "" {
		if err := os.MkdirAll(cfg.userDataDir, 0755); err != nil {
			panic("无法创建 userDataDir: " + err.Error())
		}
		l = l.UserDataDir(cfg.userDataDir)
		logrus.Infof("使用浏览器数据目录: %s", cfg.userDataDir)
	}

	browserURL := l.MustLaunch()

	b := rod.New().
		ControlURL(browserURL).
		MustConnect()

	// 加载 cookies
	cookiePath := cookies.GetCookiesFilePathWithAccount(cfg.account)
	cookieLoader := cookies.NewLoadCookie(cookiePath)
	if data, err := cookieLoader.LoadCookies(); err == nil {
		var cks []*proto.NetworkCookie
		if err := json.Unmarshal(data, &cks); err == nil {
			b.MustSetCookies(cks...)
			logrus.Debugf("loaded cookies from file successfully: %s", cookiePath)
		}
	} else {
		logrus.Warnf("failed to load cookies: %v", err)
	}

	return &Browser{
		browser:  b,
		launcher: l,
	}
}
