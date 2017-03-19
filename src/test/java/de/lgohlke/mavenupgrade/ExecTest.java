package de.lgohlke.mavenupgrade;

import org.junit.Test;

import java.io.IOException;

import static org.assertj.core.api.Assertions.assertThat;

public class ExecTest {
    @Test
    public void shouldExecuteCommandEcho() throws IOException {
        Exec.Result result = new Exec().exec("echo test");

        assertThat(result.getExitCode()).isEqualTo(0);
        assertThat(result.getStdout()).hasSize(1);
        assertThat(result.getStdout()[0]).isEqualTo("test");
    }

    @Test(expected = IOException.class)
    public void shouldFailExecuteOnUnknownCommand() throws IOException {
        new Exec().exec("xxx test");
    }

    @Test
    public void shouldExecuteCommandReturnsExitCode1() throws IOException {
        Exec.Result result = new Exec().exec("test -f test");

        assertThat(result.getExitCode()).isEqualTo(1);
    }
}
