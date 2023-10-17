# Watch and run chatbot program
chatbot.dev:
	gow run ./cmd/chatbot

# Builds chatbot binary named ./bin/chatbot
chatbot.build:
	go build -o ./bin/ ./cmd/chatbot 

# Build chatbot and run chatbot binary ./bin/chatbot
chatbot.run: chatbot.build
	./bin/chatbot
	
