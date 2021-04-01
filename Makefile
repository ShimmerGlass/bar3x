proto:
	cd lib/mirror && \
	protoc \
		-I=. \
		--go_opt=paths=source_relative \
		--go_out=. \
		--go-grpc_opt=paths=source_relative \
		--go-grpc_out=. \
		*.proto