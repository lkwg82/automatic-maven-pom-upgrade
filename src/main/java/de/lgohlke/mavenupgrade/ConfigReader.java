package de.lgohlke.mavenupgrade;

import com.esotericsoftware.yamlbeans.YamlReader;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;

@Slf4j
@RequiredArgsConstructor
class ConfigReader {
    private final static String FILENAME1 = ".autoupgrade.yml";
    private final static String FILENAME2 = ".autoupgrade.yaml";

    // needs to be testable, when working with temporary directories
    private final Path currentWorkindDirctory;

    public ConfigReader() {
        this(Paths.get(System.getProperty("user.dir")));
    }

    @SuppressWarnings("unchecked")
    public ConfigurationV1 readConfigV1() throws IOException {
        Map map = createRawConfigMap();

        log.debug("reading config version 1");
        ConfigurationV1 configurationV1 = new ConfigurationV1();

        Map<String, Object> notifications = (Map) map.getOrDefault("notifications", new HashMap<>());
        String email = (String) notifications.getOrDefault("email", "");
        configurationV1.getNotifications().setEmail(email);

        return configurationV1;
    }

    private Map createRawConfigMap() throws FileNotFoundException, com.esotericsoftware.yamlbeans.YamlException {
        FileReader reader = createReader();

        // no config file found
        if (null == reader) {
            return new HashMap();
        }

        YamlReader yamlReader = new YamlReader(reader);
        Map map = yamlReader.read(Map.class);

        // in case of empty file, avoid NPEs on the way
        if (map == null) {
            map = new HashMap();
        }
        return map;
    }

    private FileReader createReader() throws FileNotFoundException {
        log.debug("looking for config file {} or {}", FILENAME1, FILENAME2);
        File file1 = currentWorkindDirctory.resolve(FILENAME1).toFile();
        File file2 = currentWorkindDirctory.resolve(FILENAME2).toFile();

        FileReader fileReader = null;
        if (file1.exists() && file2.exists()) {
            log.debug("  found both config files");
            throw new IllegalStateException("could not have both: " + file1 + " and " + file2);
        }

        if (file1.exists()) {
            log.debug(" found {}",FILENAME1);
            fileReader = new FileReader(file1);
        } else {
            if (file2.exists()) {
                log.debug(" found {}",FILENAME2);
                fileReader = new FileReader(file2);
            }else{
                log.debug(" found no config");
            }
        }
        return fileReader;
    }
}
