package cmd

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/spf13/cobra"
)

var customersCmd = &cobra.Command{
	Use:   "customers",
	Short: "Manage customers",
}

var customersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all customers",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		q := url.Values{}
		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			q.Set("count", fmt.Sprintf("%d", count))
		}
		if skip, _ := cmd.Flags().GetInt("skip"); skip > 0 {
			q.Set("skip", fmt.Sprintf("%d", skip))
		}
		data, err := client.Get("/customers", q)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var customersFetchCmd = &cobra.Command{
	Use:   "fetch <customer_id>",
	Short: "Fetch a customer by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		data, err := client.Get("/customers/"+args[0], nil)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var customersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new customer",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		contact, _ := cmd.Flags().GetString("contact")
		params, _ := cmd.Flags().GetStringArray("param")

		body := map[string]interface{}{}
		if name != "" {
			body["name"] = name
		}
		if email != "" {
			body["email"] = email
		}
		if contact != "" {
			body["contact"] = contact
		}
		extra, err := api.ParseParams(params)
		if err != nil {
			return err
		}
		for k, v := range extra {
			body[k] = v
		}

		if len(body) == 0 {
			return fmt.Errorf("at least one of --name, --email, or --contact is required")
		}

		data, err := client.Post("/customers", body)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var customersUpdateCmd = &cobra.Command{
	Use:   "update <customer_id>",
	Short: "Update a customer",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		contact, _ := cmd.Flags().GetString("contact")
		params, _ := cmd.Flags().GetStringArray("param")

		body := map[string]interface{}{}
		if name != "" {
			body["name"] = name
		}
		if email != "" {
			body["email"] = email
		}
		if contact != "" {
			body["contact"] = contact
		}
		extra, err := api.ParseParams(params)
		if err != nil {
			return err
		}
		for k, v := range extra {
			body[k] = v
		}

		data, err := client.Patch("/customers/"+args[0], body)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	customersCmd.AddCommand(customersListCmd)
	customersCmd.AddCommand(customersFetchCmd)
	customersCmd.AddCommand(customersCreateCmd)
	customersCmd.AddCommand(customersUpdateCmd)

	customersListCmd.Flags().Int("count", 10, "Number of customers to fetch")
	customersListCmd.Flags().Int("skip", 0, "Number of customers to skip")

	customersCreateCmd.Flags().String("name", "", "Customer name")
	customersCreateCmd.Flags().String("email", "", "Customer email address")
	customersCreateCmd.Flags().String("contact", "", "Customer contact number")
	customersCreateCmd.Flags().StringArray("param", nil, "Additional parameter as key=value")

	customersUpdateCmd.Flags().String("name", "", "Customer name")
	customersUpdateCmd.Flags().String("email", "", "Customer email address")
	customersUpdateCmd.Flags().String("contact", "", "Customer contact number")
	customersUpdateCmd.Flags().StringArray("param", nil, "Additional parameter as key=value")
}
