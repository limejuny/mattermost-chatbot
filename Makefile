BINARY_NAME=plugin

.PHONY: default
default: build deploy

.PHONY: build
build:
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) -v
	@tar -czf bin/$(BINARY_NAME).tar.gz plugin.json assets/icon.png site -C bin/ $(BINARY_NAME)

.PHONY: deploy
deploy:
	@bin/mmctl plugin delete com.github.limejuny.mattermost-chatbot
	@bin/mmctl plugin add bin/$(BINARY_NAME).tar.gz
	@bin/mmctl plugin enable com.github.limejuny.mattermost-chatbot

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf bin/${BINARY_NAME} bin/$(BINARY_NAME).tar.gz
