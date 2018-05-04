#@IgnoreInspection BashAddShebang
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export CGO_ENABLED= 0
export GOOS=linux
export ENV=development
export DATABASE_URI=postgres://realworld:secret@127.0.0.1:5432/realworld?sslmode=disable
export DATABASE_URI_TEST=postgres://realworld:secret@127.0.0.1:5432/realworld_test?sslmode=disable
export LOG_LEVEL=DEBUG
export VERSION=dev
export COVERAGE_DIR=$(ROOT)/coverage
export GLIDE_HOME=$(HOME)/.glide
export SERVER_ADDRESS=127.0.0.1:8585
export APP=golang-realword
export LDFLAGS="-w -s -X main.BuildTime=`date -u +%Y/%m/%d_%H:%M:%S` -X main.BuildID=`git rev-parse HEAD` -X main.Version=`git tag -l --points-at HEAD`"
export DEBUG= 1

all: lint build citest

fetch: glide-install

contributors:
	git log --all --format='%aN <%cE>' | sort -u  > CONTRIBUTORS

#######
# Build
#######

build:
	go build -v -o bin/$(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) *.go
	bash build-workers.sh

install: fetch
	go install -v -a -installsuffix cgo -ldflags $(LDFLAGS) *.go

run:
	go run -ldflags $(LDFLAGS) cmd/main.go

run-mac:
	GOOS=darwin go run -ldflags $(LDFLAGS) cmd/main.go

######
# Lint
######

check-gometalinter:
	which gometalinter || (go get -u -v github.com/alecthomas/gometalinter && gometalinter --install --update)

fast-lint:
	gometalinter --disable-all --vendor --enable=gofmt $(ROOT)/...
	gometalinter --disable-all --vendor --enable=vet --exclude=pkg/axmlParser $(ROOT)/...
	gometalinter --disable-all --vendor --enable=vetshadow $(ROOT)/...
	gometalinter --disable-all --vendor --enable=gocyclo --cyclo-over=36 $(ROOT)/...
	gometalinter --disable-all --vendor --enable=golint --exclude=mocks --exclude=pkg/axmlParser --exclude=pkg/faktory $(ROOT)/...
	gometalinter --disable-all --vendor --enable=ineffassign $(ROOT)/...
	gometalinter --disable-all --vendor --enable=misspell $(ROOT)/...

lint:
	gometalinter --disable-all --vendor --enable=deadcode $(ROOT)/...
	gometalinter --disable-all --vendor --enable=errcheck --exclude=mocks/ $(ROOT)/...
	gometalinter --disable-all --vendor --enable=gas $(ROOT)/...
	gometalinter --disable-all --vendor --enable=goconst $(ROOT)/...
	gometalinter --disable-all --vendor --enable=gocyclo --cyclo-over=36 $(ROOT)/...
	gometalinter --disable-all --vendor --enable=gofmt $(ROOT)/...
	gometalinter --disable-all --vendor --enable=goimports $(ROOT)/...
	gometalinter --disable-all --vendor --enable=golint --exclude=mocks  $(ROOT)/...
	gometalinter --disable-all --vendor --enable=gosimple --exclude=mocks/  $(ROOT)/...
	gometalinter --disable-all --vendor --enable=gotype --exclude=mocks/ $(ROOT)/...
	gometalinter --disable-all --vendor --enable=gotypex $(ROOT)/...
	gometalinter --disable-all --vendor --enable=ineffassign $(ROOT)/...
	gometalinter --disable-all --vendor --enable=interfacer  $(ROOT)/...
	gometalinter --disable-all --vendor --enable=megacheck  $(ROOT)/...
	gometalinter --disable-all --vendor --enable=misspell $(ROOT)/...
	gometalinter --disable-all --vendor --enable=nakedret $(ROOT)/...
	gometalinter --disable-all --vendor --enable=safesql --exclude=vendor $(ROOT)/...
	gometalinter --disable-all --vendor --enable=staticcheck $(ROOT)/...
	gometalinter --disable-all --vendor --enable=structcheck $(ROOT)/...
	gometalinter --disable-all --vendor --enable=unconvert  $(ROOT)/...
	gometalinter --disable-all --vendor --enable=varcheck $(ROOT)/...
	gometalinter --disable-all --vendor --enable=vet --exclude=pkg/axmlParser $(ROOT)/...
	gometalinter --disable-all --vendor --enable=vetshadow $(ROOT)/...

format:
	which goimports || go get -u -v golang.org/x/tools/cmd/goimports
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R

#######
# Vendor
#######

check-glide:
	which glide || curl https://glide.sh/get | sh

check-glide-init:
	@[ -f $(ROOT)/glide.yaml ] || make -f $(ROOT)/Makefile glide-init

# Scan a codebase and create a glide.yaml file containing the dependencies.
glide-init: check-glide
	glide init

# Install the latest dependencies into the vendor directory matching the version resolution information.
# The complete dependency tree is installed, importing Glide, Godep, GB, and GPM configuration along the way.
# A lock file is created from the final output.
glide-update: check-glide check-glide-init
	glide update

# Install the dependencies and revisions listed in the lock file into the vendor directory.
# If no lock file exists an update will run.
glide-install: check-glide check-glide-init
	glide install

#########
# Migrate
#########

check-migrate:
	which migrate || go get -u -v github.com/xesina/migrate

migrate-create: check-migrate
	migrate create --name=$(NAME)

migrate-up: check-migrate
	migrate up

migrate-rollback: check-migrate
	migrate rollback

migrate-reset: check-migrate
	migrate reset

migrate-refresh: check-migrate
	migrate refresh

#########
# Test
#########

check-goconvey:
	which goconvey || go get -u -v github.com/smartystreets/goconvey

check-mockery:
	which mockery || go get -u -v github.com/vektra/mockery/.../

mockery-repository: check-mockery
	mockery -name $(ENTITY)Repository -dir lib/$(APPLICATION)/repositories -output common/mocks/$(APPLICATION) -outpkg $(APPLICATION)
	mockery -name RepositoryFactory -dir lib/$(APPLICATION)/repositories -output common/mocks/$(APPLICATION) -outpkg $(APPLICATION)

mockery-service: check-mockery
	mockery -name Service -dir common/services/$(SERVICE) -output common/mocks/services/$(SERVICE)/ -outpkg $(SERVICE)

test:
	ENV=test cd specs && go test
#test: fetch check-goconvey
#	make -f $(ROOT)/Makefile migrate-refresh ENV=testing STEP=0
#	ENV=testing goconvey -host=0.0.0.0 -port=8080 -workDir=$(ROOT) || true

citest: fetch
	make -f $(ROOT)/Makefile migrate-refresh ENV=testing STEP=0
	ENV=testing go list ./... | grep -v /vendor/ | xargs --max-args=1 --replace=R  go test -v -coverprofile=coverage.cov -covermode=atomic R

acceptancetest: fetch
	cd lib && go test ./...

coverage-report:
	[ -d $(COVERAGE_DIR) ] || mkdir -p $(COVERAGE_DIR)
	# Writes atomic mode on top of file
	echo 'mode: atomic' > $(COVERAGE_DIR)/full.cov
	# Collects all coverage files and skips top line with mode
	find $(ROOT)/* -type f -name coverage.cov | xargs tail -q -n +2 >> ${COVERAGE_DIR}/full.cov
	# generate full report
	go tool cover -func=${COVERAGE_DIR}/full.cov
	go tool cover -html=${COVERAGE_DIR}/full.cov -o ${COVERAGE_DIR}/coverage.html
