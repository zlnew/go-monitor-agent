#!/bin/bash
# Enable immediate exit on command failure
set -e

# --- Configuration (No changes here, keeping it clean) ---
APP_NAME="horizonx-server"
AGENT_NAME="horizonx-agent"
MIGRATE_TOOL="horizonx-migrate"
SEED_TOOL="horizonx-seed"

SERVER_SERVICE="horizonx-server"
AGENT_SERVICE="horizonx-agent"

INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/horizonx"
LOG_DIR="/var/log/horizonx"

# User System
SYS_USER="horizonx"
SYS_GROUP="horizonx"

# Source Paths
BIN_SRC="./bin"
ENV_SERVER_SRC="./.env.server.prod"
ENV_AGENT_SRC="./.env.agent.prod"

# Deployment Flags
DEPLOY_SERVER=false
DEPLOY_AGENT=false

echo "=== HorizonX Secure Deployment (Split Config) ==="

# --- Helper Functions ---

# Function to setup the base environment (user, dirs, binaries)
setup_base() {
    echo "🛠️  Building project..."
    make build

    echo "🛑 Stopping existing services..."
    # Using || true to ignore stop errors if service is not running
    sudo systemctl stop $SERVER_SERVICE || true
    sudo systemctl stop $AGENT_SERVICE || true

    # 3. Create System User
    if ! id -u $SYS_USER >/dev/null 2>&1; then
        echo "👤 Creating system user $SYS_USER..."
        sudo useradd -r -s /bin/false $SYS_USER
    fi

    # 4. Directories
    echo "📂 Setting up directories..."
    sudo mkdir -p $CONFIG_DIR $LOG_DIR
    sudo chown -R $SYS_USER:$SYS_GROUP $LOG_DIR
    sudo chmod 755 $LOG_DIR
    sudo chmod 755 $CONFIG_DIR # Config dir owned by root but readable

    # 5. Deploy Binaries (Deploy ALL, even if not used, for simplicity)
    echo "🚀 Copying binaries..."
    sudo cp $BIN_SRC/server $INSTALL_DIR/$APP_NAME
    sudo cp $BIN_SRC/agent $INSTALL_DIR/$AGENT_NAME
    sudo cp $BIN_SRC/migrate $INSTALL_DIR/$MIGRATE_TOOL
    sudo cp $BIN_SRC/seed $INSTALL_DIR/$SEED_TOOL
    sudo chmod +x $INSTALL_DIR/*
}

# Function for Server specific deployment steps
deploy_server() {
    echo "--- Server Setup ---"

    # 6A. Server Config
    echo "📄 Deploying Server Configuration..."
    if [ -f "$ENV_SERVER_SRC" ]; then
        sudo cp $ENV_SERVER_SRC $CONFIG_DIR/server.env
        # server.env secure: readable by root, writable by root, readable by system group
        sudo chown root:$SYS_GROUP $CONFIG_DIR/server.env
        sudo chmod 640 $CONFIG_DIR/server.env
        echo "    -> server.env deployed (Secure)"
    else
        echo "⚠️  FATAL: $ENV_SERVER_SRC not found! Cannot deploy Server."
        exit 1
    fi

    # 7. Run Migrations (Requires server.env)
    echo "📦 Running Database Migrations..."
    # Load env vars temporarily for the migration tool
    sudo sh -c "set -a; source $CONFIG_DIR/server.env; set +a; $INSTALL_DIR/$MIGRATE_TOOL -op=up"
    
    # 8. Setup Systemd: SERVER
    echo "⚙️  Configuring Server Service..."
    SERVER_UNIT="/etc/systemd/system/${SERVER_SERVICE}.service"
    sudo tee $SERVER_UNIT >/dev/null <<EOF
[Unit]
Description=HorizonX Core Server
After=network.target postgresql.service

[Service]
Type=simple
EnvironmentFile=$CONFIG_DIR/server.env
ExecStart=$INSTALL_DIR/$APP_NAME
Restart=always
RestartSec=5
User=$SYS_USER
Group=$SYS_GROUP
StandardOutput=append:${LOG_DIR}/server.log
StandardError=append:${LOG_DIR}/server.error.log

[Install]
WantedBy=multi-user.target
EOF

    # Start/Enable Service
    echo "🔥 Enabling and Starting Server Service..."
    sudo systemctl daemon-reload
    sudo systemctl enable $SERVER_SERVICE
    sudo systemctl start $SERVER_SERVICE
    echo "✅ Server Deployment Complete!"

    # Hint for seeding (since this is server)
    echo ""
    echo "---------------------------------------------------------"
    echo "🌱  HINT: Database Server Deployed."
    echo "    If you need to seed the database, run this command manually:"
    echo ""
    echo "    sudo sh -c \"set -a; source $CONFIG_DIR/server.env; set +a; $INSTALL_DIR/$SEED_TOOL\""
    echo "---------------------------------------------------------"
}

# Function for Agent specific deployment steps
deploy_agent() {
    echo "--- Agent Setup ---"

    # 6B. Agent Config
    echo "📄 Deploying Agent Configuration..."
    if [ -f "$ENV_AGENT_SRC" ]; then
        sudo cp $ENV_AGENT_SRC $CONFIG_DIR/agent.env
        # agent.env very secure: readable/writable by root only
        sudo chown root:root $CONFIG_DIR/agent.env
        sudo chmod 600 $CONFIG_DIR/agent.env
        echo "    -> agent.env deployed (Secure)"
    else
        echo "⚠️  WARNING: $ENV_AGENT_SRC not found! Cannot deploy Agent."
        return 1 # Non-fatal error, just skip agent start/enable
    fi

    # 9. Setup Systemd: AGENT
    echo "⚙️  Configuring Agent Service..."
    AGENT_UNIT="/etc/systemd/system/${AGENT_SERVICE}.service"
    sudo tee $AGENT_UNIT >/dev/null <<EOF
[Unit]
Description=HorizonX Metrics Agent
After=network.target

[Service]
Type=simple
EnvironmentFile=$CONFIG_DIR/agent.env
ExecStart=$INSTALL_DIR/$AGENT_NAME
Restart=always
RestartSec=5
# Run as root for hardware access (as per original script)
User=root
Group=root
StandardOutput=append:${LOG_DIR}/agent.log
StandardError=append:${LOG_DIR}/agent.error.log

[Install]
WantedBy=multi-user.target
EOF

    # Start/Enable Service
    echo "🔥 Enabling and Starting Agent Service..."
    # Systemd reload is needed if deploy_server wasn't called
    if ! $DEPLOY_SERVER; then
        sudo systemctl daemon-reload
    fi
    sudo systemctl enable $AGENT_SERVICE
    sudo systemctl start $AGENT_SERVICE
    echo "✅ Agent Deployment Complete!"
}

# --- Main Logic ---

# 1. Selection Menu
echo ""
echo "Choose deployment type:"
echo "1) Server Only (Core App, DB Migrations)"
echo "2) Agent Only (Monitoring Client)"
echo "3) All (Server + Agent)"
echo "4) Exit"
echo -n "Enter choice [1-4]: "
read choice

case "$choice" in
    1)
        DEPLOY_SERVER=true
        ;;
    2)
        DEPLOY_AGENT=true
        ;;
    3)
        DEPLOY_SERVER=true
        DEPLOY_AGENT=true
        ;;
    4)
        echo "Exiting deployment."
        exit 0
        ;;
    *)
        echo "❌ Invalid choice. Exiting."
        exit 1
        ;;
esac

echo ""
echo "Starting setup..."
setup_base # Run base setup regardless of choice

if $DEPLOY_SERVER; then
    deploy_server
fi

if $DEPLOY_AGENT; then
    deploy_agent
fi

echo ""
echo "========================================================="
echo "✅ Deployment Process Finished."
echo "Check systemctl status $SERVER_SERVICE and $AGENT_SERVICE for details."
echo "========================================================="
