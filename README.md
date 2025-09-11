

## 简介

*GOAUTO* 的目的是解放渗透测试过程中需要使用多种工具，将各种工具的安装和使用集成于一体，解放双手。

2023 年完成过一版 *GOAUTO* 但是由于工具需要手动安装，添加工具复杂，所以弃用。在这期间学习阅读各种工具源码以期待将整个流程使用单纯的 *Go* 来实现，但是个人的力量是有限的，工作后更没有精力和信心把每个模块的工具都做到很好，实现功能不难，参考已有工具就可以，但是仅实现功能好像并没有任何的意义。而且新的更好的工具也在不断出现，有开源也有闭源，如果有新的工具新的思路出现，凭个人去维护是很难的，最近空闲又拾起这件事情。改变想法，单纯的调用工具来完成整个流程。

[Osmedeus](https://github.com/j3ssie/osmedeus) 可以说是调用二进制工具实现工作流一个很好的工具，可以通过 `yaml` 来实现各种模块的工作流，由代码提供自定义 Script 配合系统命令实现一个更自由的工具流，所有参数可控，命令可控，但是由于基于 `yaml` 编写太不习惯，一些细节不好把控，所以再编写一版简单的 *GOAUTO*。

## 初始化

![image-20250911162519370](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911162519370.png)

前提所需环境：

- nmap
- libpcap
- git
- python3 | pip3
- golang
- xscan

命令执行使用 powershell 和 bash。

其余工具会通过 git、go 和下载可执行文件的方式自动下载。

一些 VPS 的环境安装 ksubdomain、naabu 这些会报错，可以尝试一下命令：

```
apt-get install build-essential
apt-get install libpcap-dev # 选 18 然后重启
```

Golang 安装：

```
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz -O go.linux-amd64.tar.gz
tar -C /usr/local -xzf go.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bash_profile
echo "export PATH=$PATH:/root/go/bin" >> ~/.bash_profile
source ~/.bash_profile
bash -c go
```

工具初始化：

```
go install -v github.com/liancccc/goauto@latest
goauto install
```

初始化包含 3 个部分：

- 虚拟环境
- 模块 各种工具的安装
- 字典

![image-20250911160604343](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911160604343.png)

![image-20250911162903696](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911162903696.png)

## 工具使用

### 列出工作流

```
goauto flows
```

![image-20250911204623102](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911204623102.png)

### 扫描模式

```
goauto scan -h
```

![image-20250911204723041](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911204723041.png)

### WEB模式

```
goauto web
```

![image-20250911204806791](https://blog-1310215391.cos.ap-beijing.myqcloud.com/images/image-20250911204806791.png)



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

	// 对非 CDN 目标进行端口扫描
	var portscanOutDir = filepath.Join(runner.workSpace, "portscan")
	if fileutil.CountLines(filepath.Join(cdncheckOutDir, "noCdn.txt")) > 0 {
		new(naabu.ModuleStruct).Run(modules.BaseParams{
			Target: filepath.Join(cdncheckOutDir, "noCdn.txt"),
			Output: filepath.Join(portscanOutDir, "noCdn-services.txt"),
		})
	}

	// URL 去重验活 + 获取截图等信息
	var httpxOutDir = filepath.Join(runner.workSpace, "httpx")
	if fileutil.CountLines(filepath.Join(cdncheckOutDir, "cdn.txt")) > 0 {
		new(httpx_unique.ModuleStruct).Run(modules.BaseParams{
			Target:          filepath.Join(cdncheckOutDir, "cdn.txt"),
			Output:          filepath.Join(httpxOutDir, "cdn-alive.txt"),
			CustomizeParams: "-mc 200,302 -p 80,443,8080,8000,8888,4848,7070,8089,8181,9080,9443,5000,8443,5001,81,8081,50805,3000,88,7547",
			Proxy:           runner.opt.Proxy,
		})
	}
	if fileutil.CountLines(filepath.Join(portscanOutDir, "noCdn-services.txt")) > 0 {
		new(httpx_unique.ModuleStruct).Run(modules.BaseParams{
			Target: filepath.Join(portscanOutDir, "noCdn-services.txt"),
			Output: filepath.Join(httpxOutDir, "noCdn-alive.txt"),
			Proxy:  runner.opt.Proxy,
		})
	}

	new(merge.ModuleStruct).Run(merge.Params{
		BaseParams: &modules.BaseParams{
			Output: filepath.Join(httpxOutDir, "merge.txt"),
		},
		Targets: []string{filepath.Join(httpxOutDir, "noCdn-alive.txt"), filepath.Join(httpxOutDir, "cdn-alive.txt")},
	})

	new(unique.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "merge.txt"),
		Output: filepath.Join(httpxOutDir, "all.txt"),
	})

	new(httpx_info.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(httpxOutDir, "web.txt"),
	})

	// 爬虫
	var spiderOutDir = filepath.Join(runner.workSpace, "spider")
	new(gospider.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(spiderOutDir, "gospider.txt"),
		Proxy:  runner.opt.Proxy,
	})
	new(katana.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(spiderOutDir, "katana.txt"),
		Proxy:  runner.opt.Proxy,
	})
	new(urlfinder.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(spiderOutDir, "urlfinder.txt"),
		Proxy:  runner.opt.Proxy,
	})
	new(merge.ModuleStruct).Run(merge.Params{
		BaseParams: &modules.BaseParams{
			Output: filepath.Join(spiderOutDir, "all.txt"),
		},
		Targets: []string{filepath.Join(spiderOutDir, "gospider.txt"), filepath.Join(spiderOutDir, "katana.txt"), filepath.Join(spiderOutDir, "urlfinder.txt")},
	})

	new(unique.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(spiderOutDir, "all.txt"),
		Output: filepath.Join(spiderOutDir, "links.txt"),
	})

	// 漏洞扫描
	var vulscanOutDir = filepath.Join(runner.workSpace, "vulscan")
	new(nuclei.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(vulscanOutDir, "nuclei.txt"),
	})
	new(xscan_spider.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(spiderOutDir, "links.txt"),
		Output: filepath.Join(vulscanOutDir, "xscan.json"),
	})
	new(xray.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(spiderOutDir, "links.txt"),
		Output: filepath.Join(vulscanOutDir, "xray.txt"),
	})
}
```

