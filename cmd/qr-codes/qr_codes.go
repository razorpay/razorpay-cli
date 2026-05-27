package qrcodes

import "github.com/spf13/cobra"

const basePath = "/v1/payments/qr_codes"

var Cmd = &cobra.Command{
	Use:   "qr-codes",
	Short: "Manage QR codes",
}
