package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
)

func main() {
	var (
		headless    bool
		binPath     string // 浏览器二进制文件路径
		account     string // 账号名称
		userDataDir string // 浏览器 userData 目录
		port        string
		workDir     string // 工作目录（必传）
	)
	flag.BoolVar(&headless, "headless", true, "是否无头模式")
	flag.StringVar(&binPath, "bin", "", "浏览器二进制文件路径")
	flag.StringVar(&account, "account", "", "账号名称（用于多账号支持）")
	flag.StringVar(&userDataDir, "user-data-dir", "", "浏览器 userData 目录（可选）")
	flag.StringVar(&port, "port", ":18060", "端口")
	flag.StringVar(&workDir, "workDir", "", "工作目录（必传），账号数据存储在 {workDir}/xhs-accounts/ 下")
	flag.Parse()

	if workDir == "" {
		fmt.Fprintln(os.Stderr, "错误：-workDir 参数为必传项")
		flag.Usage()
		os.Exit(1)
	}

	if len(binPath) == 0 {
		binPath = os.Getenv("ROD_BROWSER_BIN")
	}

	configs.InitHeadless(headless)
	configs.SetBinPath(binPath)
	configs.SetAccount(account)
	configs.SetUserDataDir(userDataDir)
	configs.SetWorkspace(workDir)

	// 初始化日志
	if err := configs.InitLogger(account); err != nil {
		logrus.Fatalf("failed to init logger: %v", err)
	}

	// 初始化服务
	xiaohongshuService := NewXiaohongshuService()

	// 创建并启动应用服务器
	appServer := NewAppServer(xiaohongshuService)
	if err := appServer.Start(port); err != nil {
		logrus.Fatalf("failed to run server: %v", err)
	}
}
