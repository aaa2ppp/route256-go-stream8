generate-all:
	cd loms && TARGET_OS=linux TARGET_ARCH=amd64 make generate
	cd cart && TARGET_OS=linux TARGET_ARCH=amd64 make generate
	cd notifications && TARGET_OS=linux TARGET_ARCH=amd64 make generate

build-all: generate-all
	cd loms && TARGET_OS=linux TARGET_ARCH=amd64 make build
	cd cart && TARGET_OS=linux TARGET_ARCH=amd64 make build
	cd notifications && TARGET_OS=linux TARGET_ARCH=amd64 make build

run-all: run-db build-all
	docker-compose up --force-recreate --build

run-db:
	docker-compose up db -d


precommit:
	cd loms && make precommit
	cd cart && make precommit
	cd notifications && make precommit

clean-all:
	cd loms && make clean
	cd cart && make clean
	cd notifications && make clean
