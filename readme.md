# monitor
> golang 监听热编译，用于golang或nodejs

### 安装 
```
go get Zhan9Yunhua/monitor
```

### 使用
```
monitor -lang=xx -o=xx -f=xx -args=xxx -s=xxx
```

- -lang =
语言,默认go（go||js||ts）
- -o =
打包路径（默认当前路径，node不需要）
- -f =
配置文件（默认不需要，"./monitor.json"）
- -args =
命令行参数（"-host=:8080"）
- -s =
只有node需要。每次文件改变执行的脚本（ts: ts-node test.ts）

### 配置文件
```
{
  // 应用名称
  "appName": "app",
  // 打包路径（node不需要）
  "output": "",
  // 监听文件后缀。默认.go
  "exts": [".go"],
  // 打包-tags（node不需要）
  "buildTags": "",
  // 命令行args（键值对）
  "cmdArgs": ["a=b"],
  // 环境变量（键值对）
  "envs": ["a=b"],
  // 语言,默认go（go||js||ts）
  "lang": "go",
  // 只有node需要。每次文件改变执行的脚本（ts: ts-node test.ts）
  "script": "node test.js"
}

```
