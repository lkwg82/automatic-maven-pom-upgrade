package de.lgohlke.mavenupgrade;

import lombok.Getter;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.exec.*;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.nio.file.Path;
import java.nio.file.Paths;

@Slf4j
public class Exec {
    private final String command;
    private final Path workingDirectory;

    public Exec(String command) {
        this(command, Paths.get(System.getProperty("user.dir")));
    }

    public Exec(String command, Path workingDirectory) {
        this.command = command;
        this.workingDirectory = workingDirectory;
        log.debug("working directory: {}", workingDirectory.toString());
    }

    public Exec.Result exec(String... args) {
        log.debug("executing: {} {}", command, String.join(" ", args));

        ByteArrayOutputStream out = new ByteArrayOutputStream();
        ByteArrayOutputStream err = new ByteArrayOutputStream();

        CombinedOutputStream out1 = new CombinedOutputStream(out, "Stdout");
        CombinedOutputStream err1 = new CombinedOutputStream(err, "Stderr");

        PumpStreamHandler streamHandler = new PumpStreamHandler(out1, err1);

        DefaultExecutor executor = new DefaultExecutor();
        executor.setStreamHandler(streamHandler);

        CommandLine commandLine = CommandLine.parse(command + " " + String.join(" ", args));
        DefaultExecuteResultHandler executeResultHandler = new DefaultExecuteResultHandler();

        try {
            executor.setWorkingDirectory(workingDirectory.toFile());
            executor.execute(commandLine, executeResultHandler);
            executeResultHandler.waitFor();
        } catch (IOException | InterruptedException e) {
            throw new IllegalStateException(e);
        }

        int exitCode = executeResultHandler.getExitValue();

        boolean commandNotFound = exitCode == Executor.INVALID_EXITVALUE;

        String[] stdoutLines = out.toString().split("\n");
        String[] stderrLines = err.toString().split("\n");
        return new Result(stdoutLines, stderrLines, exitCode, commandNotFound, executeResultHandler.getException());
    }

    @RequiredArgsConstructor
    private static class CombinedOutputStream extends OutputStream {
        private final OutputStream stream;
        private final String prefix;

        private StringBuffer buffer = new StringBuffer();

        @Override
        public void write(int b) throws IOException {
            if (b == '\n') {
                log.debug("  {} {}" ,prefix, buffer.toString().trim());
                buffer = new StringBuffer();
            }
            buffer.append((char)b);
            stream.write(b);
        }
    }

    @RequiredArgsConstructor
    @Getter
    public static class Result {
        private final String[] stdout;
        private final String[] stderr;
        private final int exitCode;
        private final boolean commandNotFound;
        private final Exception exception;
    }
}
