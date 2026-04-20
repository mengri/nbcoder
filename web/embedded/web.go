package embedded

import "embed"

//go:embed dist
var FS embed.FS

// WebFS 获取嵌入的文件系统
func WebFS() *embed.FS {
	return &FS
}
