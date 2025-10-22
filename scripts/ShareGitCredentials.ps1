# PowerShell Script to enable and start the ssh-agent service on Windows
# Must be run as Administrator

Write-Host "=== Checking OpenSSH Client ==="
$sshClient = Get-WindowsCapability -Online | Where-Object Name -like 'OpenSSH.Client*'

if ($sshClient.State -ne "Installed") {
    Write-Host "OpenSSH Client is not installed. Installing..."
    Add-WindowsCapability -Online -Name OpenSSH.Client~~~~0.0.1.0
}
else {
    Write-Host "OpenSSH Client is already installed."
}

Write-Host "`n=== Checking ssh-agent service ==="
$sshService = Get-Service -Name ssh-agent -ErrorAction SilentlyContinue

if (-not $sshService) {
    Write-Host "ssh-agent service not found. Enabling OpenSSH Authentication Agent..."
    Enable-WindowsOptionalFeature -Online -FeatureName "OpenSSH-AuthenticationAgent" -NoRestart
}
else {
    Write-Host "ssh-agent service already exists."
}

Write-Host "`n=== Setting ssh-agent service to start automatically ==="
Set-Service ssh-agent -StartupType Automatic -ErrorAction SilentlyContinue

Write-Host "=== Starting ssh-agent service ==="
Start-Service ssh-agent -ErrorAction SilentlyContinue

Write-Host "`n=== Final ssh-agent status ==="
Get-Service ssh-agent | Select-Object Name, DisplayName, Status
