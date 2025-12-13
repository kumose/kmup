# Kmup

## build

```shell
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

run docker
```bash
docker run -d \
  --name kmup-server \          # Name the container (easier to manage)
  -p 2222:22 \                  # Map host port 2222 → container port 22 (SSH for git)
  -p 3326:3326 \                # Map host port 3326 → container port 3326 (kmup web UI)
  -v /host/path/kmup-data:/data \ # Persist /data volume (critical: saves repos/config/database)
  --restart unless-stopped \    # Auto-restart container if it crashes
  ghcr.io/kumose/kmup:latest
```

run docker
```bash
docker run -d --name kmup-server -p 2222:22 -p 3326:3326 -v /host/path/kmup-data:/data --restart unless-stopped ghcr.io/kumose/kmup:latest
```
