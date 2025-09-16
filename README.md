## 介绍

安全工具工作流

2023 年完成过一版 *GOAUTO* 但是由于工具需要手动安装，添加工具复杂，所以弃用。在这期间学习阅读各种工具源码以期待将整个流程使用单纯的 *Go* 来实现，但是个人的力量是有限的，工作后更没有精力和信心把每个模块的工具都做到很好，实现功能不难，参考已有工具就可以，但是仅实现功能好像并没有任何的意义。而且新的更好的工具也在不断出现，有开源也有闭源，如果有新的工具新的思路出现，凭个人去维护是很难的，最近空闲又拾起这件事情。改变想法，单纯的调用工具来完成整个流程。

[Osmedeus](https://github.com/j3ssie/osmedeus) 可以说是调用二进制工具实现工作流一个很好的工具，可以通过 `yaml` 来实现各种模块的工作流，由代码提供自定义 Script 配合系统命令实现一个更自由的工具流，所有参数可控，命令可控，但是由于基于 `yaml` 编写太不习惯，一些细节不好把控，所以再编写一版简单的 *GOAUTO*。相较于 osm 会显得很死板，不能动态的去改变工作流，但是可以做更多的事情，需要个人去按照个人的习惯去编写。

## 安装

### 工具结构

![image-20250911162519370](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911162519370.png)

### 所需环境

- nmap
- libpcap
- git
- python3 | pip3
- golang
- xscan

命令执行使用 powershell 和 bash，测试于 windows11 和 ubuntu22。

其余工具会通过 git、go 和下载可执行文件的方式自动下载。

一些 VPS 的环境安装 ksubdomain、naabu 这些会报错，可以尝试一下命令：

```
# 选 18 然后重启
apt-get install build-essential -y
apt-get install libpcap-dev -y 
```

如果 httpx 截图报 rad 缺少依赖的话：

```
apt install -y libnss3 libatk1.0-0 libatk-bridge2.0-0 libcups2 libxss1 libxcomposite1 libxrandr2 libasound2 libpangocairo-1.0-0 libgtk-3-0

wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
dpkg -i google-chrome*; sudo apt-get -f install
apt-get install google-chrome-stable
```

xscan 加 8 佬星球获取：https://t.zsxq.com/L2U6k

```
unzip xscan_*.zip -d /root/goauto-base/tools/xscan
```

Golang 安装：

```
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz -O go.linux-amd64.tar.gz
tar -C /usr/local -xzf go.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin:/root/go/bin" >> ~/.bash_profile
echo "export PATH=$PATH:/usr/local/go/bin:/root/go/bin" >> ~/.bashrc
source ~/.bash_profile
source ~/.bashrc
bash -c go
```

### 初始化

```
go install -v github.com/liancccc/goauto@latest
goauto install
```

一些工具会报毒比如 xray 如果被删就手动下载加到白名单，其中有几个工具是 fork 到自己的仓库，如果介意就自己下载：

- oneforall：注释爆破相关，因为安装 pip 会报错
- cdncheck：添加国内源和并发扫描
- uro：windows 环境下 -o 的 gbk 问题，取消 pipx 变更单纯的 py

需要人工配置的：

- subfinder
- notify
- uncover

初始化包含 3 个部分：

- 虚拟环境
- 模块 各种工具的安装
- 字典

![image-20250911160604343](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911160604343.png)

![image-20250911162903696](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911162903696.png)

### 通知配置

Goauto 的理念是将各种优秀工具功能作为模块，然后附加到工作流中。

通知同样是借助工具实现，后续如果想添加其他功能也推荐独立成工具。

使用项目：https://github.com/projectdiscovery/notify

支持各种通常如 tg、smtp 等等，默认工作流通知的是漏扫结果

## 使用

### 列出工作流

```
goauto flows
```

这只是一个默认的工作流，更推荐大家自己按照自己的方式来编写，自定义某种模块功能然后进行编排，后续默认工作流不会做出大的更改。

默认工作流：

- 子域名收集
  - subfinder
  - ksubdomain
  - alterx
- CDN 识别 [cdncheck]
- 非 CDN 端口扫描 [ dnsx naabu]
- 根据域名从 Quake 拉取资产 [uncover]
- 去重验活 [httpx]
  - CDN 域名使用进行去重验活 [ 常见的 WEB 端口 ]
  - 非 CDN 使用进行去重验活
  - 去重后的 URL 获取指纹和截图
- 爬虫
  - gospider
  - urlfinder
  - katana
- 漏洞扫描
  - xscan
  - xray
  - nuclei
- 漏扫通知 [notify]

![image-20250911204623102](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911204623102.png)

### 扫描模式

```
goauto scan -h

goauto scan --target vulnweb.com --flow DomainALL --debug
```

![image-20250911204723041](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911204723041.png)

### WEB模式

主要参考的也是 osm ，osm 的 WEB 会把每个模块的报告展示出来，直接去查看原始的输出，这点是之前没有想到过的。这就更轻量级了，不需要什么数据库什么展示，只需要原原本本的展示就可以了，毕竟如果添加了数据库，那就会把这件事情变得复杂。

默认 8848 钛合金端口，账号密码不指定就随机

前端是基于 daisyui 写的，可以参考 https://daisyui.com/docs/editor/cursor/ 用 Cursor 自己改。

```
goauto web
```

![image-20250911204806791](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911204806791.png)

后台提供命令执行的功能方便自动化下发任务，没有做过滤，没啥必要，而且我有执行其他命令的需求。所以认证自己不要设置的太弱，或者就是默认的随机。

## 如何添加工具和自定义工作流?

### internal/venv

自定义虚拟环境或者是其他环境比如 chrome 的安装和获取路径也算。

这个没有特定的接口样式，按照自己的实现可以初始化调用即可。

### internal/dict

管理字典

- 名称、路径、下载链接

```go
var Dicts = map[string]Dict{
	"subdomain-all": {
		Name: "Subdomain-ALL",
		Path: filepath.Join(paths.DictDir, "subdomain-all.txt"),
		Link: "https://gist.githubusercontent.com/jhaddix/f64c97d0863a78454e44c2f7119c2a6a/raw",
	},
}
```

运行 install 后会自动下载不存在的，调用可以直接导这个 Dicts 的 map 来获取路径。

### internal/modules

module 表示一种功能，而不是工具，比如 ksubdomain 有两种功能，那么就是两个。

接口如下：

```go
type Module interface {
	Name() string
	Install() error
	CheckInstalled() bool
	Run(funcParams any)
}
```

按照这个格式写即可：

```go
package subfinder

import (
	"fmt"
	"strings"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/projectdiscovery/gologger"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "subfinder"
}

func (m *ModuleStruct) Install() error {
	var installCmd = "go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"
	_, err := executil.RunCommandSteamOutput(installCmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("subfinder")
	return strings.Contains(commandSteamOutput, "projectdiscovery.io")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var command string
	if params.IsFileTarget() {
		command = fmt.Sprintf("subfinder -dL %s -o %s", params.Target, params.Output)
	} else {
		command = fmt.Sprintf("subfinder -d %s -o %s", params.Target, params.Output)
	}
	_, err := executil.RunCommandSteamOutput(command, params.Timeout)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}

```

也可以参考 `internal/modules/xray/xray.go` 下载压缩包实现按照。

`modules.BaseParams` 表示基础参数的结构体，如果是其他的再创建对应的结构体即可，不过这个已经可以覆盖大部分的了。

```go
type BaseParams struct {
	Target string
	Output string

	Proxy   string
	Dict    string
	Timeout string

	CustomizeParams string
}
```

对于代理这些添加就是这样，每个支持代理的工具的参数都不同，所有需要这样：

![image-20250911162156783](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911162156783.png)

编写完成后需要在 `internal/initializer/module.go` 中导入一下：

![image-20250911162418914](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911162418914.png)

### internal/workflow

工作流，把各种功能模块串起来运行，按照文件的方式传递。

示例 DomainALL 工作流 `internal/workflow/domain.go`

```go
package workflow

import (
	"path/filepath"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/alterx"
	"github.com/liancccc/goauto/internal/modules/cdncheck"
	"github.com/liancccc/goauto/internal/modules/gospider"
	httpx_info "github.com/liancccc/goauto/internal/modules/httpx/info"
	httpx_unique "github.com/liancccc/goauto/internal/modules/httpx/unique"
	"github.com/liancccc/goauto/internal/modules/katana"
	"github.com/liancccc/goauto/internal/modules/ksubdomain/enum"
	"github.com/liancccc/goauto/internal/modules/ksubdomain/verify"
	"github.com/liancccc/goauto/internal/modules/merge"
	"github.com/liancccc/goauto/internal/modules/naabu"
	"github.com/liancccc/goauto/internal/modules/nuclei"
	"github.com/liancccc/goauto/internal/modules/oneforall"
	"github.com/liancccc/goauto/internal/modules/subfinder"
	"github.com/liancccc/goauto/internal/modules/unique"
	"github.com/liancccc/goauto/internal/modules/urlfinder"
	"github.com/liancccc/goauto/internal/modules/xray"
	xscan_spider "github.com/liancccc/goauto/internal/modules/xscan/spider"
	"github.com/projectdiscovery/gologger"
)

type DomainALLFlow struct {
}

func init() {
	RegisterWorkflow(&DomainALLFlow{})
}

func (f *DomainALLFlow) Name() string {
	return "DomainALL"
}

func (f *DomainALLFlow) Description() string {
	return "子域名收集 -> CDN识别 -> 端口扫描 -> WEB验活去重 -> WEB获取信息和截图 -> 爬虫 -> 漏洞扫描"
}

func (f *DomainALLFlow) Run(runner *Runner) {
	var subdomainOutDir = filepath.Join(runner.workSpace, "subdomains")
	new(subfinder.ModuleStruct).Run(modules.BaseParams{
		Target: runner.opt.Target,
		Output: filepath.Join(subdomainOutDir, "subfinder.txt"),
	})
	new(oneforall.ModuleStruct).Run(modules.BaseParams{
		Target: runner.opt.Target,
		Output: filepath.Join(subdomainOutDir, "oneforall.txt"),
	})
	new(ksubdomain_enum.ModuleStruct).Run(modules.BaseParams{
		Target:  runner.opt.Target,
		Output:  filepath.Join(subdomainOutDir, "ksubdomain.txt"),
		Timeout: "5m",
	})
	new(merge.ModuleStruct).Run(merge.Params{
		BaseParams: &modules.BaseParams{
			Output: filepath.Join(subdomainOutDir, "merge.txt"),
		},
		Targets: []string{filepath.Join(subdomainOutDir, "subfinder.txt"), filepath.Join(subdomainOutDir, "oneforall.txt"), filepath.Join(subdomainOutDir, "ksubdomain.txt")},
	})
	new(alterx.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(subdomainOutDir, "merge.txt"),
		Output: filepath.Join(subdomainOutDir, "alterx.txt"),
	})

	new(ksubdomain_verify.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(subdomainOutDir, "alterx.txt"),
		Output: filepath.Join(subdomainOutDir, "alterx-alive.txt"),
	})
	new(merge.ModuleStruct).Run(merge.Params{
		BaseParams: &modules.BaseParams{
			Output: filepath.Join(subdomainOutDir, "merge-alterx.txt"),
		},
		Targets: []string{filepath.Join(subdomainOutDir, "merge.txt"), filepath.Join(subdomainOutDir, "alterx-alive.txt")},
	})

	new(unique.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(subdomainOutDir, "merge-alterx.txt"),
		Output: filepath.Join(subdomainOutDir, "all.txt"),
	})

	if fileutil.CountLines(filepath.Join(subdomainOutDir, "all.txt")) == 0 {
		gologger.Error().Msgf("%s Get 0 Subdomains Exit", runner.opt.Target)
		return
	}

	fileutil.Cleaning(subdomainOutDir, []string{
		filepath.Join(subdomainOutDir, "all.txt"),
	})

	// CDN 识别
	var cdncheckOutDir = filepath.Join(runner.workSpace, "cdncheck")
	new(cdncheck.ModuleStruct).Run(cdncheck.Params{
		BaseParams: &modules.BaseParams{
			Target: filepath.Join(subdomainOutDir, "all.txt"),
		},
		CDNPath:   filepath.Join(cdncheckOutDir, "cdn.txt"),
		NoCDNPath: filepath.Join(cdncheckOutDir, "noCdn.txt"),
	})

	// ....
}
```

## 记录



## 待办

- [ ] 跑一个月看 BUG
- [ ] 加爬虫！gau craw 动态爬虫！！！！！！！！！！！！！！
- [x] 加 wlh
- [ ] nuclei 报告 xray 展示 等级分级
- [ ] httpx 的信息不要这样了 去重后输出的仍然是 json，然后解析渲染变成 html
- [ ] quake 获取 services uncover 不太行
- [ ] 路径整理 拆分
- [ ] 自己搞一个监控 SRC 的
- [ ] xray 报告拆分
- [x] alterx 超时控制
- [ ] 结果判断 是否跑偏
- [ ] 工作流重构
