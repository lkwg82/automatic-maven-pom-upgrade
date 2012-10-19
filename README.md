automatic-maven-pom-upgrade
===========================

concept and prototypical implementation of a maven plugin which upgrades dependencies step by step


# idea

Create a maven plugin and extends [versions:use-next-releases](http://mojo.codehaus.org/versions-maven-plugin/use-next-releases-mojo.html)
to only upgrade one dependency at the same time. After each single dependency upgrade a test and a compile run as a kind of acceptance test could be run.

# thoughts 

* implement parent-pom updates as well
* implement some rules/hooks
  * break the build when an increased number of compile warnings are detected
  * exclude some artifacts to be upgraded (e.g. spring etc.)
  * some version rules (how to handle version 1-rc2?) (XML/JAVA?)
* (?) risk: update dependencies with excludes 

# update strategies
* deep-update
 * for each found dependency which could be upgraded, try to update to the latest running version
 * try the next possible to be updated dependency
 * iterate until nothing left or no successfull update possible
* wide-update
 * for each found dependency which could be upgraded, try to update to the next running version
 * try the next possible to be updated dependency
 * iterate until nothing left or no successfull update possible
* new other combined strategies (combination/variations)
 

# lifecycle
