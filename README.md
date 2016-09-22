[![travisci](https://travis-ci.org/lkwg82/automatic-maven-pom-upgrade.svg)](https://travis-ci.org/lkwg82/automatic-maven-pom-upgrade)
[![GoDoc](https://godoc.org/github.com/lkwg82/automatic-maven-pom-upgrade?status.svg)](https://godoc.org/github.com/lkwg82/automatic-maven-pom-upgrade)
 [![GoReportcard](https://goreportcard.com/badge/github.com/lkwg82/automatic-maven-pom-upgrade)](https://goreportcard.com/report/github.com/lkwg82/automatic-maven-pom-upgrade)

automatic-maven-pom-upgrade
===========================

bot which upgrades maven dependencies

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
        automatic upgrade maven projects, 0.3
Options:
                    --git-no-update         skip automerge updates from master
                    --git-no-dirty-check    skip dirty check
                    --git-no-commit         skip commit
                    --hook-after=/bin/echo  command to call after commit (commit message is 1st arg)
                    --maven-settings=       path to maven settings (equivalent to -s)
  -q                --quiet                 suppress any output
                    --version               show version
  -t [help|parent]  --type=[help|parent]    type of upgrade
  -v                --verbose               output verbosely
  -h                --help                  Show usage message

   
# go to test project
$ cd test-projects/simple-parent-update/
$ upgrade --type=parent -v
2016-09-22 22:20:12.820 DEBUG [main.go:128 running in verbose mode]
2016-09-22 22:20:12.820 DEBUG [main.go:129  type: parent]
2016-09-22 22:20:12.820 DEBUG [config.go:30 looking for config file '.autoupgrade.yml']
2016-09-22 22:20:12.820 DEBUG [config.go:33 looking for config file '.autoupgrade.yaml']
2016-09-22 22:20:12.820 DEBUG [config.go:35 no config file found]
2016-09-22 22:20:12.820 DEBUG [exec.go:34 call from: (*Git).CheckIsRepo-fm -> (*Git).CheckIsRepo -> (*Git).IsRepo ]
2016-09-22 22:20:12.820 DEBUG [exec.go:35 executing: git status]
2016-09-22 22:20:12.827 INFO [maven.go:38 determine command]
2016-09-22 22:20:12.827 INFO [maven.go:41 maven wrapper script found]
2016-09-22 22:20:12.827 DEBUG [exec.go:34 call from: (*Maven).DetermineCommand-fm -> (*Maven).DetermineCommand ]
2016-09-22 22:20:12.827 DEBUG [exec.go:35 executing: ./mvnw --version]
2016-09-22 22:20:13.108 DEBUG [exec.go:104    stdout: Apache Maven 3.3.3 (7994120775791599e205a5524ec3e0dfe41d4a06; 2015-04-22T13:57:37+02:00)]
2016-09-22 22:20:13.108 DEBUG [exec.go:104    stdout: Maven home: /Users/lars.gohlke/.m2/wrapper/dists/apache-maven-3.3.3-bin/3opbjp6rgl6qp7k2a6tljcpvgp/apache-maven-3.3.3]
2016-09-22 22:20:13.108 DEBUG [exec.go:104    stdout: Java version: 1.8.0_73, vendor: Oracle Corporation]
2016-09-22 22:20:13.108 DEBUG [exec.go:104    stdout: Java home: /Library/Java/JavaVirtualMachines/jdk1.8.0_73.jdk/Contents/Home/jre]
2016-09-22 22:20:13.108 DEBUG [exec.go:104    stdout: Default locale: de_DE, platform encoding: UTF-8]
2016-09-22 22:20:13.108 DEBUG [exec.go:104    stdout: OS name: "mac os x", version: "10.11.6", arch: "x86_64", family: "mac"]
2016-09-22 22:20:13.111 DEBUG [maven.go:64 use maven settings path: ]
2016-09-22 22:20:13.111 DEBUG [maven.go:66 ignoring empty settings path]
2016-09-22 22:20:13.112 DEBUG [exec.go:34 call from: (*Git).BranchExists ]
2016-09-22 22:20:13.112 DEBUG [exec.go:35 executing: git branch --list --all *autoupdate_parent]
2016-09-22 22:20:13.123 DEBUG [git.go:75 isCurrentBranch:false, isLocalBranch:false, isRemoteBranch: false]
2016-09-22 22:20:13.123 INFO [maven.go:92 updating parent]
2016-09-22 22:20:13.123 DEBUG [exec.go:34 call from: (*Maven).UpdateParent ]
2016-09-22 22:20:13.123 DEBUG [exec.go:35 executing: ./mvnw org.codehaus.mojo:versions-maven-plugin:2.3:update-parent -DgenerateBackupPoms=false --batch-mode]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] Scanning for projects...]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO]                                                                         ]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] ------------------------------------------------------------------------]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] Building test 1.3.7.RELEASE]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] ------------------------------------------------------------------------]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] ]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] --- versions-maven-plugin:2.3:update-parent (default-cli) @ test ---]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] Updating parent from 1.3.7.RELEASE to 1.4.1.RELEASE]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] ------------------------------------------------------------------------]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] BUILD SUCCESS]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] ------------------------------------------------------------------------]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] Total time: 1.081 s]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] Finished at: 2016-09-22T22:20:15+02:00]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] Final Memory: 14M/309M]
2016-09-22 22:20:15.232 DEBUG [exec.go:104    stdout: [INFO] ------------------------------------------------------------------------]
2016-09-22 22:20:15.232 INFO [maven.go:116 updated: Updating parent from 1.3.7.RELEASE to 1.4.1.RELEASE]
result: Updating parent from 1.3.7.RELEASE to 1.4.1.RELEASE
2016-09-22 22:20:15.233 DEBUG [exec.go:34 call from: (*Git).BranchExists ]
2016-09-22 22:20:15.233 DEBUG [exec.go:35 executing: git branch --list --all *autoupdate_parent]
2016-09-22 22:20:15.240 DEBUG [git.go:75 isCurrentBranch:false, isLocalBranch:false, isRemoteBranch: false]
2016-09-22 22:20:15.240 DEBUG [exec.go:34 call from: (*Git).BranchCurrent ]
2016-09-22 22:20:15.240 DEBUG [exec.go:35 executing: git symbolic-ref --short HEAD]
2016-09-22 22:20:15.244 DEBUG [exec.go:104    stdout: yaml-configuration]
2016-09-22 22:20:15.244 DEBUG [exec.go:34 call from: (*Git).BranchCheckoutNew ]
2016-09-22 22:20:15.244 DEBUG [exec.go:35 executing: git checkout -b autoupdate_parent]
2016-09-22 22:20:15.294 DEBUG [exec.go:104    stderr: Switched to a new branch 'autoupdate_parent']
2016-09-22 22:20:15.295 DEBUG [exec.go:104    stdout: M main.go]
2016-09-22 22:20:15.295 DEBUG [exec.go:104    stdout: M test-projects/simple-parent-update/pom.xml]
2016-09-22 22:20:15.295 DEBUG [git.go:201 checking afterCommitHook '/bin/echo']
committing '[Updating parent from 1.3.7.RELEASE to 1.4.1.RELEASE]'
2016-09-22 22:20:15.295 DEBUG [exec.go:34 call from: (*Git).OptionalCommit -> (*Git).Commit ]
2016-09-22 22:20:15.295 DEBUG [exec.go:35 executing: git add pom.xml]
2016-09-22 22:20:15.305 DEBUG [exec.go:34 call from: (*Git).OptionalCommit -> (*Git).Commit ]
2016-09-22 22:20:15.305 DEBUG [exec.go:35 executing: git commit -m 'Updating parent from 1.3.7.RELEASE to 1.4.1.RELEASE' pom.xml]
2016-09-22 22:20:15.331 DEBUG [exec.go:104    stdout: [autoupdate_parent bf4388d] 'Updating parent from 1.3.7.RELEASE to 1.4.1.RELEASE']
2016-09-22 22:20:15.331 DEBUG [exec.go:104    stdout:  1 file changed, 1 insertion(+), 1 deletion(-)]
executing afterCommitHook: [/bin/echo]
2016-09-22 22:20:15.332 DEBUG [exec.go:34 call from: (*Git).OptionalCommit -> (*Git).execAfterCommitHook ]
2016-09-22 22:20:15.332 DEBUG [exec.go:35 executing: /bin/echo Updating parent from 1.3.7.RELEASE to 1.4.1.RELEASE]
2016-09-22 22:20:15.335 DEBUG [exec.go:104    stdout: Updating parent from 1.3.7.RELEASE to 1.4.1.RELEASE]
   
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

# configuration

inspired by `.travis.yml` it should be in a file named `.autoupgrade.yml`

```yaml
# configuration for https://github.com/lkwg82/automatic-maven-pom-upgrade

notifications:
    email: <your-email>
```

# hooks
## afterCommitHook

sets environmental variable `AUTOUPGRADE_NOTIFICATION_EMAIL=<your-email>` before execution when `.autoupgrade.yml` contains
 
```yaml
notifications:
    email: <your-email>
```

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
