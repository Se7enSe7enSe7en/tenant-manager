# Single VPS deployment

- As of now we are deploying this project in a single VPS
- In the future (far future) when this project becomes multi VPSed project (or multi instance?), we will change how we deploy because this won't be feasable to do to every single VPS every time

## Step 1: One time setup in VPS

### Step 1.1: create "deploy" user give permission to restart the app service with sudo

```sh
sudo adduser deploy

# option 1: for security
echo "deploy ALL=(ALL) NOPASSWD: /bin/systemctl * tenant-manager" | sudo tee /etc/sudoers.d/deploy
# option 2: alternatively, give deploy root access
sudo usermod -aG sudo deploy
```

### notes:

- `* tenant-manager` means we are giving user (deploy) the permission to do anything with that service (eg. stop tenant-manager, restart tenant-manager, etc.)
- "tenant-manager" is the current app name, replace with actual name if changed in the future
- for security, we only give deploy user the permission it needs, alternatively we could've
  give deploy root access with `sudo usermod -aG sudo deploy` to make it simpler (but less secure)

```
deploy          ALL=(ALL)   NOPASSWD:   /bin/systemctl restart tenant-manager
  │              │           │                │
  │              │           │                └── ONLY this exact command
  │              │           └── no password prompt needed
  │              └── on any host, as any user (root by default)
  └── the user this rule applies to
```

### verify step 1.1

```sh
sudo systemctl restart tenant-manager # should work, no password needed
sudo whoami # "Sorry, user deploy is not allowed to execute..."

sudo -l -U deploy # shows what is deploy allowed to run
```

### Step 1.2: Set up SSH key auth (go back to your terminal runing from your machine)

```sh
ssh-copy-id deploy@<vps-ip-address>
```

- make sure to go back to your terminal in this step (exit ssh session, type `exit`, if you are sshed)
- "vps-ip-address" replace with actual
- this step:
  1. finds your public key (the most recent ~/.ssh/id\*.pub file)
  2. SSHs then prompts you for the deploy user's password one time
  3. Appends the public key to ~/.ssh/authorized_keys on the VPS
  4. Sets the correct permissions

### Verify step 1.2:

```
ssh deploy@<vps-ip-address>
```

### Step 1.3: Create the app directory

```
sudo mkdir -p /opt/tenant-manager
sudo chown deploy:deploy /opt/tenant-manager
```

#### notes

- fun fact: in linux/macos "/opt" (optional) is the directory where "optional" apps go, it is optional because these are manually added programs (add ons) by the user, since usually programs are added using a package manager (eg. apt, brew)

```
sudo chown deploy:deploy /opt/myapp
#           ─────── ───────
#           owner   group
```

- chown (change owner), it sets who owns a file or directory

### Step 1.4: Install Docker (if not yet installed)

```sh
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker deploy
```

#### notes

```
sudo usermod -aG docker deploy
#            ── ── ────── ──────
#            │  │  group  user being added
#            │  └─ append (don't remove existing groups)
#            └── change group
```

### Step 1.5: Start PostgreSQL in docker using compose file

in the repo in your dev terminal

```sh
task setup-prod
```

this will copy paste the prod docker compose file into the VPS using scp
then in your VPS you can run the docker compose file with `docker compose up`

```sh
docker compose -f "location/of/compose/file" up
```

### Step 1.5 (ALTERNATE MANUAL WAY): Start PostgreSQL in Docker

```sh
docker run -d \
--name tenant-manager-postgres-db \
--restart unless-stopped \
-e POSTGRES_USER=postgres \
-e POSTGRES_PASSWORD=<the-password> \
-e POSTGRES_DB=db \
-v postgres-data:/var/lib/postgresql \
-p 127.0.0.1:5432:5432 \
postgres:18-alpine
```

### Step 1.6: Install Caddy

```sh
sudo apt install caddy -y
```

- Caddy automatically obtains and renews TLS certificates for all your sites.
- caddy: https://caddyserver.com/

### Step 1.7: create the app service (tenant-manager.service)

```sh
sudo tee /etc/systemd/system/tenant-manager.service << 'EOF'
[Unit]
Description=Tenant manager go app
After=network.target

[Service]
User=deploy
WorkingDirectory=/opt/tenant-manager
ExecStart=/opt/tenant-manager/tenant-manager
EnvironmentFile=/opt/tenant-manager/.env
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable tenant-manager
```

### Step 1.8: Create .env file

```sh
tee > /opt/tenant-manager/.env << 'EOF'
DATABASE_URL=postgres://postgres:<the-password>@127.0.0.1:5432/db?sslmode=disable
EOF

sudo chown deploy:deploy /opt/tenant-manager/.env
chmod 600 /opt/tenant-manager/.env
```

### Step 1.9: Reverse proxy setup (caddy)

```sh
sudo tee /etc/caddy/Caddyfile << 'EOF'
mydomain.com {
  reverse_proxy 127.0.0.1:8080
}
EOF

sudo systemctl reload caddy
```

- change "mydomain.com" to actual domain

#### if no domain yet, set it to port 80

```sh
sudo tee /etc/caddy/Caddyfile << 'EOF'
:80 {
  reverse_proxy 127.0.0.1:8080
}
EOF

sudo systemctl reload caddy
```

- to access just go to http://vps-ip:80/
- notice its http and not https, thats because Let's Encrypt (free automated Certificate Authority) wont create a TLS certificate for you (no TLS cert no https), but it will be a working reverse proxy using http

### Step 1.10: configure firewall (UFW)

```sh
sudo ufw allow OpenSSH
sudo ufw allow http
sudo ufw allow https
sudo ufw enable
```

- allowing OpenSSH means allowing port 22/tcp
- allowing http means allowing port 80
- allowing https means allowing port 443

## Step 2: run task deploy

```sh
task deploy
```

- this will:
  1. build the app binary
  2. upload the binary through scp (ssh copy) to the VPS
  3. migrate the app to the latest using embedded migration scripts
  4. restart the app service (tenant-manager.service)
