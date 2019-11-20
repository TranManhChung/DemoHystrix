package main

import "gitlab.360live.vn/zalopay/zpi-e-voucher/cmd"

var revision = ""

func main() {
	cmd.SetRevision(revision)
	cmd.Execute()
}
