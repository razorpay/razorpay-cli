# install.ps1 - Installer for the Razorpay CLI (Windows)
# Usage: powershell -ExecutionPolicy ByPass -c "irm https://raw.githubusercontent.com/razorpay/razorpay-cli/master/install.ps1 | iex"
[CmdletBinding()]
param (
    # Override the installation directory (defaults to %USERPROFILE%\.local\bin)
    [string]$InstallDir = ""
)
$ErrorActionPreference = "Stop"

$Repo      = "razorpay/razorpay-cli"
$ApiBase   = "https://api.github.com/repos/$Repo"
$BinaryName = "razorpay.exe"

function Write-Status([string]$Message) {
    Write-Host "razorpay-cli: $Message"
}

function Write-Err([string]$Message) {
    Write-Error "razorpay-cli: error: $Message"
    exit 1
}

# ---- architecture detection -------------------------------------------------

function Get-Arch {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64"  { return "x86_64" }
        "ARM64"  { return "arm64" }
        "x86"    { return "i386" }
        default  { Write-Err "unsupported architecture: $arch" }
    }
}

# ---- install directory selection -------------------------------------------

function Get-InstallDir {
    if ($InstallDir -ne "") {
        return $InstallDir
    }
    if ($env:RAZORPAY_INSTALL -ne $null -and $env:RAZORPAY_INSTALL -ne "") {
        return $env:RAZORPAY_INSTALL
    }
    return "$env:USERPROFILE\.local\bin"
}

# ---- checksum verification --------------------------------------------------

function Test-Checksum([string]$ArchivePath, [string]$ChecksumsPath) {
    $archiveName = Split-Path $ArchivePath -Leaf
    $expected = (Get-Content $ChecksumsPath | Where-Object { $_ -match $archiveName }) -replace '\s+.*', ''
    if (-not $expected) {
        Write-Status "warning: could not find checksum entry for $archiveName, skipping verification"
        return
    }
    $actual = (Get-FileHash -Algorithm SHA256 $ArchivePath).Hash.ToLower()
    if ($expected -ne $actual) {
        Write-Err "checksum mismatch for ${archiveName}`n  expected: $expected`n  actual:   $actual"
    }
    Write-Status "checksum verified"
}

# ---- PATH update ------------------------------------------------------------

function Add-ToUserPath([string]$Dir) {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -split ";" -contains $Dir) {
        return  # Already present
    }
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$currentPath;$Dir",
        "User"
    )
    Write-Status "added $Dir to your user PATH (restart your shell to take effect)"
}

# ---- main -------------------------------------------------------------------

function Main {
    $arch       = Get-Arch
    $installDir = Get-InstallDir

    Write-Status "detecting latest release..."
    $releaseInfo = Invoke-RestMethod -Uri "$ApiBase/releases/latest" -Headers @{ "User-Agent" = "razorpay-cli-installer" }
    $latestTag   = $releaseInfo.tag_name
    if (-not $latestTag) {
        Write-Err "could not determine latest release tag"
    }

    # Strip leading 'v' for filenames
    $version = $latestTag.TrimStart("v")

    $archiveName    = "razorpay-cli_Windows_${arch}.zip"
    $checksumsName  = "razorpay-cli_${version}_checksums.txt"
    $baseUrl        = "https://github.com/$Repo/releases/download/$latestTag"

    Write-Status "installing razorpay-cli $latestTag (Windows/$arch)"

    $tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.IO.Path]::GetRandomFileName())
    New-Item -ItemType Directory -Path $tmpDir | Out-Null

    try {
        $archivePath   = Join-Path $tmpDir $archiveName
        $checksumsPath = Join-Path $tmpDir $checksumsName

        Write-Status "downloading $archiveName..."
        Invoke-WebRequest -Uri "$baseUrl/$archiveName"   -OutFile $archivePath   -UseBasicParsing
        Invoke-WebRequest -Uri "$baseUrl/$checksumsName" -OutFile $checksumsPath -UseBasicParsing

        Test-Checksum $archivePath $checksumsPath

        Write-Status "extracting..."
        Expand-Archive -Path $archivePath -DestinationPath $tmpDir -Force

        # Locate the binary — GoReleaser may or may not wrap in a subdirectory
        $binaryPath = Get-ChildItem -Path $tmpDir -Recurse -Filter $BinaryName | Select-Object -First 1 -ExpandProperty FullName
        if (-not $binaryPath) {
            Write-Err "could not find '$BinaryName' binary in the downloaded archive"
        }

        if (-not (Test-Path $installDir)) {
            New-Item -ItemType Directory -Path $installDir -Force | Out-Null
        }

        $destBinary = Join-Path $installDir $BinaryName
        Write-Status "installing to $destBinary"
        Copy-Item -Path $binaryPath -Destination $destBinary -Force

        # Verify
        $testOutput = & $destBinary --help 2>&1
        if ($LASTEXITCODE -ne 0 -and $LASTEXITCODE -ne 1) {
            Write-Err "installed binary failed to run — please report this at https://github.com/$Repo/issues"
        }

        Add-ToUserPath $installDir

        Write-Status ""
        Write-Status "razorpay-cli $latestTag installed successfully!"
        Write-Status "run 'razorpay --help' to get started (you may need to restart your terminal)"
    }
    finally {
        Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue
    }
}

Main
