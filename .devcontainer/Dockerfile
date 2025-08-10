FROM mcr.microsoft.com/devcontainers/go:1.21

# Install additional development tools
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        curl \
        wget \
        git \
        vim \
        nano \
        htop \
        tree \
        jq \
        yq \
        unzip \
        zip \
        ca-certificates \
        gnupg \
        lsb-release \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install Go tools
RUN go install golang.org/x/tools/cmd/goimports@latest \
    && go install golang.org/x/lint/golint@latest \
    && go install github.com/fatih/gomodifytags@latest \
    && go install github.com/josharian/impl@latest \
    && go install github.com/cweill/gotests/gotests@latest \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install github.com/ramya-rao-a/go-outline@latest \
    && go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest \
    && go install github.com/axw/gocov/gocov@latest

# Set up shell configuration
RUN echo 'export PATH=$PATH:/go/bin' >> /home/vscode/.bashrc \
    && echo 'export PATH=$PATH:/go/bin' >> /home/vscode/.zshrc \
    && echo 'alias ll="ls -la"' >> /home/vscode/.bashrc \
    && echo 'alias ll="ls -la"' >> /home/vscode/.zshrc

# Create a basic go.mod file if it doesn't exist
RUN mkdir -p /workspaces/envvars-cli

# Set proper permissions
RUN chown -R vscode:vscode /go /workspaces

USER vscode

