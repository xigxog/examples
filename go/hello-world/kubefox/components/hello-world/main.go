package main

import (
	"fmt"

	"github.com/xigxog/kubefox/libs/core/kit"
)

func main() {
	k := kit.New()
	k.Route("Path(`/examples/hello-world`)", hello)
	k.Start()
}

func hello(ktx kit.Kontext) error {
	msg := fmt.Sprintf("👋 Hello %s!", ktx.EnvDef("who", "World"))
	ktx.Log().Info(msg)

	return ktx.Resp().SendStr(msg)
}
