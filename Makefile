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
SIGN_IDENTITY ?=
NOTARY_PROFILE ?=

.PHONY: build run icon package adhoc-sign zip dmg sign notarize release clean

build:
	mkdir -p "$(BUILD_DIR)"
	go build -trimpath -ldflags "-s -w -X main.appVersion=$(VERSION)" -o $(BUILD_DIR)/$(BIN_NAME) .

run:
	go run .

icon:
	mkdir -p "$(BUILD_DIR)"
	go run ./cmd/iconrender -in assets/openmoji-cat-face.svg -out "$(ICONSET)"
	iconutil -c icns "$(ICONSET)" -o "$(ICNS)"

package: build icon
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
	rm -f "$(DMG)"
	hdiutil create -volname "$(APP_NAME)" -srcfolder "$(APP_DIR)" -ov -format UDZO "$(DMG)"

sign: package
	@test -n "$(SIGN_IDENTITY)" || (echo 'SIGN_IDENTITY is required, for example: make sign SIGN_IDENTITY="Developer ID Application: Your Name (TEAMID)"' && exit 1)
	codesign --force --deep --options runtime --timestamp --sign "$(SIGN_IDENTITY)" "$(APP_DIR)"
	codesign --verify --deep --strict --verbose=2 "$(APP_DIR)"

notarize: sign
	@test -n "$(NOTARY_PROFILE)" || (echo 'NOTARY_PROFILE is required, for example: make notarize NOTARY_PROFILE=catkeyper-notary' && exit 1)
	rm -f "$(ZIP)"
	ditto -c -k --keepParent "$(APP_DIR)" "$(ZIP)"
	xcrun notarytool submit "$(ZIP)" --keychain-profile "$(NOTARY_PROFILE)" --wait
	xcrun stapler staple "$(APP_DIR)"
	ditto -c -k --keepParent "$(APP_DIR)" "$(ZIP)"

release: zip dmg

clean:
	rm -rf "$(BUILD_DIR)"
