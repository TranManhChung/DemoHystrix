package main

import "github.com/TranManhChung/DemoHystrix/cmd"

var revision = ""

// ./main admin --config=zpi-e-voucher.local.toml
func main() {
	cmd.SetRevision(revision)
	cmd.Execute()
}
