# Bsimp

Bsimp is a minimalist audio library with multiple back-ends. It allows you to play audio files from an storage with any arbitrary directory structure.

It works with local storage or any S3 API compatible storage such as  AWS S3, DigitalOcean Spaces, Backblaze B2, Cloudflare R2 or MinIO.

## Why
I need a straightforward audio player to listen to my personal audiobooks and music. I store these audio files on my personal NAS so I can play them anywhere.

I find the original [bsimp](https://github.com/akrylysov/bsimp) project to be excellent, so I forked it to meet some of my personal requirements.

## Features
- Various back-ends such as local file systems or S3
- Cover art support
- Responsive design
- Sync global play progress rate
- Simplest password login
- Docker deployment

## Screenshots

Directory Browser

<img src="docs/directory-browser.png" alt="directory-browser" width="631"/>

Audio Player

<img src="docs/player.png" alt="player" width="631"/>

## Configuring

DigitalOcean Spaces config example:

```toml
[s3]
type = "s3"
region = "nyc3"
endpoint = "https://nyc3.digitaloceanspaces.com"
bucket = "foo"

[s3.credentials]
id = "SPACES KEY"
secret = "SPACES SECRET"
```

MinIO config example:
```toml
[s3]
type = "s3"
region = "local"
endpoint = "http://localhost:9000"
bucket = "music"
force_path_style = true

[s3.credentials]
id = "minioadmin"
secret = "minioadmin"
```

Local File System
```toml
[s3]
type = "local"
endpoint = "http://localhost:9000"
```

## Running

```sh
bsimp -config=/etc/bsimp/config.toml -http=":8080"
```

## FAQ

### What audio formats does it support?

All audio formats [supported](https://caniuse.com/?search=audio%20format) by the web browser.

### Is there a mobile app?

No, but the web interface works well on mobile phones. The Media Session API lets you control the playback from the notification bar or the lock screen.

### Does it support playlists?

No, Bsimp follows the directory structure.

### Does it support transcoding?

No, audio files are streamed as is.
