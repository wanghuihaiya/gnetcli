package fsm

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestName(t *testing.T) {
	// 创建临时模板目录
	tmpDir, err := ioutil.TempDir("", "textfsm_demo")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 写入示例模板文件（模拟 display interface brief 命令）
	sampleTpl := `Value IFACE (\\S+)
Value STATUS (\\S+)
Start
^${IFACE}\\s+${STATUS} -> Record
`
	tplPath := filepath.Join(tmpDir, "huawei_display_interface_brief.textfsm")
	if err := ioutil.WriteFile(tplPath, []byte(sampleTpl), 0644); err != nil {
		log.Fatalf("写入模板失败: %v", err)
	}

	// 模拟设备输出
	sampleOutput := `
GE1/0/1 up
GE1/0/2 down
GE1/0/3 up
`

	// 初始化模板管理器
	tm, err := NewTemplateManager(128)
	if err != nil {
		log.Fatalf("初始化模板管理器失败: %v", err)
	}
	if err := tm.BuildIndexFromDir(tmpDir); err != nil {
		log.Fatalf("构建模板索引失败: %v", err)
	}

	// 执行解析
	pe := NewParserEngine(tm)
	rows, err := pe.ParseByCommandAndVersion("huawei", "S5720-28X-SI-AC", "V200R005", "display interface brief", sampleOutput)
	if err != nil {
		log.Fatalf("解析出错: %v", err)
	}

	fmt.Println("✅ 解析成功！结果如下：")
	PrintResult(rows)
}
