package de.lgohlke.mavenupgrade;

import org.junit.Test;

import java.nio.file.Paths;

import static org.assertj.core.api.Assertions.assertThat;

public class ExecTest {

    @Test
    public void shouldExecuteCommandEcho() {
        Exec.Result result = new Exec("echo").exec("test");

        assertThat(result.getExitCode()).isEqualTo(0);
        assertThat(result.getStdout()).hasSize(1);
        assertThat(result.getStdout()[0]).isEqualTo("test");
    }

    @Test
    public void shouldFailExecuteOnUnknownCommand() {
        Exec.Result result = new Exec("xxx").exec("test");

        assertThat(result.isCommandNotFound()).isTrue();
    }

    @Test
    public void shouldExecuteCommandReturnsExitCode1() {
        Exec.Result result = new Exec("test").exec("-f", "test");

        assertThat(result.getExitCode()).isEqualTo(1);
    }

    @Test
    public void shouldExecuteInGivenDirectory() {
        Exec exec = new Exec("pwd", Paths.get("/tmp"));

        Exec.Result result = exec.exec();

        assertThat(result.getStdout()[0]).isEqualTo("/tmp");
    }
}
