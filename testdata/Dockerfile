FROM golang:1.24

# Install dpkg/apt tools
RUN apt-get update && apt-get install -y \
    dpkg \
    apt \
    && rm -rf /var/lib/apt/lists/*

# Set working directory inside container
WORKDIR /app

# Copy everything (mod files + code)
COPY . .

# Download modules
RUN go mod download

# Run tests in apt package
CMD ["go", "test", "-v", "./"]
