FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>
RUN apt update -y
RUN apt install -y protobuf-compiler

USER vscode
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN sudo mv migrate /usr/bin/migrate
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install go.uber.org/mock/mockgen@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN mkdir -p ~/.local/bin && wget -qO- https://github.com/ktr0731/evans/releases/download/v0.10.11/evans_linux_amd64.tar.gz | tar xvz -C ~/.local/bin/
# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>
# USER root

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
