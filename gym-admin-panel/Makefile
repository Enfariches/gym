# Determine shell type
ifeq ($(OS),Windows_NT)
SHELL_TYPE = windows
else
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
SHELL_TYPE = linux
endif
ifeq ($(UNAME_S),Darwin)
SHELL_TYPE = darwin
endif
endif

# Set variables based on shell type
ifeq ($(SHELL_TYPE),windows)
PROTO_FILES = ../api/v1/auth/auth.proto ../api/v1/users/admin.proto ../api/v1/users/employee.proto ../api/v1/schedule/schedule.proto ../api/v1/media/media.proto
MAKE_DIR = mkdir protogen
RM_CMD = if exist protogen rmdir /s /q protogen
else
PROTO_FILES = $(shell find ../api/v1 -name '*.proto')
MAKE_DIR = mkdir -p ./protogen
RM_CMD = rm -rf ./protogen
endif

# Генерация gRPC TypeScript Интерфейсов. 
gen:
	$(RM_CMD)
	$(MAKE_DIR)
	npx protoc -I ../api $(PROTO_FILES) \
		--ts_out ./protogen \
		--ts_opt ts_nocheck \
		--ts_opt eslint_disable

clean:
	$(RM_CMD)

