GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_TIME = $(shell TZ=Asia/Shanghai date +'%Y-%m-%d.%H:%M:%S%Z')
BUILD_FLAGS = -ldflags "-X 'github.com/rockiecn/platform/cmd.BuildFlag=$(GIT_COMMIT)+$(BUILD_TIME)'"

platform:
	        go build $(BUILD_FLAGS) -o $@ ./


clean:
	        rm -rf platform

all: clean platform
