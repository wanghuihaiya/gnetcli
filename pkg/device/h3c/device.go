/*
Package h3c implements huawei CLI using genericcli.
It is a copy of huawei device with small changes
*/
package h3c

import (
	"fmt"
	"regexp"

	"github.com/annetutil/gnetcli/pkg/cmd"
	"github.com/annetutil/gnetcli/pkg/device/genericcli"
	"github.com/annetutil/gnetcli/pkg/expr"
	"github.com/annetutil/gnetcli/pkg/streamer"
)

const (
	questionExpression = `\n(?P<question>.*Continue\? \[Y/N\]:)$`
	// User-view <hostname> and system-view [hostname]. Do not use [<\[]...[>\]] —
	// that treats a MAC Port/Nickname like [CORE-SW] as a prompt and stops the read early.
	promptExpression = `(\r\n|^|\x00)(?P<prompt>((\(M\))?<[/\w\-.:]+>|(\(M\))?\[[~*]?[/\w\-.:]+\]))\s*$`
	errorExpression  = `(` +
		`(\^\r\n)?( % )?Error:(?P<error>.+) at '\^' position\.` +
		`|\r\n % (Unrecognized command|Too many parameters|Incomplete command) found at '\^' position\.` +
		`)`
	pagerExpression = `(?P<store>(\r\n|\n))?(?:\x1b\[[0-9;]*m)*\s*-{0,4}\s*More\s*-{0,4}(?:\x1b\[[0-9;]*m)*\s*$`
)

var autoCommands = []cmd.Cmd{
	cmd.NewCmd("screen-length 0 temporary", cmd.WithErrorIgnore()),
	cmd.NewCmd("screen-length disable", cmd.WithErrorIgnore()),
	cmd.NewCmd("terminal echo-mode line", cmd.WithErrorIgnore()),
	cmd.NewCmd("undo terminal monitor", cmd.WithErrorIgnore()),
}

func NewDevice(connector streamer.Connector, opts ...genericcli.GenericDeviceOption) genericcli.GenericDevice {
	cli := genericcli.MakeGenericCLI(
		expr.NewSimpleExprLast200().FromPattern(promptExpression),
		expr.NewSimpleExprLast200().FromPattern(errorExpression),
		genericcli.WithPager(
			expr.NewSimpleExprLast200().FromPattern(pagerExpression),
		),
		genericcli.WithAutoCommands(autoCommands),
		genericcli.WithQuestion(
			expr.NewSimpleExprLast200().FromPattern(questionExpression),
		),
		genericcli.WithSFTPEnabled(),
		genericcli.WithTerminalParams(400, 0),
		// h3c adds extra \r in the echo
		genericcli.WithEchoExprFn(func(c cmd.Cmd) expr.Expr {
			return expr.NewSimpleExpr().FromPattern(fmt.Sprintf(`%s\r*\n`, regexp.QuoteMeta(string(c.Value()))))
		}),
	)
	return genericcli.MakeGenericDevice(cli, connector, opts...)
}
