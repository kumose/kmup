#!/usr/bin/env bash
# This is an update script for kmup installed via the binary distribution
# from dl.kmup.com on linux as systemd service. It performs a backup and updates
# Kmup in place.
# NOTE: This adds the GPG Signing Key of the Kmup maintainers to the keyring.
# Depends on: bash, curl, xz, sha256sum. optionally jq, gpg
#   See section below for available environment vars.
#   When no version is specified, updates to the latest release.
# Examples:
#   upgrade.sh 1.15.10
#   kmuphome=/opt/kmup kmupconf=$kmuphome/app.ini upgrade.sh

# Check if kmup service is running
if ! pidof kmup &> /dev/null; then
  echo "Error: kmup is not running."
  exit 1
fi

# Continue with rest of the script if kmup is running
echo "Kmup is running. Continuing with rest of script..."

# apply variables from environment
: "${kmupbin:="/usr/local/bin/kmup"}"
: "${kmuphome:="/var/lib/kmup"}"
: "${kmupconf:="/etc/kmup/app.ini"}"
: "${kmupuser:="git"}"
: "${sudocmd:="sudo"}"
: "${arch:="linux-amd64"}"
: "${service_start:="$sudocmd systemctl start kmup"}"
: "${service_stop:="$sudocmd systemctl stop kmup"}"
: "${service_status:="$sudocmd systemctl status kmup"}"
: "${backupopts:=""}" # see `kmup dump --help` for available options

function kmupcmd {
  if [[ $sudocmd = "su" ]]; then
    # `-c` only accept one string as argument.
    "$sudocmd" - "$kmupuser" -c "$(printf "%q " "$kmupbin" "--config" "$kmupconf" "--work-path" "$kmuphome" "$@")"
  else
    "$sudocmd" --user "$kmupuser" "$kmupbin" --config "$kmupconf" --work-path "$kmuphome" "$@"
  fi
}

function require {
  for exe in "$@"; do
    command -v "$exe" &>/dev/null || (echo "missing dependency '$exe'"; exit 1)
  done
}

# parse command line arguments
while true; do
  case "$1" in
    -v | --version ) kmupversion="$2"; shift 2 ;;
    -y | --yes ) no_confirm="yes"; shift ;;
    --ignore-gpg) ignore_gpg="yes"; shift ;;
    "" | -- ) shift; break ;;
    * ) echo "Usage:  [<environment vars>] upgrade.sh [-v <version>] [-y] [--ignore-gpg]"; exit 1;; 
  esac
done

# exit once any command fails. this means that each step should be idempotent!
set -euo pipefail

if [[ -f /etc/os-release ]]; then
  os_release=$(cat /etc/os-release)

  if [[ "$os_release" =~ "OpenWrt" ]]; then
    sudocmd="su"
    service_start="/etc/init.d/kmup start"
    service_stop="/etc/init.d/kmup stop"
    service_status="/etc/init.d/kmup status"
  else
    require systemctl
  fi
fi

require curl xz sha256sum "$sudocmd"

# select version to install
if [[ -z "${kmupversion:-}" ]]; then
  require jq
  kmupversion=$(curl --connect-timeout 10 -sL https://dl.kmup.com/kmup/version.json | jq -r .latest.version)
  echo "Latest available version is $kmupversion"
fi

# confirm update
echo "Checking currently installed version..."
current=$(kmupcmd --version | cut -d ' ' -f 3)
[[ "$current" == "$kmupversion" ]] && echo "$current is already installed, stopping." && exit 0
if [[ -z "${no_confirm:-}"  ]]; then
  echo "Make sure to read the changelog first: https://github.com/go-kmup/kmup/blob/main/CHANGELOG.md"
  echo "Are you ready to update Kmup from ${current} to ${kmupversion}? (y/N)"
  read -r confirm
  [[ "$confirm" == "y" ]] || [[ "$confirm" == "Y" ]] || exit 1
fi

echo "Upgrading kmup from $current to $kmupversion ..."

pushd "$(pwd)" &>/dev/null
cd "$kmuphome" # needed for kmup dump later

# download new binary
binname="kmup-${kmupversion}-${arch}"
binurl="https://dl.kmup.com/kmup/${kmupversion}/${binname}.xz"
echo "Downloading $binurl..."
curl --connect-timeout 10 --silent --show-error --fail --location -O "$binurl{,.sha256,.asc}"

# validate checksum & gpg signature
sha256sum -c "${binname}.xz.sha256"
if [[ -z "${ignore_gpg:-}" ]]; then
  require gpg
  gpg --keyserver keys.openpgp.org --recv 7C9E68152594688862D62AF62D9AE806EC1592E2
  gpg --verify "${binname}.xz.asc" "${binname}.xz" || { echo 'Signature does not match'; exit 1; }
fi
rm "${binname}".xz.{sha256,asc}

# unpack binary + make executable
xz --decompress --force "${binname}.xz"
chown "$kmupuser" "$binname"
chmod +x "$binname"

# stop kmup, create backup, replace binary, restart kmup
echo "Flushing kmup queues at $(date)"
kmupcmd manager flush-queues
echo "Stopping kmup at $(date)"
$service_stop
echo "Creating backup in $kmuphome"
kmupcmd dump $backupopts
echo "Updating binary at $kmupbin"
cp -f "$kmupbin" "$kmupbin.bak" && mv -f "$binname" "$kmupbin"
$service_start
$service_status

echo "Upgrade to $kmupversion successful!"

popd
