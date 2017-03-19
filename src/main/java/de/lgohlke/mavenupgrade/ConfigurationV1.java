package de.lgohlke.mavenupgrade;

import lombok.Getter;
import lombok.Setter;

@Getter
public class ConfigurationV1 {
    private Notification notifications = new Notification();

    @Getter
    @Setter
    public static class Notification {
        private String email;
    }
}
