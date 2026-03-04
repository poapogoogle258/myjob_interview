FROM golang:1.25.4-alpine
WORKDIR /app
COPY . .
RUN rm -f .env
RUN go mod download
RUN go build -o main ./cmd/app/
EXPOSE 8077
# Create an empty .env file because main.go panics if it's missing
RUN touch .env
CMD ["./main"]