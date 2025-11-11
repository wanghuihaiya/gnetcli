package fsm

import (
	"encoding/json"
	"fmt"

	"github.com/sirikothe/gotextfsm"
)

// ParserEngine 封装模板管理器并负责执行解析
type ParserEngine struct {
	TM *TemplateManager
}

// NewParserEngine 创建解析引擎
func NewParserEngine(tm *TemplateManager) *ParserEngine {
	return &ParserEngine{TM: tm}
}

// ParseByCommandAndVersion 支持厂商+型号+版本+命令解析
func (pe *ParserEngine) ParseByCommandAndVersion(vendor, model, version, command, input string) ([]map[string]interface{}, error) {
	meta, err := pe.TM.FindTemplate(vendor, model, version, command)
	if err != nil {
		return nil, err
	}
	fsm, err := pe.TM.GetFSM(meta.Path)
	if err != nil {
		return nil, err
	}
	var parser gotextfsm.ParserOutput
	if err := parser.ParseTextString(input, *fsm, true); err != nil {
		return nil, fmt.Errorf("解析输入失败: %w", err)
	}
	return parser.Dict, nil
}

// PrintResult 打印解析结果
func PrintResult(rows []map[string]interface{}) {
	b, _ := json.MarshalIndent(rows, "", " ")
	fmt.Printf("解析结果:\n%s\n", string(b))
}
