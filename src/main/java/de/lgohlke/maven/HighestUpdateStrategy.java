package de.lgohlke.maven;

import org.apache.maven.artifact.Artifact;
import org.apache.maven.artifact.factory.DefaultArtifactFactory;
import org.apache.maven.artifact.metadata.ArtifactMetadataRetrievalException;
import org.apache.maven.artifact.versioning.ArtifactVersion;
import org.apache.maven.artifact.versioning.InvalidVersionSpecificationException;
import org.apache.maven.artifact.versioning.VersionRange;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.project.MavenProject;
import org.codehaus.mojo.versions.api.ArtifactVersions;
import org.codehaus.mojo.versions.api.VersionsHelper;

/**
 * User: lars
 */
public class HighestUpdateStrategy implements UpdateStrategy {

  private final DefaultArtifactFactory artifactFactory;
  private final VersionsHelper helper;

  public HighestUpdateStrategy(VersionsHelper helper, DefaultArtifactFactory artifactFactory) {
    this.helper = helper;
    this.artifactFactory = artifactFactory;
  }

  @Override
  public ArtifactVersion findNextVersion(MavenProject project, String currentVersion) throws MojoExecutionException {
    VersionRange versionRange = retrieveVersionRange(currentVersion);
    Artifact artifact = artifactFactory.createDependencyArtifact(project.getGroupId(), project.getArtifactId(), versionRange, "pom", null, null);

    try {
      ArtifactVersions artifactVersions = helper.lookupArtifactVersions(artifact, true);
      return artifactVersions.getNewestVersion(versionRange, false);
    } catch (ArtifactMetadataRetrievalException e) {
      throw new MojoExecutionException(e.getMessage(), e);
    }
  }

  private VersionRange retrieveVersionRange(String version) throws MojoExecutionException {
    try {
      return VersionRange.createFromVersionSpec(version);
    } catch (InvalidVersionSpecificationException e) {
      throw new MojoExecutionException(e.getMessage(), e);
    }
  }
}
