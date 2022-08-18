# Compile runtime binary inside alpine-based builder-container
# (avoid using local go-binary to prevent "works on my machine" problems)
FROM golang:1.18-alpine AS build

WORKDIR /app
ADD . /app

ENV VENDORED=true
RUN apk --update --no-cache add make
RUN make build

# Build Runtime Image
# (Ref: https://github.com/GoogleContainerTools/distroless)
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=build /app/go-app /go-app
USER 65532:65532

ENTRYPOINT ["/go-app"]
