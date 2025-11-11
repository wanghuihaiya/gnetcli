package fsm

import (
	"fmt"
	"github.com/sirikothe/gotextfsm"
	"os"
	"path/filepath"
	"strings"
	"sync"

	lru "github.com/hashicorp/golang-lru"
)

// TemplateMeta 表示模板的元信息
// 包括厂商、命令和模板路径
type TemplateMeta struct {
	Vendor  string // 厂商名称
	Model   string // 型号
	Version string // 版本
	Command string // 命令名称（已标准化）
	Path    string // 模板文件路径
}

// TemplateManager 管理模板索引和缓存
type TemplateManager struct {
	mu    sync.RWMutex
	index map[string]map[string]map[string]map[string]*TemplateMeta // vendor -> model -> version -> command -> meta
	cache *lru.Cache                                                // path -> *gotextfsm.TextFSM
}

// NewTemplateManager 创建模板管理器并初始化缓存
func NewTemplateManager(cacheSize int) (*TemplateManager, error) {
	c, err := lru.New(cacheSize)
	if err != nil {
		return nil, err
	}
	return &TemplateManager{
		index: make(map[string]map[string]map[string]map[string]*TemplateMeta),
		cache: c,
	}, nil
}

// BuildIndexFromDir 扫描目录建立索引，支持默认模型和版本
// 文件命名示例: vendor_command_model_version.textfsm
func (tm *TemplateManager) BuildIndexFromDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.textfsm"))
	if err != nil {
		return err
	}
	for _, f := range files {
		base := filepath.Base(f)
		parts := strings.SplitN(base, "_", 4)
		if len(parts) != 4 {
			continue // 不符合规则的文件
		}

		vendor := strings.ToLower(parts[0])
		model := strings.ToLower(parts[1])
		version := strings.ToLower(parts[2])
		command := strings.ReplaceAll(parts[3], "_", " ")
		command = strings.ToLower(strings.TrimSuffix(command, ".textfsm"))

		// 建立索引
		if tm.index[vendor] == nil {
			tm.index[vendor] = make(map[string]map[string]map[string]*TemplateMeta)
		}
		if tm.index[vendor][model] == nil {
			tm.index[vendor][model] = make(map[string]map[string]*TemplateMeta)
		}
		if tm.index[vendor][model][version] == nil {
			tm.index[vendor][model][version] = make(map[string]*TemplateMeta)
		}
		meta := &TemplateMeta{Vendor: vendor, Model: model, Version: version, Command: command, Path: f}
		tm.index[vendor][model][version][command] = meta
	}
	return nil
}

// FindTemplate 查找模板，优先匹配厂商+型号+版本+命令，支持默认模板回退
func (tm *TemplateManager) FindTemplate(vendor, model, version, command string) (*TemplateMeta, error) {
	vendor = strings.ToLower(strings.TrimSpace(vendor))
	model = strings.ToLower(strings.TrimSpace(model))
	version = strings.ToLower(strings.TrimSpace(version))
	command = strings.ToLower(strings.Join(strings.Fields(command), " "))

	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// 匹配优先级
	variants := []struct{ mdl, ver string }{
		{model, version}, // 完整匹配
		{model, "*"},     // 模型 + 默认版本
		{"*", "*"},       // 默认模板
	}

	if models, ok := tm.index[vendor]; ok {
		for _, v := range variants {
			if versions, ok := models[v.mdl]; ok {
				if cmds, ok := versions[v.ver]; ok {
					if meta, ok := cmds[command]; ok {
						return meta, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("未找到模板: vendor=%s, model=%s, version=%s, command=%s", vendor, model, version, command)
}

// GetFSM 获取已编译的 TextFSM，对路径做缓存和懒加载
func (tm *TemplateManager) GetFSM(path string) (*gotextfsm.TextFSM, error) {
	// 检查缓存
	if v, ok := tm.cache.Get(path); ok {
		return v.(*gotextfsm.TextFSM), nil
	}

	// 读取文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取模板失败 %s: %w", path, err)
	}

	// 编译模板
	var fsm gotextfsm.TextFSM
	if err := fsm.ParseString(string(data)); err != nil {
		return nil, fmt.Errorf("解析模板失败 %s: %w", path, err)
	}

	// 添加到缓存
	tm.cache.Add(path, &fsm)
	return &fsm, nil
}
