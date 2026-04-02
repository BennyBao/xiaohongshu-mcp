# 多账号支持使用说明

## 功能说明

本项目现已支持多小红书账号模式，主要特性：

1. **独立的浏览器 userData 目录**：每个账号可以指定独立的浏览器数据目录
2. **独立的 Cookie 文件**：Cookie 文件保存在 `{workDir}/xhs-accounts/{account}/cookies.json`
3. **独立的日志文件**：日志文件保存在 `{workDir}/xhs-accounts/{account}/app.log`
4. **自动日志清理**：每次服务重启时自动清空日志文件

## 启动参数

| 参数 | 必填 | 说明 |
|------|------|------|
| `-workDir` | **是** | 工作目录根路径，账号数据存储在 `{workDir}/xhs-accounts/` 下，支持相对路径或绝对路径 |

未传入 `-workDir` 时，程序将拒绝启动。

## 目录结构

```
{workDir}/
└── xhs-accounts/
    ├── account1/
    │   ├── cookies.json      # 账号1的 Cookie
    │   ├── app.log           # 账号1的日志
    │   └── browser-data/     # 账号1的浏览器用户数据目录
    └── account2/
        ├── cookies.json      # 账号2的 Cookie
        ├── app.log           # 账号2的日志
        └── browser-data/     # 账号2的浏览器用户数据目录
```

## 使用方法

### 1. 登录多个账号

使用 `-account` 参数指定账号名，`-workDir` 指定工作目录：

```bash
# 登录账号1（数据保存在 {workDir}/xhs-accounts/account1/ 下）
./xiaohongshu-login -workDir=. -account=account1

# 登录账号2
./xiaohongshu-login -workDir=. -account=account2
```

如果需要指定自定义的浏览器数据目录，可以使用 `-user-data-dir` 参数：

```bash
./xiaohongshu-login -workDir=. -account=account1 -user-data-dir=./custom/userdata
```

登录成功后，文件会保存到：
- Cookie：`{workDir}/xhs-accounts/account1/cookies.json`
- 日志：`{workDir}/xhs-accounts/account1/app.log`

### 2. 启动 MCP 服务

启动服务时指定工作目录和账号：

```bash
# 使用账号1
./xiaohongshu-mcp -workDir=. -account=account1

# 使用账号2（使用不同端口）
./xiaohongshu-mcp -workDir=. -account=account2 -port=:18061

# 使用绝对路径
./xiaohongshu-mcp -workDir=/data -account=account1
```

### 3. 参数说明

**登录工具参数：**
- `-workDir`：工作目录（**必传**）
- `-account`：账号名称（用于多账号支持）
- `-user-data-dir`：浏览器 userData 目录（可选，如果不指定则自动使用 `{workDir}/xhs-accounts/{account}/browser-data`）
- `-bin`：浏览器二进制文件路径（可选）

**MCP 服务参数：**
- `-workDir`：工作目录（**必传**）
- `-account`：账号名称（用于多账号支持）
- `-user-data-dir`：浏览器 userData 目录（可选，如果不指定则自动使用 `{workDir}/xhs-accounts/{account}/browser-data`）
- `-headless`：是否无头模式（默认 true）
- `-bin`：浏览器二进制文件路径（可选）
- `-port`：服务端口（默认 :18060）

## 注意事项

1. `-workDir` 为必传参数，未传入时程序立即退出
2. 所有账号数据统一存储在 `{workDir}/xhs-accounts/` 目录下
3. 每个账号会自动使用独立的 `{workDir}/xhs-accounts/{account}/browser-data` 目录作为浏览器数据目录
4. 多个账号同时运行时，需要使用不同的端口
5. 日志会同时输出到控制台和日志文件
6. 每次服务重启时，日志文件会自动清空
