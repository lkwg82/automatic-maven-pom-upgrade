FROM golang:1.6

RUN apt-get update && apt-get install -y --no-install-recommends \
		default-jdk \
		coreutils \
		&& rm -rf /var/lib/apt/lists/*

CMD "./run_in_docker.sh"