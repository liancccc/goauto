## 命令模块

### 执行命令

请求方法：GET

请求接口：/exec 

请求参数：command

请求成功响应：

```json
{"data":null,"message":"success","success":true}
```

### 获取命令模块信息

请求方法：GET

请求接口：/execHelp

请求成功响应：

```json
{
    "data": {
        "flows": [
            "DomainALL"
        ],
        "help": "\n  _________  ___  __  ____________ \n / ___/ __ \\/ _ |/ / / /_  __/ __ \\\n/ (_ / /_/ / __ / /_/ / / / / /_/ /\n\\___/\\____/_/ |_\\____/ /_/  \\____/\n\n\t\tgithub.com/liancccc/goauto\n\nUsage:\n  goauto scan [flags]\n\nFlags:\n      --flow string        \n  -h, --help               help for scan\n      --log-file           \n      --target string      \n      --task-name string\n\nGlobal Flags:\n      --debug\n"
    },
    "message": "success",
    "success": true
}
```

此处的 flows 替换掉前端的模块部分，变更为工作流

### 批量上传

请求方法：POST

请求接口：/upload/targets

请求参数：targets

请求成功响应：

```json
{
    "data": {
        "filePath": "C:\\Users\\admin\\goauto-base\\targets\\20250911185929.txt"
    },
    "message": "upload success",
    "success": true
}
```

## 系统信息

请求方法：GET

请求接口：/system/info

请求成功响应：

```json
{
    "data": {
        "hostname": "DESKTOP-M8UVIRV",
        "cpu_usage": 7.293952353230934,
        "memory_usage": 53,
        "memory_total": 27.69232940673828,
        "memory_used": 14.733139038085938,
        "os": "Microsoft Windows 11 Pro 10.0.26100.6584 Build 26100.6584",
        "arch": "amd64",
        "go_version": "go1.25.1"
    },
    "message": "success",
    "success": true
}
```

## 任务信息

### 任务状态

请求方法：GET

请求接口：/task/status

请求成功响应：

```json
{
    "data": {
        "running": 3,
        "waiting": 2
    },
    "message": "success",
    "success": true
}
```

### 任务列表

请求方法：GET

请求接口：/task/list

请求成功响应：

```
{
    "data": [
        {
            "task": "xazlsec.com",
            "flow": "DomainALL",
            "start_at": "2025-09-11 19:16:22",
            "end_at": "",
            "pid": 38004,
            "status": "running",
            "command": "C:\\Users\\admin\\AppData\\Local\\go-build\\c4\\c4648469ec0667a966dbc3180db76070d2db6564a567101af80994327deab3e1-d\\main.exe scan --target xazlsec.com --log-file --flow DomainALL"
        }
    ],
    "message": "success",
    "success": true
}
```

之前的运行模式现在更换为工作流

### 任务详情

请求方法：GET

请求接口：/task/detail

请求参数：task

请求成功响应：

```json
{
    "data": {
        "task": "xazlsec.com",
        "flow": "DomainALL",
        "start_at": "2025-09-11 19:16:22",
        "end_at": "",
        "pid": 38004,
        "status": "running",
        "command": "C:\\Users\\admin\\AppData\\Local\\go-build\\c4\\c4648469ec0667a966dbc3180db76070d2db6564a567101af80994327deab3e1-d\\main.exe scan --target xazlsec.com --log-file --flow DomainALL",
        "child_process": [
            {
                "pid": 31368,
                "command": "powershell -Command \"ksubdomain enum -d xazlsec.com -o C:\\Users\\admin\\goauto-workspace\\xazlsec.com\\subdomains\\ksubdomain-enum.txt\"",
                "create_at": "2025-09-11 19:17:59"
            }
        ],
        "reports": {
            "subdomains": [
                {
                    "name": "ksubdomain-enum.txt",
                    "count": 0,
                    "link": "workspace/xazlsec.com/subdomains/ksubdomain-enum.txt"
                },
                {
                    "name": "oneforall.txt",
                    "count": 15,
                    "link": "workspace/xazlsec.com/subdomains/oneforall.txt"
                },
                {
                    "name": "subfinder.txt",
                    "count": 10,
                    "link": "workspace/xazlsec.com/subdomains/subfinder.txt"
                }
            ]
        }
    },
    "message": "success",
    "success": true
}
```

