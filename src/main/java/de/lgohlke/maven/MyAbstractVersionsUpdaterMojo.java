package de.lgohlke.maven;

import org.apache.maven.artifact.factory.ArtifactFactory;
import org.apache.maven.artifact.manager.WagonManager;
import org.apache.maven.artifact.metadata.ArtifactMetadataSource;
import org.apache.maven.artifact.repository.ArtifactRepository;
import org.apache.maven.execution.MavenSession;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.project.MavenProject;
import org.apache.maven.project.MavenProjectBuilder;
import org.apache.maven.project.path.PathTranslator;
import org.apache.maven.settings.Settings;
import org.codehaus.mojo.versions.AbstractVersionsUpdaterMojo;
import org.codehaus.mojo.versions.api.DefaultVersionsHelper;
import org.codehaus.mojo.versions.api.VersionsHelper;
import org.sonatype.aether.impl.ArtifactResolver;

import java.util.List;

/**
 * User: lars
 */
public abstract class MyAbstractVersionsUpdaterMojo extends org.codehaus.mojo.versions.AbstractVersionsUpdaterMojo{
  /**
   * The Maven Project.
   *
   * @parameter expression="${project}"
   * @required
   * @readonly
   * @since 1.0-alpha-1
   */
  private MavenProject project;

  /**
   * @component
   * @since 1.0-alpha-1
   */
  protected ArtifactFactory artifactFactory;

  /**
   * @component
   * @since 1.0-alpha-1
   */
  protected ArtifactResolver resolver;

  /**
   * @component allo
   * @since 1.0-alpha-1
   */
  protected MavenProjectBuilder projectBuilder;

  /**
   * @parameter expression="${reactorProjects}"
   * @required
   * @readonly
   * @since 1.0-alpha-1
   */
  protected List reactorProjects;

  /**
   * The artifact metadata source to use.
   *
   * @component
   * @required
   * @readonly
   * @since 1.0-alpha-1
   */
  protected ArtifactMetadataSource artifactMetadataSource;

  /**
   * @parameter expression="${project.remoteArtifactRepositories}"
   * @readonly
   * @since 1.0-alpha-3
   */
  protected List remoteArtifactRepositories;

  /**
   * @parameter expression="${project.pluginArtifactRepositories}"
   * @readonly
   * @since 1.0-alpha-3
   */
  protected List remotePluginRepositories;

  /**
   * @parameter expression="${localRepository}"
   * @readonly
   * @since 1.0-alpha-1
   */
  protected ArtifactRepository localRepository;

  /**
   * @component
   * @since 1.0-alpha-3
   */
  private WagonManager wagonManager;

  /**
   * @parameter expression="${settings}"
   * @readonly
   * @since 1.0-alpha-3
   */
  protected Settings settings;

  /**
   * settings.xml's server id for the URL.
   * This is used when wagon needs extra authentication information.
   *
   * @parameter expression="${maven.version.rules.serverId}" default-value="serverId";
   * @since 1.0-alpha-3
   */
  private String serverId;

  /**
   * The Wagon URI of a ruleSet file containing the rules that control how to compare version numbers.
   *
   * @parameter expression="${maven.version.rules}"
   * @since 1.0-alpha-3
   */
  private String rulesUri;

  /**
   * Controls whether a backup pom should be created (default is true).
   *
   * @parameter expression="${generateBackupPoms}"
   * @since 1.0-alpha-3
   */
  private Boolean generateBackupPoms;

  /**
   * Whether to allow snapshots when searching for the latest version of an artifact.
   *
   * @parameter expression="${allowSnapshots}" default-value="false"
   * @since 1.0-alpha-1
   */
  protected Boolean allowSnapshots;

  /**
   * Our versions helper.
   */
  private VersionsHelper helper;

  /**
   * The Maven Session.
   *
   * @parameter expression="${session}"
   * @required
   * @readonly
   * @since 1.0-alpha-1
   */
  protected MavenSession session;

  /**
   * @component
   */
  protected PathTranslator pathTranslator;

  /**
   * @component
   */
  protected ArtifactResolver artifactResolver;
}
