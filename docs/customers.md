# Customers

Customer records store contact information for reuse across orders and payments. Associating a customer with a payment enables features like saved cards and payment history.

## Create a customer

```bash
razorpay customers create [--name <name>] [--email <email>] [--contact <phone>] [--param <key=value>]
```

At least one of `--name`, `--email`, or `--contact` is required.

Flags:

| Flag        | Type   | Description                              |
| ----------- | ------ | ---------------------------------------- |
| `--name`    | string | Customer's full name                     |
| `--email`   | string | Customer's email address                 |
| `--contact` | string | Customer's phone number (with country code) |
| `--param`   | string | Additional parameter as `key=value`, repeatable |

Examples:

```bash
# Create a customer with all fields
razorpay customers create \
  --name "Gaurav Kumar" \
  --email gaurav.kumar@example.com \
  --contact "+919000090000"

# Create a customer with just an email
razorpay customers create --email gaurav.kumar@example.com
```

## List customers

```bash
razorpay customers list [flags]
```

Flags:

| Flag      | Type | Default | Description                    |
| --------- | ---- | ------- | ------------------------------ |
| `--count` | int  | 10      | Number of customers to fetch   |
| `--skip`  | int  | 0       | Records to skip for pagination |

Example:

```bash
razorpay customers list --count 25
```

## Fetch a customer

```bash
razorpay customers fetch <customer_id>
```

Example:

```bash
razorpay customers fetch cust_K6fNE0WJZWGqtN
```

## Update a customer

```bash
razorpay customers update <customer_id> [--name <name>] [--email <email>] [--contact <phone>] [--param <key=value>]
```

Example:

```bash
# Update a customer's email address
razorpay customers update cust_K6fNE0WJZWGqtN --email newemail@example.com
```
