FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /hydra-cleanup .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /hydra-cleanup /usr/bin/hydra-cleanup

USER nonroot:nonroot

ENTRYPOINT ["hydra-cleanup"]
CMD ["cron", "-f"]
