FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    clang \
    llvm \
    iproute2 \
    linux-tools-common \
    linux-tools-generic \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY test_prog.c .
COPY gobpftool .

# Compile the eBPF program
RUN clang -O2 -g -target bpf -c test_prog.c -o test_prog.o || true

CMD ["/bin/bash"]
