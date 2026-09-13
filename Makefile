APP_NAME := Hammock
BUNDLE_DIR := $(APP_NAME).app

.PHONY: build-mac
build-mac:
	@echo "Building for macos-arm64 app bundle..."
	rm -rf $(BUNDLE_DIR)
	mkdir -p $(BUNDLE_DIR)/Contents/MacOS
	mkdir -p $(BUNDLE_DIR)/Contents/Resources

	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o $(BUNDLE_DIR)/Contents/MacOS/hammock ./cmd/hammock

	cp build/macos/Info.plist $(BUNDLE_DIR)/Contents/Info.plist
	if [ -f build/macos/AppIcon.icns ]; then cp build/macos/AppIcon.icns $(BUNDLE_DIR)/Contents/Resources/; fi

	codesign --force --deep --sign - $(BUNDLE_DIR)
	@echo "yay we built $(BUNDLE_DIR)!"

.PHONY: clean
clean:
	rm -rf $(BUNDLE_DIR) $(APP_NAME)-macOS-arm64.zip
