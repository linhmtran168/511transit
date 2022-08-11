# Go
FROM golang:1.19-alpine as go_builder

RUN apk update \ 
  && apk upgrade \
  && apk add --no-cache ca-certificates \
  && update-ca-certificates

RUN mkdir -p /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . . 

RUN mkdir -p bin
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags '-extldflags "-static"' -o bin/transit cmd/transit.go

# JS, CSS
FROM node:lts-alpine as asset_builder

COPY web /workspace/web
WORKDIR /workspace/web

RUN npm ci
RUN npm run build

# Final image
FROM scratch

# CA certificate
COPY --from=go_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go_builder /app/bin/ /app
COPY --from=go_builder /app/.env /app
COPY --from=asset_builder /workspace/web/dist /app/web/dist

WORKDIR /app

CMD [ "./transit" ]