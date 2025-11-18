# Start from the Golang base image to build the Go application
FROM golang:1.25 AS build
WORKDIR /app
# Set the target OS and architecture for cross-compilation
ENV GOOS=linux
ENV GOARCH=arm64

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the necessary source directories
COPY lambda/ lambda/
COPY bin/ bin/

# Build the Go application
RUN go build -o main ./lambda/scrape/main.go
RUN go build -o pw-install ./bin/pw-install/main.go

# Start from the Amazon Linux 2023 base image for deployment (make sure this base supports ARM64)
FROM public.ecr.aws/lambda/provided:al2023

# Install Node.js using NodeSource and DNF (confirm Node.js and npm support ARM64 in this environment)
RUN curl -sL https://rpm.nodesource.com/setup_20.x | bash - && \
    dnf install -y nodejs

# Install Playwright and its dependencies
RUN dnf install -y \
    nss \
    atk \
    cups-libs \
    libXcomposite \
    libXrandr \
    libxkbcommon \
    libXScrnSaver \
    pango \
    at-spi2-atk \
    at-spi2-core \
    libXdamage \
    libXfixes \
    mesa-libgbm \
    alsa-lib-devel

# Copy the built Go binary from the build stage
COPY --from=build /app/main ./main
COPY --from=build /app/pwinstall ./pwinstall

RUN chmod +x ./pw-install
RUN ./pw-install

# Copy the Playwright installation cache
RUN cp -r /root/.cache ./.cache

# Adjust permissions to allow execution
RUN chmod -R a+rx ./.cache/ms-playwright-go/

# Set a home directory environment variable
ENV HOME=.

# Set the entrypoint to the Go binary
ENTRYPOINT ["./main"]