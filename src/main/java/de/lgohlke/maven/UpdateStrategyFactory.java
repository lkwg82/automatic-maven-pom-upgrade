package de.lgohlke.maven;

import org.apache.maven.artifact.factory.DefaultArtifactFactory;
import org.codehaus.mojo.versions.api.VersionsHelper;
import org.codehaus.plexus.component.annotations.Component;

/**
 * User: lars
 */
public class UpdateStrategyFactory {

  private final VersionsHelper helper;
  private final DefaultArtifactFactory artifactFactory;

  public UpdateStrategyFactory(VersionsHelper helper,DefaultArtifactFactory artifactFactory) {
    this.helper = helper;
    this.artifactFactory = artifactFactory;
  }

  public UpdateStrategy createFor(UpdateStrategy.Name strategyName) {
    UpdateStrategy strategy;
    switch (strategyName) {
      case HIGHEST:
        strategy = new HighestUpdateStrategy(helper,artifactFactory);
        break;
      default:
        throw new IllegalArgumentException("could not handle update strategy");
    }
    return strategy;
  }
}
