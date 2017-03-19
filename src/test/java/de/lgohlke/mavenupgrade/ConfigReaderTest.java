package de.lgohlke.mavenupgrade;

import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

import static org.assertj.core.api.Assertions.assertThat;

public class ConfigReaderTest {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder();

    private ConfigReader configReader;
    private Path path;

    @Before
    public void setUp() throws Exception {
        path = temporaryFolder.newFolder().toPath();
        configReader = new ConfigReader(path);
    }

    @Test(expected = IllegalStateException.class)
    public void failOnBothYamlAndYmlConfig() throws IOException {
        writeYaml(".autoupgrade.yaml", "");
        writeYaml(".autoupgrade.yml", "");

        configReader.readConfigV1();
    }

    @Test
    public void doNotFailOnMissingConfig() throws IOException {
        configReader.readConfigV1();
    }

    @Test
    public void shouldReadYaml() throws IOException {
        writeYaml(".autoupgrade.yaml", "");

        configReader.readConfigV1();
    }

    @Test
    public void shouldReadYml() throws IOException {
        writeYaml(".autoupgrade.yml", "");

        configReader.readConfigV1();
    }

    @Test
    public void shouldReadNotificationEmail() throws IOException {
        String ymlStr = "" +
                "notifications:\n" +
                "   email: me@email.de";
        writeYaml(".autoupgrade.yml", ymlStr);

        ConfigurationV1 configurationV1 = configReader.readConfigV1();

        assertThat(configurationV1.getNotifications().getEmail()).isEqualTo("me@email.de");
    }

    private void writeYaml(String file, String content) throws IOException {
        Files.write(path.resolve(file), content.getBytes());
    }

}
