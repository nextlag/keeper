postgres:
	docker run --name keeper -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=Xer_0101 -d postgres:14-alpine
	sleep 10
	$(MAKE) db_create
	sleep 1
	$(MAKE) db_uuid
db_create:
	docker exec -it keeper psql -U postgres -c "CREATE DATABASE keeper;"
db_uuid:
	docker exec -it keeper psql -U postgres -d keeper -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"

# Key constant
KEY := key.pem
BASE64 := key_base64
INLINE := key_base64_inline

PUBLIC := public_key.pem
PUBLIC_BASE64 := public_key_base64
PUBLIC_INLINE := public_key_base64_inline

collectKeys:
	$(MAKE) key_pem && sleep 1 && $(MAKE) base64 && sleep 1 && $(MAKE) inline && sleep 1 && $(MAKE) public_pem && sleep 1 && \
	$(MAKE) public_base64 && sleep 1 && $(MAKE) public_inline && sleep 1 && $(MAKE) base64

key_pem:
	openssl genpkey -algorithm RSA -out $(KEY)
base64:
	openssl base64 -in $(KEY) -out $(BASE64)
inline:
	cat $(BASE64) | tr -d '\n' > $(INLINE)
public_pem:
	openssl rsa -in $(KEY) -pubout -out $(PUBLIC)
public_base64:
	openssl base64 -in $(PUBLIC) -out $(PUBLIC_BASE64)
public_inline:
	cat $(PUBLIC_BASE64) | tr -d '\n' > $(PUBLIC_INLINE)

clean:
	rm -f $(KEY) $(BASE64) $(INLINE) $(PUBLIC) $(PUBLIC_BASE64) $(PUBLIC_INLINE)

# Linter constants
LINTER := golangci-lint
lint:
	@echo === Lint
	$(LINTER) --version
	$(LINTER) cache clean && $(LINTER) run

generate:
	go generate ./...

reload_postgres:
	docker stop keeper > /dev/null
	docker rm keeper > /dev/null
	$(MAKE) postgres

test:
	go test ./...


# Download files
URL := http://localhost:8080/api/v1/auth/login
EMAIL := example@mail.ru
PASSWORD := 12345
FILE_ID := e9b45ca4-a92c-46a4-9590-af5bc950060a # current uuid is required
DOWNLOAD_URL := http://localhost:8080/api/v1/user/binary/$(FILE_ID)

# Define variables to store tokens
ACCESS_TOKEN := $(shell echo $$(curl -s -X POST --location "$(URL)" -d '{"email": "$(EMAIL)", "password": "$(PASSWORD)"}' | jq -r '.access_token'))
REFRESH_TOKEN := $(shell echo $$(curl -s -X POST --location "$(URL)" -d '{"email": "$(EMAIL)", "password": "$(PASSWORD)"}' | jq -r '.refresh_token'))

# Target to download the file using the access token
download:
	curl -X GET "$(DOWNLOAD_URL)" -H "Authorization: Bearer $(ACCESS_TOKEN)" -o downloaded_file.txt

.PHONY: db_create db_uuid postgres lint generate reload_postgres test key_pem base64 inline public_pem public_base64 public_inline collectKeys clean download
