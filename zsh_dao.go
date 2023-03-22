package penv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-ps"
)

var (
	zshShell = &shell{
		configFileName: filepath.Join(os.Getenv("HOME"), ".zshrc"),
		commentSigil:   " #",
		quote: func(value string) string {
			r := strings.NewReplacer(
				"\\", "\\\\",
				"'", "\\'",
				"\n", `'"\n"'`,
				"\r", `'"\r"'`,
			)
			return "'" + r.Replace(value) + "'"
		},
		mkSet: func(sh *shell, nv NameValue) string {
			return fmt.Sprintf(
				"export %s=%s",
				nv.Name, sh.quote(nv.Value),
			)
		},
		mkAppend: func(sh *shell, nv NameValue) string {
			return fmt.Sprintf(
				"export %s=${%s}${%s:+:}%s",
				nv.Name, nv.Name, nv.Name, sh.quote(nv.Value),
			)
		},
		mkUnset: func(sh *shell, nv NameValue) string {
			return fmt.Sprintf(
				"unset %s",
				nv.Name,
			)
		},
	}
)

type (
	zshOp struct {
		op        string
		nameValue NameValue
	}
	// ZshDAO is a data access object for zsh
	ZshDAO struct{}
)

func init() {
	RegisterDAO(1000, func() bool {
		pid := os.Getpid()
		for pid > 0 {
			p, err := ps.FindProcess(pid)
			if err != nil || p == nil {
				break
			}
			if p.Executable() == "zsh" {
				return true
			}
			pid = p.PPid()
		}
		return false
	}, zshShell)
}
