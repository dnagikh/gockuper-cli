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


## 🔐 How to Configure Dropbox Access

To enable backup uploads to Dropbox, you need to generate an access token for a scoped app.

### ✅ Step-by-step:

1. Go to the Dropbox App Console:  
   👉 https://www.dropbox.com/developers/apps

2. Click **"Create App"**

3. Choose:
    - **Scoped access**
    - **App folder** (Recommended: limits the app to a single folder inside `/Apps/`)

4. Give your app a name (e.g., `backups`) and click **Create App**

5. Go to the **Permissions** tab and enable:
    - ✅ `files.content.write` *(required for uploading)*
    - ❌ Do **not** enable `files.content.delete` *(for safety — no deletions)*

6. Go to the **OAuth 2** tab and click **"Generate access token"** (only after scopes set up)

7. Copy the generated token and paste it into your `config.env` file:

```env
DROPBOX_API_TOKEN=your_generated_token
STORAGE_TYPE=dropbox
STORAGE_FILE_PATH=/backups  # or just "/" for root inside app folder
```