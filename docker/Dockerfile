FROM golang:1.7

RUN apt-get update && apt-get install -y --no-install-recommends \
		default-jdk \
		coreutils \
		upx \
		&& rm -rf /var/lib/apt/lists/*

CMD "docker/run_in_docker.sh"
