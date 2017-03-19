package de.lgohlke.mavenupgrade;

import lombok.Getter;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.exec.*;

import java.io.ByteArrayOutputStream;
import java.io.IOException;

@Slf4j
public class Exec {
    public Result exec(String command, String... args) throws IOException {
        log.debug("executing: {} {}", command, String.join(" ", args));

        ByteArrayOutputStream out = new ByteArrayOutputStream();
        ByteArrayOutputStream err = new ByteArrayOutputStream();
        DefaultExecutor executor = new DefaultExecutor();
        PumpStreamHandler streamHandler = new PumpStreamHandler(out, err);
        executor.setStreamHandler(streamHandler);

        CommandLine commandLine = CommandLine.parse(command + " " + String.join(" ", args));
        DefaultExecuteResultHandler executeResultHandler = new DefaultExecuteResultHandler();

        executor.execute(commandLine, executeResultHandler);

        try {
            executeResultHandler.waitFor();
        } catch (InterruptedException e) {
            throw new IOException(e);
        }
        int exitCode = executeResultHandler.getExitValue();
        ExecuteException exception = executeResultHandler.getException();

        if (null != exception && exception.getCause() instanceof IOException) {
            throw exception;
        }


        String[] stdoutLines = out.toString().split("\n");
        String[] stderrLines = err.toString().split("\n");
        return new Result(stdoutLines, stderrLines, exitCode);
    }

    @RequiredArgsConstructor
    @Getter
    public static class Result {
        private final String[] stdout;
        private final String[] stderr;
        private final int exitCode;
    }
}
