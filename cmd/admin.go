package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.360live.vn/zalopay/go-common/log"
	"gitlab.360live.vn/zalopay/go-common/tracing"
	"gitlab.360live.vn/zalopay/zpi-e-voucher/services/admin"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Starts EVoucher service",
	Long:  `Starts EVoucher service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.NewStandardFactory(viper.GetString("admin.logs"), "admin")

		server := admin.NewServer(
			tracing.Init("admin.EVoucherService", logger),
			logger,
		)

		return server.Run()
	},
}

func init() {
	RootCmd.AddCommand(adminCmd)
}
