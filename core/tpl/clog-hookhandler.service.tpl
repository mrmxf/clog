# clog Hookhandler system installer ©2025 Mr MXF Ltd
# install in /etc/systemd/system as clog-hookhandler.service
#            sudo systemctl daemon-reload
#            sudo systemctl enable clog-hookhandler.service
#            sudo systemctl start  clog-hookhandler.service

[Unit]
Description=clog-hookhandler: Listen for webhooks and take actions
After=network.target

[Install]
WantedBy=multi-user.target

# start clog in the admin folder. The @ prefix allows an argument to the exe

[Service]
Type=simple
# ensure this file exists
EnvironmentFile=/home/admin/clogrc/clog-hookhandler.service.env
ExecStart=/usr/local/bin/clog ${ARG1} ${ARG2} {$ARG3} {$ARG4}
WorkingDirectory=/home/admin
Restart=always
RestartSec=5

# Send output to systemd journal
StandardOutput=journal
StandardError=journal

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
#ProtectHome=true
ReadWritePaths=/home/admin

# Network security
RestrictAddressFamilies=AF_INET AF_INET6
