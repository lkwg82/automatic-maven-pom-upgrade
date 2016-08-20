FROM golang:1.6

RUN apt-get update && apt-get install -y --no-install-recommends \
		default-jdk && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y --no-install-recommends \
		coreutils && rm -rf /var/lib/apt/lists/*

ENV PROJECT_DIR src/github.com/lkwg82/automatic-maven-pom-upgrade

RUN mkdir -p /go/${PROJECT_DIR}
WORKDIR /go/${PROJECT_DIR}

CMD "./run_in_docker.sh"