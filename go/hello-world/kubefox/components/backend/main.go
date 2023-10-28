package main

import (
	"embed"

	"github.com/xigxog/kubefox/libs/core/kit"
)

//go:embed 1mb.txt
var EFS embed.FS

var file []byte

func main() {
	k := kit.New()

	f, err := EFS.ReadFile("1mb.txt")
	if err != nil {
		k.Log().Fatal(err)
	}
	file = f

	k.Default(sayWho)
	k.Start()
}

func sayWho(k kit.Kontext) error {
	who := k.EnvDef("who", "World")
	k.Log().Debugf("The who is '%s'!", who)

	return k.Resp().SendStr(who)
}
