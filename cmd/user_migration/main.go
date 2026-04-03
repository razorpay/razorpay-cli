package main

import (
	"github.com/razorpay/foundation/database/migration"
)

func main() {
	migration.Run(
		"./config/user",
		"./migrations",
		"primary_db",
	)
}
