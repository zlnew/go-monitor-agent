#!/bin/bash
set -e

APP_NAME="horizonx-server"
MIGRATE_NAME="horizonx-server-migrate"
SEED_NAME="horizonx-server-seed"

SERVICE_NAME="horizonx-server"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/horizonx-server"
DATA_DIR="/var/lib/horizonx-server"
LOG_DIR="/var/log/horizonx-server"
USER="horizonx"
GROUP="horizonx"

BINARY_PATH="./bin/${APP_NAME}"
MIGRATE_BINARY_PATH="./bin/${MIGRATE_NAME}"
SEED_BINARY_PATH="./bin/${SEED_NAME}"

ENV_PATH="./.env.production"
PROD_DB_PATH="${DATA_DIR}/horizonx.db"

echo "=== HorizonX Deployment Script (Root Mode) ==="

# 1. Build binaries (Server + Tools)
echo "Building binaries..."
make build

# 2. Stop existing service
if systemctl list-units --full -all | grep -q "${SERVICE_NAME}.service"; then
    echo "Stopping existing systemd service..."
    sudo systemctl stop $SERVICE_NAME || true
fi

# 3. Create system user
if ! id -u $USER >/dev/null 2>&1; then
    echo "Creating system user $USER..."
    sudo useradd -r -s /bin/false $USER
fi

# 4. Setup directories
echo "Setting up directories..."
sudo mkdir -p $CONFIG_DIR $DATA_DIR $LOG_DIR
# Config & Log owned by horizonx user
sudo chown -R $USER:$GROUP $CONFIG_DIR $DATA_DIR $LOG_DIR

# 5. Deploy Binaries
echo "Deploying binaries to ${INSTALL_DIR}..."
# Copy Server
sudo cp $BINARY_PATH $INSTALL_DIR/$APP_NAME
sudo chown root:root $INSTALL_DIR/$APP_NAME
sudo chmod +x $INSTALL_DIR/$APP_NAME

# Copy Migration Tool
sudo cp $MIGRATE_BINARY_PATH $INSTALL_DIR/$MIGRATE_NAME
sudo chown root:root $INSTALL_DIR/$MIGRATE_NAME
sudo chmod +x $INSTALL_DIR/$MIGRATE_NAME

# Copy Seeder Tool
sudo cp $SEED_BINARY_PATH $INSTALL_DIR/$SEED_NAME
sudo chown root:root $INSTALL_DIR/$SEED_NAME
sudo chmod +x $INSTALL_DIR/$SEED_NAME

# 6. Copy .env
echo "Copying config..."
sudo cp $ENV_PATH $CONFIG_DIR/.env
sudo chown $USER:$GROUP $CONFIG_DIR/.env
sudo chmod 600 $CONFIG_DIR/.env

# 7. RUN MIGRATIONS
echo "Running Database Migrations..."
sudo $INSTALL_DIR/$MIGRATE_NAME -op=up -db=$PROD_DB_PATH

# 8. Create systemd service
echo "Updating systemd service..."
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
sudo tee $SERVICE_FILE >/dev/null <<EOF
[Unit]
Description=HorizonX Server (Root Mode)
After=network.target

[Service]
Type=simple
EnvironmentFile=${CONFIG_DIR}/.env
ExecStart=${INSTALL_DIR}/${APP_NAME}
WorkingDirectory=${DATA_DIR}
Restart=on-failure
RestartSec=5
# Run as root for full sensor/SSH access
User=root
Group=root
# Output logs
StandardOutput=file:${LOG_DIR}/out.log
StandardError=file:${LOG_DIR}/error.log

[Install]
WantedBy=multi-user.target
EOF

# 9. Start Service
echo "Reloading systemd and starting service..."
sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME
sudo systemctl start $SERVICE_NAME
sudo systemctl status $SERVICE_NAME --no-pager

echo ""
echo "=== Deployment Complete! ==="
echo "Binaries installed:"
echo "- Server:  ${INSTALL_DIR}/${APP_NAME}"
echo "- Migrate: ${INSTALL_DIR}/${MIGRATE_NAME}"
echo "- Seed:    ${INSTALL_DIR}/${SEED_NAME}"
echo ""
echo "To seed the database (Owner User), run:"
echo "sudo OWNER_PASSWORD='yourpass' ${INSTALL_DIR}/${SEED_NAME} -db=${PROD_DB_PATH}"
