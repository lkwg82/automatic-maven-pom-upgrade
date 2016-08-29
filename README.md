[![travisci](https://travis-ci.org/lkwg82/automatic-maven-pom-upgrade.svg)](https://travis-ci.org/lkwg82/automatic-maven-pom-upgrade)

automatic-maven-pom-upgrade
===========================

bot which upgrades dependencies (parent-pom)


# idea

Having a bot which automatically try to find new dependencies. If it finds some updates, it will update the `pom.xml` and commit on a branch.

# Roadmap
- DONE parent pom
- TODO dependency updates
- TODO plugin updates

# Requirements
- git installed
- java installed and JAVA_HOME set
- a maven project (of course ;))

# development

## building

with Docker

```bash
./docker/docker_run.sh
```

with installed GO and JDK

```bash
./build.sh
```

and use the `bin/upgrade` binary

## continuous testing

install  `inotify-tools` on Linux (or add OSX support)

```bash
./keep_tests_running.sh
```
