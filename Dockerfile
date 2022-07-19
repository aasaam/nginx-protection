FROM node:16-bullseye AS static

ADD . /src

RUN cd /src/web \
  && echo "npm install" \
  && npm install \
  && npm run lint \
  && npm run build \
  && cd ../static \
  && echo "front build ended" \
  && rm -rf .gitignore

FROM golang:1.17-buster AS builder

ADD . /src

COPY --from=static /src/static /src/static

RUN cd /src \
  && echo "download go lint" \
  && go get -u -v golang.org/x/lint/golint \
  && echo "go mod tidy" \
  && go mod tidy \
  && go get -u -v \
  && echo "go mod download" \
  && go mod download \
  && echo "run lint" \
  && golint . \
  && echo "prepare to test" \
  && export CI=1 \
  && go test -covermode=count -coverprofile=coverage.out \
  && echo "prepare to coverage" \
  && cat coverage.out | grep -v "main.go" > coverage.txt \
  && TOTAL_COVERAGE_FOR_CI_F=$(go tool cover -func coverage.txt | grep total | grep -Eo '[0-9]+.[0-9]+') \
  && echo "TOTAL_COVERAGE_FOR_CI_F: $TOTAL_COVERAGE_FOR_CI_F" \
  && echo "perpare build" \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o nginx-protection \
  && echo "go build ended"

FROM scratch

COPY --from=builder /src/nginx-protection /usr/bin/nginx-protection

ENTRYPOINT ["/usr/bin/nginx-protection"]
