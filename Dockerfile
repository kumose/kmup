# syntax=docker/dockerfile:1
# Build stage - compile kmup binary from source
FROM docker.io/library/golang:1.25-alpine3.22 AS build-env

# Go module proxy (Alibaba Cloud mirror for faster download in China, fallback to direct if proxy missing)
ARG GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# Build arguments for kmup compilation
ARG KMUP_VERSION
ARG TAGS="sqlite sqlite_unlock_notify"
ENV TAGS="bindata timetzdata $TAGS"
ARG CGO_EXTRA_CFLAGS

# Fix 1: Correct BusyBox sed syntax for Alpine repo replacement
# Use double quotes + ensure repo path matches Alpine 3.22 structure
RUN sed -i "s|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g" /etc/apk/repositories

# Fix 2: Install build dependencies (pnpm via npm instead of apk)
# apk add nodejs (includes npm) → then install pnpm globally (more reliable)
RUN apk --no-cache add \
    build-base \
    git \
    nodejs \
    npm && \
    # Install pnpm globally (fixes "pnpm not found" in Alpine 3.22)
    npm install -g pnpm && \
    # Configure pnpm/npm registry (npmmirror for faster frontend dependency download)
    pnpm config set registry https://registry.npmmirror.com && \
    npm config set registry https://registry.npmmirror.com

# Set working directory for kmup source code
WORKDIR ${GOPATH}/src/github.com/kumose/kmup

# Copy source code (exclude .git - mounted separately for version data)
# Avoid mount for node_modules (platform-dependent content needs exclusion)
COPY --exclude=.git/ . .

# Build kmup binary with cached dependencies (speed up rebuilds)
# - /go/pkg/mod: Go module cache
# - /root/.cache/go-build: Go build cache
# - /root/.local/share/pnpm/store: pnpm dependency cache
# - .git mount: Required for version metadata (commit hash/version)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target="/root/.cache/go-build" \
    --mount=type=cache,target=/root/.local/share/pnpm/store \
    --mount=type=bind,source=".git/",target=".git/" \
    make

# Copy Docker runtime configs to temp directory
COPY docker/root /tmp/local

# Fix executable permissions (Windows builds strip execute bits from files)
# Targets: entrypoint script, s6 process manager scripts, kmup binary
RUN chmod 755 /tmp/local/usr/bin/entrypoint \
              /tmp/local/usr/local/bin/* \
              /tmp/local/etc/s6/kmup/* \
              /tmp/local/etc/s6/openssh/* \
              /tmp/local/etc/s6/.s6-svscan/* \
              /go/src/github.com/kumose/kmup/kmup

# Runtime stage - minimal Alpine image for kmup execution
FROM docker.io/library/alpine:3.22 AS kmup

# Configure Alpine apk mirror (Alibaba Cloud) for runtime dependencies
RUN sed -i "s|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g" /etc/apk/repositories

# Expose ports: 22 (SSH for git access), 3326 (kmup web interface)
EXPOSE 22 3326

# Install runtime dependencies (minimal set for security/performance)
# bash: Shell for entrypoint scripts
# ca-certificates: SSL/TLS certificates for HTTPS
# curl: HTTP requests (health checks/config fetch)
# gettext: Environment variable substitution in config files
# git: Git binary for repository operations
# linux-pam: PAM authentication support
# openssh: SSH server for git clone/push
# s6: Lightweight process manager (manage kmup + sshd)
# sqlite: Embedded database (default for kmup)
# su-exec: Safe user switching (avoid root execution)
# gnupg: GPG signature verification (optional for secure git)
RUN apk --no-cache add \
    bash \
    ca-certificates \
    curl \
    gettext \
    git \
    linux-pam \
    openssh \
    s6 \
    sqlite \
    su-exec \
    gnupg

# Create git user/group (UID/GID 1000) - run kmup as non-root (security best practice)
# -S: System group/user
# -H: No home directory (mount /data/git instead)
# -D: Disable password login (SSH key only)
# -h: Home directory for git repositories
# -s: Default shell (bash for git operations)
RUN addgroup \
    -S -g 1000 \
    git && \
  adduser \
    -S -H -D \
    -h /data/git \
    -s /bin/bash \
    -u 1000 \
    -G git \
    git && \
  # Disable password login for git user (* = no valid password)
  echo "git:*" | chpasswd -e

# Copy runtime configs and compiled kmup binary from build stage
COPY --from=build-env /tmp/local /
COPY --from=build-env /go/src/github.com/kumose/kmup/kmup /app/kumose/kmup

# Environment variables for runtime
# USER: Default user to run kmup (non-root git user)
# KMUP_CUSTOM: Custom config directory (mounted as volume)
ENV USER=git
ENV KMUP_CUSTOM=/data/kmup

# Persistent volume for kmup data (repositories, config, database, logs)
VOLUME ["/data"]

# Entrypoint script: Initialize environment, switch user, start services
ENTRYPOINT ["/usr/bin/entrypoint"]

# Default command: Start s6 process manager (manages kmup + sshd)
CMD ["/usr/bin/s6-svscan", "/etc/s6"]