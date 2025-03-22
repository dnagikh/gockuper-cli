# gockuper-cli

**gockuper-cli** is a simple and powerful CLI utility for PostgreSQL backups, built in Go.  
It creates compressed backups, uploads them to cloud or local storage, and rotates old archives automatically.

---

## ✅ Currently Supported

**Databases**
- ✅ PostgreSQL

**Storage**
- ✅ Local file system
- ✅ Dropbox

**Compression**
- ✅ Gzip (`.gz`)
- ✅ None

**Rotation**
- ✅ Max backup count (`MAX_BACKUPS`)

---

## 🛠️ Coming Soon

**Databases**
- [ ] MySQL

**Storage**
- [ ] Google Drive
- [ ] Yandex Disk
- [ ] Mail.ru Cloud

---

## 📦 Usage

```shell
gockuper-cli

gockuper-cli is a CLI utility to backup your database,
compress it, upload to storage, and clean up old backups.

Usage:
  gockuper-cli [command]

Available Commands:
  backup       Create and upload a new backup
  daemon       Run daemon for refreshing cloud service tokens
  help         Help about any command

Flags:
  -h, --help   Help for gockuper-cli

Use "gockuper-cli [command] --help" for more information about a command.
```

## ⚙️ Configuration

Configuration is loaded from a .env file or system environment variables.

### 🧱 Base config

```
DB_TYPE=postgres             # only "postgres" is supported for now
DB_HOST=localhost
DB_PORT=5432
DB_NAME=db_name
DB_USER=username
DB_PASSWORD=password
```

### ☁️ Storage

```
STORAGE_TYPE=dropbox         # or "file"
STORAGE_FILE_PATH=/          # Dropbox path or local folder (Must be / for root)
MAX_BACKUPS=3                # Max number of backups to keep

DROPBOX_API_TOKEN=           # required for Dropbox
```

### 🗜 Compression

```
COMPRESS=gzip                # gzip or none
```

### 📃 Logging
```
LOG_NAME=./gockuper.log      # log file name
LOG_MAX_SIZE=50              # max size (MB)
LOG_MAX_BACKUPS=3            # rotated backups to keep
LOG_MAX_AGE=90               # max age (days)
LOG_COMPRESSION=true         # compress rotated logs
LOG_LEVEL=info               # debug, info, warn, error
LOG_TARGET=stdout            # or file
```

## 🚀 Installation

Install via go install:
```
go install github.com/your-org/gockuper-cli@latest
```
Or download prebuilt binary from the Releases page.

Make sure the binary is executable:
```
chmod +x ./gockuper-cli
```

## 🧪 Development
```
make build        # Build the binary
make lint         # Run golangci-lint
make test         # Run unit tests
```

## ⏱ Example CRON Job

Add to your crontab (crontab -e):
```
0 3 * * * /path_to_dir/gockuper-cli backup > /dev/null 2>&1
```

## 🪪 License
[![MIT](https://github.com/dnagikh/gockuper-cli/blob/main/LICENSE)]


# 🔐 Dropbox OAuth Setup Guide for gockuper-cli

This guide will help you generate the required `access_token` and `refresh_token` to enable automatic token refresh with Dropbox.

---

## Step-by-step Instructions

### 1. Create Dropbox App

1. Visit: [https://www.dropbox.com/developers/apps](https://www.dropbox.com/developers/apps)
2. Click **Create App**
3. Choose:
   - **Scoped access**
   - **App folder** (recommended)
4. Name your app (e.g., `gockuper-backups`) and click **Create App**

---

### 2. Allow Offline Access

1. Go to the app's settings
2. Scroll to **OAuth 2** section
3. Make sure the following scopes are enabled in **Permissions**:
   - `files.content.write`
   - `files.metadata.read`
4. Take note of your:
   - **App key** (client_id)
   - **App secret** (client_secret)

---

### 3. Generate Authorization Code

In your browser, open:

```
https://www.dropbox.com/oauth2/authorize?client_id=<APP_KEY>&token_access_type=offline&response_type=code
```

- Replace `APP_KEY` with your app key
- Login and approve access
- You'll be redirected to a blank page with a `code=XYZ...`

Copy that `code`.

---

### 4. Exchange Code for Tokens

Run this `curl` command:

```bash
curl -X POST https://api.dropboxapi.com/oauth2/token \
  -d code=<CODE> \
  -d grant_type=authorization_code \
  -d client_id=<APP_KEY> \
  -d client_secret=<APP_SECRET>
```

In response, you will get:

```json
{
  "access_token": "...",
  "expires_in": 14400,
  "refresh_token": "...",
  ...
}
```

---

### 5. Create `gockuper_token.json`

Save the data into a file:

```json
{
  "access_token": "...",
  "refresh_token": "...",
  "expires_at": "2025-03-22T18:00:00Z"
}
```
The expires_at field must be formatted using the RFC3339 standard.

---

### 6. Configure `.env`

```env
DROPBOX_CLIENT_ID=your_app_key
DROPBOX_CLIENT_SECRET=your_app_secret
DROPBOX_TOKEN_FILE=~/.gockuper/token.json
```

Or leave `DROPBOX_TOKEN_FILE` unset and it will default to `token.json` in current working directory.

---

### Done!

Now `gockuper-cli` will:
- Use the `access_token` to talk to Dropbox
- Auto-refresh it in the background using `refresh_token`

No need to manually update it ever again.
