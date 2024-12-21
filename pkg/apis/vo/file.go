package vo

// FileInfo 上传得文件信息
type FileInfo struct {
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	Src      string `json:"src"`
	FileType string `json:"fileType"`
}
