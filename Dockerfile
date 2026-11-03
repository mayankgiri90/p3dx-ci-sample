# Sample service image: static Go binary on distroless, running as the nonroot user.
#   docker build --build-arg VERSION=1.0.0 --build-arg GIT_SHA=$(git rev-parse HEAD) -t p3dx-ci-sample .

FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
ARG VERSION=dev
ARG GIT_SHA=none
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
      -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${GIT_SHA}" \
      -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/server-missing"]
