package de.lgohlke.mavenupgrade;

import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.concurrent.TimeUnit;

import static org.assertj.core.api.Assertions.assertThat;

public class GitTest {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder();

    private Git git;
    private Path path;
    private Exec gitExec;

    @Before
    public void setUp2() throws Exception {
        path = temporaryFolder.newFolder().toPath();
        git = new Git(path);
        gitExec = new Exec("git", path);
        assertThat(git.isInstalled()).isTrue();
    }

    @Test
    public void shouldDetectMissingRepository() throws Exception {
        assertThat(git.isRepo()).isFalse();
    }

    @Test
    public void shouldDetectRepositoryWhenInSubdirectory() throws Exception {
        createRepoWithSingleCommit();

        Path subdir = path.resolve(Paths.get("x", "s", "a"));
        subdir.toFile().mkdirs();

        Git git = new Git(subdir);
        assertThat(git.isRepo()).isTrue();
    }

    @Test
    public void IsNotDirtyGitRepository() throws Exception {
        execGit("init");

        assertThat(git.isDirty()).isFalse();
    }

    @Test
    public void IsDirtyGitRepository() throws Exception {
        execGit("init");

        path.resolve("test").toFile().createNewFile();

        assertThat(git.isDirty()).isTrue();
    }


    @Test
    public void BranchExists() throws IOException {
        createRepoWithSingleCommit();

        execGit("checkout", "-b", "test");

        assertThat(git.BranchExists("test")).isTrue();
    }

    @Test
    public void BranchExistsNotOnBranch() throws IOException {
        createRepoWithSingleCommit();

        execGit("checkout", "-b", "test");
        execGit("checkout", "master");

        assertThat(git.BranchExists("test")).isTrue();
    }

    private void execGit(String... args) {
        Exec.Result result = gitExec.exec(args);
        if (result.hasError()) {
            throw new IllegalStateException(result.getStderrLines());
        }
    }

    private void createRepoWithSingleCommit() throws IOException {
        execGit("init");

        execGit("config", "user.email", "test@ci.com");
        execGit("config", "user.name", "test");

        path.resolve("test").toFile().createNewFile();

        execGit("add", "test");
        execGit("commit", "-m", "'test'", "test");
    }

    @Test
    public void BranchExistsRemote() throws Exception {
        // remote repository
        Path remote = path.resolve("remote");
        remote.toFile().mkdir();

        gitExec = new Exec("git", remote);
        execGit("init");
        execGit("config", "user.email", "test@ci.com");
        execGit("config", "user.name", "test");
        remote.resolve("test").toFile().createNewFile();
        execGit("add", "test");
        execGit("commit", "-m", "'test'", "test");
        execGit("checkout", "-b", "test");
        execGit("checkout", "master");

        // local
        gitExec = new Exec("git", path);
        execGit("clone", remote.toString(), "local");

        Git gitExec = new Git(path.resolve("local"));
        assertThat(gitExec.BranchExists("test")).isTrue();
    }

    @Test
    public void BranchCheckoutNew() throws Exception {
        createRepoWithSingleCommit();

        git.BranchCheckoutNew("test");

        assertThat(git.BranchExists("test")).as("missing branch test").isTrue();
        assertThat(git.BranchCurrent()).isEqualTo("test");
    }

    @Test
    public void BranchCheckoutExisting() throws Exception {
        createRepoWithSingleCommit();
        execGit("checkout", "-b", "test");
        execGit("checkout", "master");

        git.BranchCheckoutExisting("test");

        assertThat(git.BranchExists("test")).as("missing branch test").isTrue();
        assertThat(git.BranchCurrent()).isEqualTo("test");
    }

    @Test
    public void Commit() throws Exception {
        createRepoWithSingleCommit();

        path.resolve("pom.xml").toFile().createNewFile();

        git.Commit("initial update");

        assertThat(git.isDirty()).isFalse();
    }

    @Test
    public void IsInSyncWithMasterSame() throws Exception {
        createRepoWithSingleCommit();

        execGit("checkout", "-b", "test");

        assertThat(git.IsInSyncWith("master")).isTrue();
    }

    @Test
    public void IsInSyncWithMasterAhead() throws Exception {
        createRepoWithSingleCommit();

        execGit("checkout", "-b", "test");
        path.resolve("test2").toFile().createNewFile();
        execGit("add", "test2");
        execGit("commit", "-m", "'test'", "test2");

        assertThat(git.IsInSyncWith("master")).isTrue();
    }

    @Test
    public void IsNotInSyncWithMaster() throws Exception {
        createRepoWithSingleCommit();

        path.resolve("test2").toFile().createNewFile();
        execGit("add", "test2");
        execGit("commit", "-m", "'test'", "test2");

        execGit("checkout", "-b", "test");
        execGit("reset", "--hard", "HEAD~1");

        assertThat(git.IsInSyncWith("master")).isFalse();
        assertThat(git.isDirty()).isFalse();
    }

    @Test
    public void MergeMasterIntoBranch() throws Exception {
        createRepoWithSingleCommit();

        path.resolve("test2").toFile().createNewFile();
        execGit("add", "test2");
        execGit("commit", "-m", "'test'", "test2");

        execGit("checkout", "-b", "test");
        execGit("reset", "--hard", "HEAD~1");

        Files.write(path.resolve("test2"), "test".getBytes());
        execGit("add", "test2");
        execGit("commit", "-m", "'test in branch test2'","test2");

        assertThat(git.IsInSyncWith("master")).isFalse();
        assertThat(git.HasMergeConflict("master")).isFalse();

        // action
        git.DominantMergeFrom("master", "updates from master");

        assertThat(git.IsInSyncWith("master")).isTrue();
    }

    @Test
    public void MergeMasterIntoBranchWithConflict() throws Exception {
        createRepoWithSingleCommit();

        // write to test in branch test
        execGit("checkout", "-b", "test");

        Files.write(path.resolve("test"), "test".getBytes());
        execGit("add", "test");
        execGit("commit", "-m", "'update test'", "test");

        // write to test in master different
        execGit("checkout", "master");

        Files.write(path.resolve("test"), "afafd".getBytes());
        execGit("add", "test");
        execGit("commit", "-m", "'update test'", "test");

        execGit("checkout", "test");

        assertThat(git.HasMergeConflict("master")).isTrue();

        // action
        git.DominantMergeFrom("master", "updates from master");

        TimeUnit.MILLISECONDS.sleep(10);

        assertThat(git.IsInSyncWith("master")).isTrue();
        assertThat("afafd".getBytes()).isEqualTo(Files.readAllBytes(path.resolve("test")));
    }
}
