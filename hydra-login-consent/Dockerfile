FROM node:21 AS node-builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

ENV NODE_ENV=production
RUN npm run build

FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /hydra-login-consent .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/templates ./templates
COPY --from=builder --chown=nonroot:nonroot /app/static ./static
COPY --from=node-builder --chown=nonroot:nonroot /app/static/css ./static/css
COPY --from=builder --chown=nonroot:nonroot /hydra-login-consent /usr/bin/hydra-login-consent

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["hydra-login-consent"]
CMD ["hydra-login-consent"]
