#!/bin/sh

os=$(uname -s)
arch=$(uname -m)

link="https://github.com/jonasroussel/hyve/releases/latest/download/hyve_${os}_${arch}.tar.gz"

# Download the binary with wget or curl
if command -v wget &> /dev/null
then
	wget -O hyve.tar.gz "$link"
elif command -v curl &> /dev/null
then
	curl -o hyve.tar.gz "$link"
else
	echo "wget or curl is required"
	exit 1
fi

# Extract the binary
tar -xzf hyve.tar.gz

# Remove the archive
rm -f hyve.tar.gz

# Make the binary executable
chmod +x hyve

# Move the binary to /usr/bin
sudo mv hyve /usr/bin

# Ask for the target url
read -p "Please enter the target URL to which requests will be proxied: " target

# Create the config.env file
sudo mkdir -p /etc/hyve
sudo cat <<EOF | sudo tee /etc/hyve/config.env
TARGET=$target
DATA_DIR=/var/lib/hyve
USER_DIR=/var/lib/hyve/user
STORE=file
STORE_DIR=/var/lib/hyve/certificates
EOF

# Create the data directory
sudo mkdir -p /var/lib/hyve

# Create the systemd service file
sudo cat <<EOF | sudo tee /usr/lib/systemd/system/hyve.service
[Unit]
Description=Hyve
Documentation=https://github.com/jonasroussel/hyve#readme
After=network.target network-online.target
Requires=network-online.target

[Service]
Type=exec
User=root
Group=root
ExecStart=/usr/bin/hyve
EnvironmentFile=/etc/hyve/config.env
TimeoutStopSec=5s
LimitNOFILE=1048576
LimitNPROC=512
ProtectSystem=full
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
EOF

# Reload the systemd daemon
sudo systemctl daemon-reload

# Enable and start the service
sudo systemctl enable hyve
sudo systemctl start hyve

# Print the admin API URL
echo "Installation complete!!"