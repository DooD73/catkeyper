APP_NAME := CatKeyper
BUNDLE_ID := com.local.catkeyper
VERSION ?= 1.0.0
BUILD_DIR := build
APP_DIR := $(BUILD_DIR)/$(APP_NAME).app
BIN_NAME := catkeyper
ICONSET := $(BUILD_DIR)/AppIcon.iconset
ICNS := assets/AppIcon.icns
ZIP := $(BUILD_DIR)/$(BIN_NAME)-$(VERSION)-macos.zip
DMG := $(BUILD_DIR)/$(BIN_NAME)-$(VERSION)-macos.dmg
DMG_BACKGROUND := assets/dmg-background.png
DMG_STAGING_DIR := $(BUILD_DIR)/.dmg-staging

.PHONY: build run test icon package adhoc-sign zip dmg release clean

build:
	mkdir -p "$(BUILD_DIR)"
	go build -trimpath -ldflags "-s -w -X main.appVersion=$(VERSION)" -o $(BUILD_DIR)/$(BIN_NAME) .

run:
	go run .

test:
	go test -v ./...

icon:
	mkdir -p "$(BUILD_DIR)"
	go run ./cmd/iconrender -in assets/openmoji-cat-face.svg -out "$(ICONSET)"
	iconutil -c icns "$(ICONSET)" -o "$(ICNS)"

package: build
	test -f "$(ICNS)"
	rm -rf "$(APP_DIR)"
	mkdir -p "$(APP_DIR)/Contents/MacOS" "$(APP_DIR)/Contents/Resources"
	printf '%s\n' \
		'<?xml version="1.0" encoding="UTF-8"?>' \
		'<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">' \
		'<plist version="1.0">' \
		'<dict>' \
		'  <key>CFBundleExecutable</key>' \
		'  <string>$(BIN_NAME)</string>' \
		'  <key>CFBundleIdentifier</key>' \
		'  <string>$(BUNDLE_ID)</string>' \
		'  <key>CFBundleName</key>' \
		'  <string>$(APP_NAME)</string>' \
		'  <key>CFBundleDisplayName</key>' \
		'  <string>$(APP_NAME)</string>' \
		'  <key>CFBundlePackageType</key>' \
		'  <string>APPL</string>' \
		'  <key>CFBundleIconFile</key>' \
		'  <string>AppIcon</string>' \
		'  <key>CFBundleShortVersionString</key>' \
		'  <string>$(VERSION)</string>' \
		'  <key>CFBundleVersion</key>' \
		'  <string>$(VERSION)</string>' \
		'  <key>LSMinimumSystemVersion</key>' \
		'  <string>12.0</string>' \
		'  <key>NSHighResolutionCapable</key>' \
		'  <true/>' \
		'</dict>' \
		'</plist>' > "$(APP_DIR)/Contents/Info.plist"
	cp "$(BUILD_DIR)/$(BIN_NAME)" "$(APP_DIR)/Contents/MacOS/$(BIN_NAME)"
	cp "$(ICNS)" "$(APP_DIR)/Contents/Resources/AppIcon.icns"
	chmod +x "$(APP_DIR)/Contents/MacOS/$(BIN_NAME)"

adhoc-sign: package
	codesign --force --deep --sign - "$(APP_DIR)"
	codesign --verify --deep --strict --verbose=2 "$(APP_DIR)"

zip: adhoc-sign
	rm -f "$(ZIP)"
	ditto -c -k --keepParent "$(APP_DIR)" "$(ZIP)"

dmg: adhoc-sign
	@set -eu; \
	staging_dir="$(DMG_STAGING_DIR)"; \
	rm -rf "$$staging_dir"; \
	mkdir -p "$$staging_dir"; \
	trap 'rm -rf "$$staging_dir"' EXIT HUP INT TERM; \
	ditto "$(APP_DIR)" "$$staging_dir/$(APP_NAME).app"; \
	rm -f "$(DMG)"; \
	create-dmg \
		--volname "$(APP_NAME)" \
		--volicon "$(ICNS)" \
		--background "$(DMG_BACKGROUND)" \
		--window-pos 200 120 \
		--window-size 660 400 \
		--text-size 13 \
		--icon-size 116 \
		--icon "$(APP_NAME).app" 175 235 \
		--hide-extension "$(APP_NAME).app" \
		--app-drop-link 485 235 \
		--no-internet-enable \
		--format UDZO \
		--filesystem HFS+ \
		"$(DMG)" \
		"$$staging_dir"

release: zip dmg

clean:
	rm -rf "$(BUILD_DIR)"
