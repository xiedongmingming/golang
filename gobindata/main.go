package gobindata


// go get -u github.com/jteeuwen/go-bindata/...
// 最后的三个点不可少:
// 会分析所有子目录并下载依赖编译子目录内容
// GO-BINDATA的命令工具在子目录中

// 使用示例:
// go-bindata -o=app/asset/asset.go -pkg=asset source/... theme/... doc/source/... doc/theme/...
// -o: 输出文件到
// -pkg=asset 包名
// 然后是需要打包的目录(三个点包括所有子目录)这样就可以把所有相关文件打包到ASSET.GO且开头是PACKAGE ASSET保持和目录一致