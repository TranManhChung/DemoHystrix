package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.360live.vn/zalopay/go-common/log"
	"gitlab.360live.vn/zalopay/go-common/tracing"
	"gitlab.360live.vn/zalopay/zpi-e-voucher/services/voucher"
)

var voucherCmd = &cobra.Command{
	Use:   "voucher",
	Short: "Starts EVoucher service",
	Long:  `Starts EVoucher service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.NewStandardFactory(viper.GetString("voucher.logs"), "voucher")

		server := voucher.NewServer(
			tracing.Init("voucher.EVoucherService", logger),
			logger,
		)

		return server.Run()
	},
}

func init() {
	RootCmd.AddCommand(voucherCmd)
}
