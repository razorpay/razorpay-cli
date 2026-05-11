# Installing the Razorpay CLI

---

## macOS

1. Download the binary for your CPU architecture:

   | Architecture | Download |
   |---|---|
   | Apple Silicon (M1/M2/M3) | [razorpay_mac-os_arm64.tar.gz](https://<YOUR_S3_BUCKET>.s3.<YOUR_AWS_REGION>.amazonaws.com/razorpay-cli/latest/razorpay_mac-os_arm64.tar.gz) |
   | Intel | [razorpay_mac-os_x86_64.tar.gz](https://<YOUR_S3_BUCKET>.s3.<YOUR_AWS_REGION>.amazonaws.com/razorpay-cli/latest/razorpay_mac-os_x86_64.tar.gz) |

2. Extract the archive:
   ```bash
   tar -xvf razorpay_mac-os_<arch>.tar.gz
   ```

3. Move the binary to your execution path:
   ```bash
   sudo mv razorpay /usr/local/bin/
   ```

4. Verify:
   ```bash
   razorpay --version
   ```

---

## Linux

1. Download the binary for your architecture:

   | Architecture | Download |
   |---|---|
   | x86-64 (.tar.gz) | [razorpay_linux_x86_64.tar.gz](https://<YOUR_S3_BUCKET>.s3.<YOUR_AWS_REGION>.amazonaws.com/razorpay-cli/latest/razorpay_linux_x86_64.tar.gz) |
   | ARM64 (.tar.gz) | [razorpay_linux_arm64.tar.gz](https://<YOUR_S3_BUCKET>.s3.<YOUR_AWS_REGION>.amazonaws.com/razorpay-cli/latest/razorpay_linux_arm64.tar.gz) |
   | x86-64 (.deb) | [razorpay_linux_amd64.deb](https://<YOUR_S3_BUCKET>.s3.<YOUR_AWS_REGION>.amazonaws.com/razorpay-cli/latest/razorpay_linux_amd64.deb) |
   | ARM64 (.deb) | [razorpay_linux_arm64.deb](https://<YOUR_S3_BUCKET>.s3.<YOUR_AWS_REGION>.amazonaws.com/razorpay-cli/latest/razorpay_linux_arm64.deb) |
   | x86-64 (.rpm) | [razorpay_linux_amd64.rpm](https://<YOUR_S3_BUCKET>.s3.<YOUR_AWS_REGION>.amazonaws.com/razorpay-cli/latest/razorpay_linux_amd64.rpm) |
   | ARM64 (.rpm) | [razorpay_linux_arm64.rpm](https://<YOUR_S3_BUCKET>.s3.<YOUR_AWS_REGION>.amazonaws.com/razorpay-cli/latest/razorpay_linux_arm64.rpm) |

2. Install based on the format downloaded:

   **tar.gz:**
   ```bash
   tar -xvf razorpay_linux_<arch>.tar.gz
   sudo mv razorpay /usr/local/bin/
   ```

   **deb (Debian / Ubuntu):**
   ```bash
   sudo dpkg -i razorpay_linux_<arch>.deb
   ```

   **rpm (RHEL / Fedora / Amazon Linux):**
   ```bash
   sudo rpm -i razorpay_linux_<arch>.rpm
   ```

3. Verify:
   ```bash
   razorpay --version
   ```

---

## Windows

1. Download the binary for your architecture:

   | Architecture | Download |
   |---|---|
   | x86-64 | [razorpay_windows_x86_64.zip](https://<YOUR_S3_BUCKET>.s3.<YOUR_AWS_REGION>.amazonaws.com/razorpay-cli/latest/razorpay_windows_x86_64.zip) |
   | x86 (32-bit) | [razorpay_windows_i386.zip](https://<YOUR_S3_BUCKET>.s3.<YOUR_AWS_REGION>.amazonaws.com/razorpay-cli/latest/razorpay_windows_i386.zip) |

2. Unzip the file:
   ```powershell
   Expand-Archive razorpay_windows_<arch>.zip -DestinationPath .
   ```

3. Add the path to `razorpay.exe` to your `Path` environment variable:
   - Open **System Properties → Environment Variables**
   - Under **System variables**, select `Path` → **Edit** → **New**
   - Add the folder path where `razorpay.exe` is located

4. Open a new terminal and verify:
   ```powershell
   razorpay --version
   ```

---

## Next steps

Configure your API credentials:

```bash
razorpay configure
```

You will be prompted for your **Key ID** and **Key Secret** from the [Razorpay Dashboard](https://dashboard.razorpay.com/app/website-app-settings/api-keys).

For CI/CD environments, use environment variables instead:

```bash
export RAZORPAY_KEY_ID=rzp_test_xxxxxxxxxxxx
export RAZORPAY_KEY_SECRET=xxxxxxxxxxxxxxxxxxxx
```
