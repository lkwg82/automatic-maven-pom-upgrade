package de.lgohlke.mavenupgrade;

import lombok.extern.slf4j.Slf4j;

import java.nio.file.Path;
import java.nio.file.Paths;

@Slf4j
public class Git extends Exec {
    public Git(Path workingDirectory) {
        super("git", workingDirectory);
    }

    public Git() {
        this(Paths.get(System.getProperty("user.dir")));
    }

    public boolean isInstalled() {
        return !exec("--version").hasError();
    }

    public boolean isRepo() {
        return !exec("status").hasError();
    }

    public boolean isDirty() {
        String content = execExitOnError("status", "--porcelain").getStdoutLines();

        if ("".equals(content)) {
            return false;
        }
        return !content.isEmpty();
    }

    public boolean BranchExists(String branch) {
        Result result = execExitOnError("branch", "--list", "--all", "*" + branch);

        String[] lines = result.getStdout();
        boolean isCurrentBranch = lines[0].equals("* " + branch);
        boolean isLocalBranch = lines[0].equals("  " + branch);
        boolean isRemoteBranch = lines[0].equals("  remotes/origin/" + branch);
        log.debug("isCurrentBranch: {}, isLocalBranch: {}, isRemoteBranch: {}", isCurrentBranch, isLocalBranch, isRemoteBranch);
        return isCurrentBranch || isLocalBranch || isRemoteBranch;
    }

    public void BranchCheckoutNew(String branch) {
        execExitOnError("checkout", "-b", branch);
    }

    public String BranchCurrent() {
        Result result = execExitOnError("symbolic-ref", "--short", "HEAD");

        return result.getStdout()[0].replace("refs/heads/", "");
    }

    public void BranchCheckoutExisting(String branch) {
        execExitOnError("checkout", branch);
    }

    public void Commit(String message) {
        execExitOnError("add", "pom.xml");
        execExitOnError("commit", "-m", "'" + message + "'", "pom.xml");
    }

    public boolean IsInSyncWith(String branch) {
        Result result = execExitOnError("merge", "--no-commit","--no-ff", branch);

        String output = result.getStdoutLines().trim();
        if ("Already up-to-date.".equals(output)){
            return true;
        }

        execExitOnError("merge","--abort");
        return false;
    }

    public boolean HasMergeConflict(String branch) {
        Result result = exec("merge", "--no-commit", branch);

        String content = result.combinedOutput().trim();
        if (content.equals("Already up-to-date.")) {
            return false;
        } else {
            execExitOnError("merge", "--abort");
            return result.hasError();
        }
    }

    public void DominantMergeFrom(String branch, String message) {
        execExitOnError("merge", "--commit", "--strategy-option=theirs", branch, "--message", "\""+message+"\"");
    }
}
