FROM golang:1.19.3-bullseye

WORKDIR /app

RUN useradd --create-home gopher && \
  chown -R gopher:gopher /app

USER gopher

COPY --chown=gopher:gopher . .

ENV NVMDIR .nevermind

# Add our hopeful new directory to the PATH, for ease of development
RUN echo "\nexport PATH=\"\$HOME/${NVMDIR}/bin:\$PATH\"" >> ~/.bashrc && \
  # runs go:generate in nvm-shim.go
  # re-run if you change nvm-shim.go
  go generate ./...
