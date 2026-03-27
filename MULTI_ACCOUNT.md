# 多账号支持使用说明

## 功能说明

本项目现已支持多小红书账号模式，主要特性：

1. **独立的浏览器 userData 目录**：每个账号可以指定独立的浏览器数据目录
2. **独立的 Cookie 文件**：Cookie 文件保存在 `accounts/{account}/cookies.json`
3. **独立的日志文件**：日志文件保存在 `accounts/{account}/app.log`

## 目录结构

```
accounts/
├── account1/
│   ├── cookies.json    # 账号1的 Cookie
│   └── app.log         # 账号1的日志
└── account2/
    ├── cookies.json    # 账号2的 Cookie
    └── app.log         # 账号2的日志
```

## 使用方法

### 1. 登录多个账号

使用 `-account` 参数指定账号名，使用 `-user-data-dir` 参数指定浏览器数据目录：

```bash
# 登录账号1
./xiaohongshu-login -account=account1 -user-data-dir=./userdata/account1

# 登录账号2
./xiaohongshu-login -account=account2 -user-data-dir=./userdata/account2
```

登录成功后，Cookie 文件会保存到：
- `accounts/account1/cookies.json`
- `accounts/account2/cookies.json`

日志文件会保存到：
- `accounts/account1/app.log`
- `accounts/account2/app.log`

### 2. 启动 MCP 服务

启动服务时指定要使用的账号：

```bash
# 使用账号1
./xiaohongshu-mcp -account=account1 -user-data-dir=./userdata/account1

# 使用账号2
./xiaohongshu-mcp -account=account2 -user-data-dir=./userdata/account2 -port=:18061
```

### 3. 参数说明

**登录工具参数：**
- `-account`：账号名称（可选，用于多账号支持）
- `-user-data-dir`：浏览器 userData 目录（可选）
- `-bin`：浏览器二进制文件路径（可选）

**MCP 服务参数：**
- `-account`：账号名称（可选，用于多账号支持）
- `-user-data-dir`：浏览器 userData 目录（可选）
- `-headless`：是否无头模式（默认 true）
- `-bin`：浏览器二进制文件路径（可选）
- `-port`：服务端口（默认 :18060）

## 向后兼容

如果不指定 `-account` 参数，程序会使用默认的 Cookie 路径，保持与旧版本的兼容性。

## 注意事项

1. 每个账号需要使用独立的 `user-data-dir` 目录
2. 多个账号同时运行时，需要使用不同的端口
3. Cookie 和日志文件会自动保存到 `accounts/{account}/` 目录中
4. 日志会同时输出到控制台和日志文件
