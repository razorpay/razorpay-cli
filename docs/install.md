# Install the Razorpay CLI

The Razorpay CLI lets you interact with Razorpay APIs directly from your terminal. Use it to test integrations, trigger events, manage resources and automate workflows in your CI/CD pipelines.

To install the Razorpay CLI, you must complete the following actions:

1. [Download and Install the Binary](#step-1-download-and-install-the-binary)
2. [Configure Your API Credentials](#step-2-configure-your-api-credentials)

## Step 1: Download and Install the Binary

### Quick Install (macOS and Linux)

Run the install script to automatically detect your OS and architecture, download the latest binary, and install it to `/usr/local/bin`:

```bash
curl -fsSL https://razorpay.com/cli/latest/install.sh | bash
```

### macOS

1. Download the binary that matches your CPU architecture.

   | Architecture | Download |
   |---|---|
   | Apple Silicon | [razorpay_mac-os_arm64.tar.gz](https://razorpay.com/cli/latest/razorpay_mac-os_arm64.tar.gz) |
   | Intel | [razorpay_mac-os_x86_64.tar.gz](https://razorpay.com/cli/latest/razorpay_mac-os_x86_64.tar.gz) |

   Or download using curl:

   **Apple Silicon:**
   ```bash
   curl -fsSL https://razorpay.com/cli/latest/razorpay_mac-os_arm64.tar.gz | tar -xz
   ```

   **Intel:**
   ```bash
   curl -fsSL https://razorpay.com/cli/latest/razorpay_mac-os_x86_64.tar.gz | tar -xz
   ```

2. Extract the archive.

   ```bash
   tar -xvf razorpay_mac-os_<arch>.tar.gz
   ```

3. Make the binary executable and remove the macOS quarantine attribute.

   ```bash
   chmod +x ./razorpay
   xattr -d com.apple.quarantine ./razorpay
   ```

4. Move the binary to your execution path.

   ```bash
   sudo mv razorpay /usr/local/bin/
   ```

5. Verify the installation.

   ```bash
   razorpay --version
   ```

---

### Linux

1. Download the binary that matches your architecture and preferred package format.

   | Architecture | Format | Download |
   |---|---|---|
   | x86-64 | tar.gz | [razorpay_linux_x86_64.tar.gz](https://razorpay.com/cli/latest/razorpay_linux_x86_64.tar.gz) |
   | ARM64 | tar.gz | [razorpay_linux_arm64.tar.gz](https://razorpay.com/cli/latest/razorpay_linux_arm64.tar.gz) |
   | x86-64 | deb | [razorpay_linux_amd64.deb](https://razorpay.com/cli/latest/razorpay_linux_amd64.deb) |
   | ARM64 | deb | [razorpay_linux_arm64.deb](https://razorpay.com/cli/latest/razorpay_linux_arm64.deb) |
   | x86-64 | rpm | [razorpay_linux_amd64.rpm](https://razorpay.com/cli/latest/razorpay_linux_amd64.rpm) |
   | ARM64 | rpm | [razorpay_linux_arm64.rpm](https://razorpay.com/cli/latest/razorpay_linux_arm64.rpm) |

   Or download using curl:

   **x86-64 (tar.gz):**
   ```bash
   curl -fsSL https://razorpay.com/cli/latest/razorpay_linux_x86_64.tar.gz | tar -xz
   ```

   **ARM64 (tar.gz):**
   ```bash
   curl -fsSL https://razorpay.com/cli/latest/razorpay_linux_arm64.tar.gz | tar -xz
   ```

   **x86-64 (deb):**
   ```bash
   curl -fsSL https://razorpay.com/cli/latest/razorpay_linux_amd64.deb -o /tmp/razorpay.deb \
   && sudo dpkg -i /tmp/razorpay.deb
   ```

   **ARM64 (deb):**
   ```bash
   curl -fsSL https://razorpay.com/cli/latest/razorpay_linux_arm64.deb -o /tmp/razorpay.deb \
   && sudo dpkg -i /tmp/razorpay.deb
   ```

   **x86-64 (rpm):**
   ```bash
   curl -fsSL https://razorpay.com/cli/latest/razorpay_linux_amd64.rpm -o /tmp/razorpay.rpm \
   && sudo rpm -i /tmp/razorpay.rpm
   ```

   **ARM64 (rpm):**
   ```bash
   curl -fsSL https://razorpay.com/cli/latest/razorpay_linux_arm64.rpm -o /tmp/razorpay.rpm \
   && sudo rpm -i /tmp/razorpay.rpm
   ```

2. Install based on the format you downloaded.

   For `tar.gz`:

   ```bash
   tar -xvf razorpay_linux_<arch>.tar.gz
   sudo mv razorpay /usr/local/bin/
   ```

   For `deb` (Debian and Ubuntu):

   ```bash
   sudo dpkg -i razorpay_linux_<arch>.deb
   ```

   For `rpm` (RHEL, Fedora and Amazon Linux):

   ```bash
   sudo rpm -i razorpay_linux_<arch>.rpm
   ```

3. Verify the installation.

   ```bash
   razorpay --version
   ```

---

### Windows

1. Download the binary that matches your architecture.

   | Architecture | Download |
   |---|---|
   | x86-64 | [razorpay_windows_x86_64.zip](https://razorpay.com/cli/latest/razorpay_windows_x86_64.zip) |
   | x86 (32-bit) | [razorpay_windows_i386.zip](https://razorpay.com/cli/latest/razorpay_windows_i386.zip) |

   Or download using curl in PowerShell:

   **x86-64:**
   ```powershell
   curl.exe -fsSL https://razorpay.com/cli/latest/razorpay_windows_x86_64.zip
   ```

   **x86 (32-bit):**
   ```powershell
   curl.exe -fsSL https://razorpay.com/cli/latest/razorpay_windows_i386.zip
   ```

2. Extract the archive.

   Open PowerShell in the folder where you downloaded the file and run:

   ```powershell
   Expand-Archive razorpay_windows_<arch>.zip -DestinationPath .
   ```

3. Add the binary to your Path environment variable.

   1. Open **System Properties** > **Environment Variables**.
   2. Under **System variables**, select **Path** > **Edit** > **New**.
   3. Add the folder path where `razorpay.exe` is located.
   4. Click **OK** to save the changes.

4. Verify the installation.

   Open a new terminal and run:

   ```bash
   razorpay --version
   ```

---

## Step 2: Configure Your API Credentials

After installing the CLI, configure it with your API credentials to start making requests. Generate your keys from the [Razorpay Dashboard](https://dashboard.razorpay.com/app/keys).

You can configure credentials in two ways:

**Option 1: Pass credentials directly in the command**

```bash
razorpay configure --key-id rzp_test_xxxxxxxxxxxx --key-secret xxxxxxxxxxxxxxxxxxxx
```

**Option 2: Use the interactive prompt**

Run the command without arguments and enter your credentials when prompted:

```bash
razorpay configure
```

```
Enter your Razorpay Key ID: rzp_test_xxxxxxxxxxxx
Enter your Razorpay Key Secret: xxxxxxxxxxxxxxxxxxxx
```

> **Test Mode and Live Mode**
>
> The CLI supports both Test Mode and Live Mode keys. Use Test Mode keys (prefixed with `rzp_test_`) while building and testing your integration.

> **Verify Your Setup**
>
> Once configured, run `razorpay --version` to confirm the CLI is installed and ready to use.

## Related Information

- [Authentication](https://razorpay.com/docs/api/authentication/)
- [Sandbox Setup](https://razorpay.com/docs/api/sandbox-setup/)
- [API Reference Guide](https://razorpay.com/docs/api/)
