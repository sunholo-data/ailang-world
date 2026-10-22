#!/bin/sh
# Build the release binary and assemble World.app by hand (no tauri-cli needed).
# Run ./fetch-fonts.sh first: the fonts are fetched, not tracked.
set -eu
cd "$(dirname "$0")"
# The data, also as a script, so the page renders even where fetch() of a bundled file is refused.
{ printf 'window.SPIKE_DATA = '; cat dist/spike-data.json; printf ';\n'; } > dist/data.js
cargo build --release
APP=World.app
rm -rf "$APP"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"
cp target/release/world "$APP/Contents/MacOS/World"
# Icon: the AILANG mark, rendered to an .icns, when icons/icon.png exists (it is not tracked:
# the repo's tracked-binary gate keeps PNGs out of git). Without it the app uses the default icon.
if [ -f icons/icon.png ]; then
SET=$(mktemp -d)/World.iconset; mkdir -p "$SET"
for s in 16 32 128 256 512; do
  sips -z $s $s icons/icon.png --out "$SET/icon_${s}x${s}.png" >/dev/null
  d=$((s*2)); [ $d -le 512 ] && sips -z $d $d icons/icon.png --out "$SET/icon_${s}x${s}@2x.png" >/dev/null
done
cp icons/icon.png "$SET/icon_512x512@2x.png"
iconutil -c icns "$SET" -o "$APP/Contents/Resources/World.icns"
fi
cat > "$APP/Contents/Info.plist" <<'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>CFBundleName</key><string>World</string>
  <key>CFBundleDisplayName</key><string>World</string>
  <key>CFBundleExecutable</key><string>World</string>
  <key>CFBundleIdentifier</key><string>com.sunholo.ailang.world.spike</string>
  <key>CFBundleIconFile</key><string>World</string>
  <key>CFBundlePackageType</key><string>APPL</string>
  <key>CFBundleShortVersionString</key><string>0.1.0</string>
  <key>CFBundleVersion</key><string>1</string>
  <key>LSMinimumSystemVersion</key><string>11.0</string>
  <key>NSHighResolutionCapable</key><true/>
  <key>NSPrincipalClass</key><string>NSApplication</string>
</dict></plist>
PLIST
codesign --force --deep -s - "$APP" 2>/dev/null || true
du -sh "$APP"
