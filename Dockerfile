FROM golang:1.22 AS base

WORKDIR /opt/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o ./app

FROM base AS production

ENV PORT "8080"
ENV REDIS_URL "redis://redis:6379/0"
COPY --from=base /opt/app/app /usr/bin/beats
CMD ["beats"]

