# Check if running as administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "This script requires administrator privileges." -ForegroundColor Red
    Write-Host "Please run PowerShell as administrator and try again." -ForegroundColor Red
    Write-Host ""
    Write-Host "Press any key to exit..."
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit
}

$hostsFile = "$env:windir\System32\drivers\etc\hosts"

# Check if entry already exists
$entryExists = Get-Content -Path $hostsFile | Where-Object { $_ -match "127\.0\.0\.1\s+minio" }

if ($entryExists) {
    Write-Host "Host entry already exists in $hostsFile" -ForegroundColor Green
} else {
    try {
        # Add entry to hosts file
        Add-Content -Path $hostsFile -Value "127.0.0.1 minio" -Force
        Write-Host "Entry added successfully to $hostsFile" -ForegroundColor Green
    } catch {
        Write-Host "Failed to write to hosts file: $_" -ForegroundColor Red
    }
}

Write-Host "`nDone! To test, try: ping minio"
Write-Host "`nPress any key to exit..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown") 