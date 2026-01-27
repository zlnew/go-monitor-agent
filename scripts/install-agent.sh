#!/bin/bash
set -e
# =============================
# Config
# =============================
USER_NAME="horizonx"
GROUP_NAME="horizonx"
DATA_DIR="/var/lib/horizonx"
SSH_DIR="$DATA_DIR/.ssh"
HOME_DIR="$DATA_DIR"
INSTALL_BIN="/usr/local/bin/horizonx-agent"
BIN_SOURCE="./bin/agent"
CONFIG_DIR="$DATA_DIR/agent.env"
LOG_DIR="/var/log/horizonx"
SERVICE_NAME="horizonx-agent"
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"

# =============================
# User & Group
# =============================
if ! id -u "$USER_NAME" >/dev/null 2>&1; then
  echo "[*] Creating user $USER_NAME"
  groupadd -f "$GROUP_NAME"
  useradd -r \
    -g "$GROUP_NAME" \
    -d "$HOME_DIR" \
    -s /usr/sbin/nologin \
    "$USER_NAME"
else
  echo "[*] User exists, skipping"
fi

# =============================
# Directories
# =============================
echo "[*] Creating directories"
mkdir -p "$SSH_DIR" "$LOG_DIR"
touch "$LOG_DIR/agent.log" "$LOG_DIR/agent.error.log"
chown -R "$USER_NAME:$GROUP_NAME" "$DATA_DIR" "$LOG_DIR"
chmod 700 "$SSH_DIR"

# =============================
# SSH key
# =============================
SSH_KEY="$SSH_DIR/id_ed25519"
if [ ! -f "$SSH_KEY" ]; then
  echo "[*] Generating SSH key"
  sudo -u "$USER_NAME" env HOME="$HOME_DIR" \
    ssh-keygen -t ed25519 \
    -f "$SSH_KEY" \
    -N "" \
    -C "horizonx-agent@$(hostname)"
else
  echo "[*] SSH key exists"
fi
chmod 600 "$SSH_KEY"
chmod 644 "$SSH_KEY.pub"
chown "$USER_NAME:$GROUP_NAME" "$SSH_KEY" "$SSH_KEY.pub"

# =============================
# SSH config
# =============================
SSH_CONFIG="$SSH_DIR/config"
echo "[*] Writing SSH config"
cat > "$SSH_CONFIG" <<EOF
Host *
  IdentityFile $SSH_KEY
  UserKnownHostsFile $SSH_DIR/known_hosts
  StrictHostKeyChecking yes
  IdentitiesOnly yes
EOF
chmod 600 "$SSH_CONFIG"
chown "$USER_NAME:$GROUP_NAME" "$SSH_CONFIG"

# =============================
# Auto-add Git Provider Known Hosts
# =============================
KNOWN_HOSTS_FILE="$SSH_DIR/known_hosts"
touch "$KNOWN_HOSTS_FILE"
chmod 644 "$KNOWN_HOSTS_FILE"
chown "$USER_NAME:$GROUP_NAME" "$KNOWN_HOSTS_FILE"

echo "[*] Adding common Git providers to known_hosts"
GIT_PROVIDERS=(
  "github.com"
  "gitlab.com"
  "bitbucket.org"
  "ssh.dev.azure.com"
  "vs-ssh.visualstudio.com"
)

for provider in "${GIT_PROVIDERS[@]}"; do
  if ! grep -q "$provider" "$KNOWN_HOSTS_FILE" 2>/dev/null; then
    echo "  [+] Scanning $provider..."
    ssh-keyscan -H -t rsa,ed25519 "$provider" >> "$KNOWN_HOSTS_FILE" 2>/dev/null || \
      echo "  [!] Failed to scan $provider (skipping)"
  else
    echo "  [√] $provider already in known_hosts"
  fi
done

# Verify known_hosts populated
if [ ! -s "$KNOWN_HOSTS_FILE" ]; then
  echo "[!] WARNING: known_hosts is empty! Manual intervention may be needed."
fi

sort -u "$KNOWN_HOSTS_FILE" -o "$KNOWN_HOSTS_FILE"
chown "$USER_NAME:$GROUP_NAME" "$KNOWN_HOSTS_FILE"

# =============================
# Git SSH Wrapper
# =============================
GIT_SSH_WRAPPER="$DATA_DIR/git-ssh-wrapper.sh"
echo "[*] Creating Git SSH wrapper"
cat > "$GIT_SSH_WRAPPER" <<'EOF'
#!/bin/bash
exec ssh -i "$SSH_KEY" -F "$SSH_CONFIG" "$@"
EOF
# Replace placeholders
sed -i "s|\$SSH_KEY|$SSH_KEY|g" "$GIT_SSH_WRAPPER"
sed -i "s|\$SSH_CONFIG|$SSH_CONFIG|g" "$GIT_SSH_WRAPPER"
chmod 700 "$GIT_SSH_WRAPPER"
chown "$USER_NAME:$GROUP_NAME" "$GIT_SSH_WRAPPER"

# =============================
# Install binary
# =============================
echo "[*] Installing agent binary"

# Stop service if running
if systemctl is-active --quiet "$SERVICE_NAME" 2>/dev/null; then
  echo "  [*] Stopping $SERVICE_NAME to update binary..."
  systemctl stop "$SERVICE_NAME"
fi

cp "$BIN_SOURCE" "$INSTALL_BIN"
chmod 755 "$INSTALL_BIN"
chown root:root "$INSTALL_BIN"

# =============================
# Capabilities
# =============================
setcap cap_dac_read_search,cap_sys_ptrace+ep "$INSTALL_BIN" || \
  echo "[!] setcap skipped"

# =============================
# Env file
# =============================
if [ ! -f "$CONFIG_DIR" ]; then
  echo "[*] Creating env file"
  cat > "$CONFIG_DIR" <<EOF
HORIZONX_API_URL=http://localhost:3000
HORIZONX_WS_URL=ws://localhost:3000/agent/ws
HORIZONX_SERVER_API_TOKEN=hzx_secret
HORIZONX_SERVER_ID=123
LOG_LEVEL=info
LOG_FORMAT=text
EOF
  chown "$USER_NAME:$GROUP_NAME" "$CONFIG_DIR"
  chmod 600 "$CONFIG_DIR"
fi

# =============================
# systemd service
# =============================
echo "[*] Writing systemd service"
cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=HorizonX Agent
After=network.target docker.service
Wants=docker.service

[Service]
Type=simple
User=$USER_NAME
Group=$GROUP_NAME
Environment=HOME=$HOME_DIR
Environment=GIT_SSH=$GIT_SSH_WRAPPER
EnvironmentFile=$CONFIG_DIR
ExecStart=$INSTALL_BIN --config $CONFIG_DIR
Restart=always
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=true
ProtectHome=false
ReadWritePaths=$DATA_DIR $LOG_DIR
StandardOutput=file:$LOG_DIR/agent.log
StandardError=file:$LOG_DIR/agent.error.log

[Install]
WantedBy=multi-user.target
EOF

# =============================
# Enable service
# =============================
systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
systemctl restart "$SERVICE_NAME"

echo
echo "[✓] HorizonX Agent installed"
echo "[*] Public SSH key (add to Git provider):"
echo "----------------------------------------"
cat "$SSH_KEY.pub"
echo "----------------------------------------"
echo
echo "[*] Verify known_hosts:"
wc -l "$KNOWN_HOSTS_FILE"
