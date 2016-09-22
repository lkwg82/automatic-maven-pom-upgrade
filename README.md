[![travisci](https://travis-ci.org/lkwg82/automatic-maven-pom-upgrade.svg)](https://travis-ci.org/lkwg82/automatic-maven-pom-upgrade)
[![GoDoc](https://godoc.org/github.com/lkwg82/automatic-maven-pom-upgrade?status.svg)](https://godoc.org/github.com/lkwg82/automatic-maven-pom-upgrade)
 [![GoReportcard](https://goreportcard.com/badge/github.com/lkwg82/automatic-maven-pom-upgrade)](https://goreportcard.com/report/github.com/lkwg82/automatic-maven-pom-upgrade)

automatic-maven-pom-upgrade
===========================

bot which upgrades dependencies (parent-pom)


# idea

Having a bot which automatically try to find new dependencies. If it finds some updates, it will update the `pom.xml` and commit on a branch.

# scope
**goals**
- assisted updates
- local git integration (branch/commit/merge)

**non-goals**
- fully managed automatic update
- automatic test suite execution to verify update
- git remote interaction (push/pull)

        
# Requirements
- git installed
- java installed and JAVA_HOME set
- a maven project (of course ;))

# how it looks like
```bash
$ git clone https://github.com/lkwg82/automatic-maven-pom-upgrade.git
Klone nach 'automatic-maven-pom-upgrade' ...
...
Prüfe Konnektivität ... Fertig

$ cd automatic-maven-pom-upgrade
$ git checkout tags/v0.1
HEAD ist jetzt bei 72cc87a... Update README.md

# copy the binary upgrade to ~/bin
export PATH=~/bin:$PATH

# check if accessible
$ upgrade
Usage of upgrade:
	automatic upgrade maven projects
Options:
  -v  --verbose               output verbosely
      --hook-after=/bin/echo  command to call after commit (commit message is 1st arg)
      --type=[help|parent]    type of upgrade
  -h  --help                  Show usage message
   
# go to test project
$ cd test-projects/simple-parent-update/
$ upgrade --type=parent -v
2016-08-29 23:18:27.377 DEBUG [exec.go:27 executing: git status]
2016-08-29 23:18:27.378 DEBUG [exec.go:27 executing: git status --porcelain]
2016-08-29 23:18:27.379 INFO [maven.go:31 determine command]
2016-08-29 23:18:27.379 INFO [maven.go:34 maven wrapper script found]
2016-08-29 23:18:27.379 DEBUG [exec.go:27 executing: ./mvnw --version]
2016-08-29 23:18:27.467 DEBUG [exec.go:41 Apache Maven 3.3.3 (7994120775791599e205a5524ec3e0dfe41d4a06; 2015-04-22T13:57:37+02:00)]
2016-08-29 23:18:27.467 DEBUG [exec.go:41 Maven home: /home/lars/.m2/wrapper/dists/apache-maven-3.3.3-bin/3opbjp6rgl6qp7k2a6tljcpvgp/apache-maven-3.3.3]
2016-08-29 23:18:27.467 DEBUG [exec.go:41 Java version: 1.8.0_91, vendor: Oracle Corporation]
2016-08-29 23:18:27.467 DEBUG [exec.go:41 Java home: /usr/lib/jvm/java-8-openjdk-amd64/jre]
2016-08-29 23:18:27.467 DEBUG [exec.go:41 Default locale: de_DE, platform encoding: UTF-8]
2016-08-29 23:18:27.467 DEBUG [exec.go:41 OS name: "linux", version: "4.2.0-42-generic", arch: "amd64", family: "unix"]
2016-08-29 23:18:27.470 DEBUG [exec.go:27 executing: git branch --list autoupdate_parent]
2016-08-29 23:18:27.470 DEBUG [git.go:41 output ]
2016-08-29 23:18:27.470 DEBUG [exec.go:27 executing: git checkout -b autoupdate_parent]
2016-08-29 23:18:27.471 DEBUG [exec.go:41 Zu neuem Branch 'autoupdate_parent' gewechselt]
2016-08-29 23:18:27.472 INFO [maven.go:56 updating parent]
2016-08-29 23:18:27.472 DEBUG [exec.go:27 executing: ./mvnw org.codehaus.mojo:versions-maven-plugin:2.3:update-parent -DgenerateBackupPoms=false --batch-mode]
2016-08-29 23:18:28.638 DEBUG [maven.go:75 [INFO] Scanning for projects...
[INFO]                                                                         
[INFO] ------------------------------------------------------------------------
[INFO] Building test 1.3.7.RELEASE
[INFO] ------------------------------------------------------------------------
[INFO] 
[INFO] --- versions-maven-plugin:2.3:update-parent (default-cli) @ test ---
[INFO] Updating parent from 1.3.7.RELEASE to 1.4.0.RELEASE
[INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS
[INFO] ------------------------------------------------------------------------
[INFO] Total time: 0.604 s
[INFO] Finished at: 2016-08-29T23:18:28+02:00
[INFO] Final Memory: 15M/286M
[INFO] ------------------------------------------------------------------------
]
2016-08-29 23:18:28.638 INFO [maven.go:80 updated: Updating parent from 1.3.7.RELEASE to 1.4.0.RELEASE]
2016-08-29 23:18:28.639 DEBUG [exec.go:27 executing: git add pom.xml]
2016-08-29 23:18:28.640 DEBUG [exec.go:27 executing: git commit -m 'Updating parent from 1.3.7.RELEASE to 1.4.0.RELEASE' pom.xml]
2016-08-29 23:18:28.641 DEBUG [exec.go:41 [autoupdate_parent 392d679] 'Updating parent from 1.3.7.RELEASE to 1.4.0.RELEASE']
2016-08-29 23:18:28.641 DEBUG [exec.go:41  1 file changed, 1 insertion(+), 1 deletion(-)]
2016-08-29 23:18:28.641 DEBUG [exec.go:27 executing: /bin/echo Updating parent from 1.3.7.RELEASE to 1.4.0.RELEASE]
2016-08-29 23:18:28.641 DEBUG [exec.go:41 Updating parent from 1.3.7.RELEASE to 1.4.0.RELEASE]
   
# see the line with  
#   Updating parent from 1.3.7.RELEASE to 1.4.0.RELEASE
   
# see the commit history
$ git log -p -1
commit 392d6792609ea88f8ede0959c89b7094e8d29154
Author: Lars K.W. Gohlke <lkwg82@gmx.de>
Date:   Mon Aug 29 23:18:28 2016 +0200

    'Updating parent from 1.3.7.RELEASE to 1.4.0.RELEASE'

diff --git a/test-projects/simple-parent-update/pom.xml b/test-projects/simple-parent-update/pom.xml
index 47fb9cd..cee402c 100644
--- a/test-projects/simple-parent-update/pom.xml
+++ b/test-projects/simple-parent-update/pom.xml
@@ -5,7 +5,7 @@
     <parent>
         <groupId>org.springframework.boot</groupId>
         <artifactId>spring-boot-starter-parent</artifactId>
-        <version>1.3.7.RELEASE</version>
+        <version>1.4.0.RELEASE</version>
     </parent>
     <artifactId>test</artifactId>
 </project>
\ No newline at end of file
   
```

# Roadmap
- DONE parent pom
- TODO dependency updates
- TODO plugin updates

# functionality
- supports auto merge updates from master into update branch

# development

## requirements
- go 1.7
- `govendor` (install with `go get -u github.com/kardianos/govendor`)
- java and `JAVA_HOME` set

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

## add new packages

```bash
go get ...
# copy them into vendor directory
govendor add +external
```
