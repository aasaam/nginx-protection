FROM node:18-bullseye AS static

ADD . /src

RUN cd /src/web \
  && echo "npm install" \
  && npm install \
  && npm run lint \
  && npm run build \
  && cd ../static \
  && echo "front build ended" \
  && rm -rf .gitignore

FROM golang:1.19-buster AS builder

ADD . /src

COPY --from=static /src/static /src/static

RUN cd /src \
  && go mod tidy \
  && go test -short -covermode=count -coverprofile=coverage.out \
  && cat coverage.out | grep -v "main.go" > coverage.txt \
  && TOTAL_COVERAGE_FOR_CI_F=$(go tool cover -func coverage.txt | grep total | grep -Eo '[0-9]+.[0-9]+') \
  && echo "TOTAL_COVERAGE_FOR_CI_F: $TOTAL_COVERAGE_FOR_CI_F" \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o nginx-protection \
  && ls -lah /src/nginx-protection

FROM alpine

COPY --from=builder /src/nginx-protection /usr/bin/nginx-protection

ENTRYPOINT ["/usr/bin/nginx-protection"]
