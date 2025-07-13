build:
	go build -o ai-cat-insta .

run:
	go run .

test-video:
	go run . test-video

clean:
	rm -f ai-cat-insta

install:
	go mod tidy

.PHONY: build run test-video clean install