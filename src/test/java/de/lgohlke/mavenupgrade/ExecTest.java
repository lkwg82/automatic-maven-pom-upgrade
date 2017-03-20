package de.lgohlke.mavenupgrade;

import org.junit.Test;

import java.nio.file.Paths;

import static org.assertj.core.api.Assertions.assertThat;

public class ExecTest {
    private final Exec exec = new Exec();

    @Test
    public void shouldExecuteCommandEcho() {
        Exec.Result result = exec.exec("echo test");

        assertThat(result.getExitCode()).isEqualTo(0);
        assertThat(result.getStdout()).hasSize(1);
        assertThat(result.getStdout()[0]).isEqualTo("test");
    }

    @Test
    public void shouldFailExecuteOnUnknownCommand() {
        Exec.Result result = exec.exec("xxx test");

        assertThat(result.isCommandNotFound()).isTrue();
    }

    @Test
    public void shouldExecuteCommandReturnsExitCode1() {
        Exec.Result result = exec.exec("test -f test");

        assertThat(result.getExitCode()).isEqualTo(1);
    }

    @Test
    public void shouldExecuteInGivenDirectory() {
        Exec exec = new Exec(Paths.get("/tmp"));

        Exec.Result result = exec.exec("pwd");

        assertThat(result.getStdout()[0]).isEqualTo("/tmp");
    }
}
