package configs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// InitLogger 初始化日志配置
// 如果指定了账号，日志会保存到 {OPENCLAW_WORKSPACE}/xhs-accounts/{account}/app.log
// 每次初始化时会清空之前的日志文件
func InitLogger(account string) error {
	if account == "" {
		// 无账号时，只输出到控制台
		logrus.SetOutput(os.Stdout)
		return nil
	}

	// 创建账号目录
	accountDir := filepath.Join(GetWorkspace(), "xhs-accounts", account)
	if err := os.MkdirAll(accountDir, 0755); err != nil {
		return err
	}

	// 创建日志文件（每次启动时清空）
	logFile := filepath.Join(accountDir, "app.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	// 同时输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(multiWriter)

	logrus.Infof("日志文件: %s", logFile)
	return nil
}
