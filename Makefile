PACKAGES?=$$(go list ./...)

all: test install

test:
	@go test $(PACKAGES)

install:
	@go install

fmt:
	@go fmt $(PACKAGES)

vet:
	@echo "go vet ."
	@go vet $(PACKAGES) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "go vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi
