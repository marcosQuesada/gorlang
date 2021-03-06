
test:
	rm -rf pkg/
	@echo "Test suite"
	@find src/actor -not -path "*.*" | awk '{print "./" $$0 "/..."}' | xargs go test -v
	@find src/register -not -path "*.*" | awk '{print "./" $$0 "/..."}' | xargs go test -v
	@find src/supervisor -not -path "*.*" | awk '{print "./" $$0 "/..."}' | xargs go test -v
	@go test -v

