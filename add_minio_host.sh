#!/bin/bash

# Function to check if the entry already exists
check_entry_exists() {
    if grep -q "127.0.0.1[[:space:]]minio" "$1"; then
        return 0 # Entry exists
    else
        return 1 # Entry doesn't exist
    fi
}

# Function to add the entry to hosts file
add_entry() {
    echo "127.0.0.1 minio" >> "$1"
    echo "Entry added successfully to $1"
}

# Detect OS
if [[ "$OSTYPE" == "linux-gnu"* ]] || [[ "$OSTYPE" == "darwin"* ]]; then
    # Linux or macOS
    HOSTS_FILE="/etc/hosts"
    
    if check_entry_exists "$HOSTS_FILE"; then
        echo "Host entry already exists in $HOSTS_FILE"
    else
        echo "Adding entry to $HOSTS_FILE (requires sudo permissions)"
        sudo bash -c "echo '127.0.0.1 minio' >> $HOSTS_FILE"
        echo "Entry added successfully"
    fi
    
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ -n "$WINDIR" ]]; then
    # Windows
    HOSTS_FILE="$WINDIR/System32/drivers/etc/hosts"
    
    if [[ ! -n "$WINDIR" ]]; then
        HOSTS_FILE="C:/Windows/System32/drivers/etc/hosts"
    fi
    
    if check_entry_exists "$HOSTS_FILE"; then
        echo "Host entry already exists in $HOSTS_FILE"
    else
        echo "Adding entry to $HOSTS_FILE (requires administrator privileges)"
        echo "Please run the following command as administrator in PowerShell:"
        echo "Add-Content -Path \"$HOSTS_FILE\" -Value \"127.0.0.1 minio\" -Force"
        
        # Try to add automatically if running with admin privileges
        if command -v powershell.exe &> /dev/null; then
            powershell.exe -Command "Start-Process powershell -Verb RunAs -ArgumentList '-Command Add-Content -Path \"$HOSTS_FILE\" -Value \"127.0.0.1 minio\" -Force'"
        fi
    fi
else
    echo "Unsupported operating system"
    exit 1
fi

echo "Done! To test, try to ping minio: ping minio" 