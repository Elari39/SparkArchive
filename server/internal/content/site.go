package content

// LoadSite 只读取站点配置，供运行时动态加载文案使用。
func LoadSite(path string) (Site, error) {
	return readYAML[Site](path)
}
