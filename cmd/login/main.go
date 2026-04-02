package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/browser"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
	"github.com/xpzouying/xiaohongshu-mcp/cookies"
	"github.com/xpzouying/xiaohongshu-mcp/xiaohongshu"
)

func main() {
	var (
		binPath     string // 浏览器二进制文件路径
		account     string // 账号名称
		userDataDir string // 浏览器 userData 目录
		workDir     string // 工作目录（必传）
	)
	flag.StringVar(&binPath, "bin", "", "浏览器二进制文件路径")
	flag.StringVar(&account, "account", "", "账号名称（用于多账号支持）")
	flag.StringVar(&userDataDir, "user-data-dir", "", "浏览器 userData 目录（可选）")
	flag.StringVar(&workDir, "workDir", "", "工作目录（必传），账号数据存储在 {workDir}/xhs-accounts/ 下")
	flag.Parse()

	if workDir == "" {
		fmt.Fprintln(os.Stderr, "错误：-workDir 参数为必传项")
		flag.Usage()
		os.Exit(1)
	}

	configs.SetWorkspace(workDir)

	// 初始化日志
	if account != "" {
		if err := configs.InitLogger(account); err != nil {
			logrus.Fatalf("failed to init logger: %v", err)
		}
	}

	// 设置配置
	if account != "" {
		configs.SetAccount(account)
	}
	if userDataDir != "" {
		configs.SetUserDataDir(userDataDir)
	}

	// 构建浏览器选项
	opts := []browser.Option{browser.WithBinPath(binPath)}
	if account != "" {
		opts = append(opts, browser.WithAccount(account))
	}
	if userDataDir != "" {
		opts = append(opts, browser.WithUserDataDir(userDataDir))
	}

	// 登录的时候，需要界面，所以不能无头模式
	b := browser.NewBrowser(false, opts...)
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	action := xiaohongshu.NewLogin(page)

	status, err := action.CheckLoginStatus(context.Background())
	if err != nil {
		logrus.Fatalf("failed to check login status: %v", err)
	}

	logrus.Infof("当前登录状态: %v", status)

	if status {
		// 已登录，刷新保存一次 cookies
		if err := saveCookies(page, account); err != nil {
			logrus.Fatalf("failed to save cookies: %v", err)
		}
		logrus.Info("已登录，cookies 已刷新")
		return
	}

	// 开始登录流程
	logrus.Info("开始登录流程...")
	if err = action.Login(context.Background()); err != nil {
		logrus.Fatalf("登录失败: %v", err)
	} else {
		if err := saveCookies(page, account); err != nil {
			logrus.Fatalf("failed to save cookies: %v", err)
		}
	}

	// 再次检查登录状态确认成功
	status, err = action.CheckLoginStatus(context.Background())
	if err != nil {
		logrus.Fatalf("failed to check login status after login: %v", err)
	}

	if status {
		logrus.Info("登录成功！")
	} else {
		logrus.Error("登录流程完成但仍未登录")
	}

}

func saveCookies(page *rod.Page, account string) error {
	cks, err := page.Browser().GetCookies()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cks)
	if err != nil {
		return err
	}

	cookiePath := cookies.GetCookiesFilePathWithAccount(account)
	cookieLoader := cookies.NewLoadCookie(cookiePath)
	return cookieLoader.SaveCookies(data)
}
