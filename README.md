# miniblog 项目

## 编译命令

```sh
# gofmt 用来格式化当前目录及其子目录中的所有 Go 源文件
$ gofmt -s -w ./
# go build 命令用来编译Go源码
# -o 指定输出的可执行文件名称
# -v 选项用于显示详细的编译过程信息
$ go build -o output/mb-apiserver -v cmd/mb-apiserver/main.go # 运行编译生成的可执行文件
$output/mb-apiserverHello World!
```
