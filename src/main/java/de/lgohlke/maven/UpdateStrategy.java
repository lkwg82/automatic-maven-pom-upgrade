package de.lgohlke.maven;

import org.apache.maven.artifact.metadata.ArtifactMetadataRetrievalException;
import org.apache.maven.artifact.versioning.ArtifactVersion;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.project.MavenProject;

/**
 * User: lars
 */
public interface UpdateStrategy {
  enum Name {
    HIGHEST(HighestUpdateStrategy.class),
    NEXT(HighestUpdateStrategy.class),
    BINARY(HighestUpdateStrategy.class);

    private final Class<? extends UpdateStrategy> updateStrategyClass;

    private Name(Class<? extends UpdateStrategy> updateStrategyClass) {
      this.updateStrategyClass = updateStrategyClass;
    }

    public Class<? extends UpdateStrategy> getUpdateStrategyClass() {
      return updateStrategyClass;
    }
  }

  ArtifactVersion findNextVersion(MavenProject project, String currentVersion) throws MojoExecutionException;
}
