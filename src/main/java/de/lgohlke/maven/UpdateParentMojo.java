package de.lgohlke.maven;

/*
 * Copyright 2001-2005 The Apache Software Foundation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import org.apache.maven.artifact.factory.DefaultArtifactFactory;
import org.apache.maven.artifact.versioning.ArtifactVersion;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.plugin.MojoFailureException;
import org.apache.maven.project.MavenProject;
import org.codehaus.mojo.versions.AbstractVersionsUpdaterMojo;
import org.codehaus.mojo.versions.api.PomHelper;
import org.codehaus.mojo.versions.api.VersionsHelper;
import org.codehaus.mojo.versions.rewriting.ModifiedPomXMLEventReader;

import javax.xml.stream.XMLStreamException;

/**
 * @goal update-parent
 * @requiresProject true
 * @requiresDirectInvocation true
 */
public class UpdateParentMojo extends MyAbstractVersionsUpdaterMojo {


  /**
   * @parameter default-value="HIGHEST"
   */
  private UpdateStrategy.Name updateStrategyName;


  /**
   * @param pom the pom to update.
   * @throws MojoExecutionException when things go wrong
   * @throws MojoFailureException   when things go wrong in a very bad way
   * @throws XMLStreamException     when things go wrong with XML streaming
   * @see AbstractVersionsUpdaterMojo#update(ModifiedPomXMLEventReader)
   * @since 1.0-alpha-1
   */
  protected void update(ModifiedPomXMLEventReader pom) throws MojoExecutionException, MojoFailureException, XMLStreamException {
    if (hasParent() && !parentIsPartOfAReactorProject()) {
      MavenProject parent = getProject().getParent();
      String version = parent.getVersion();

      VersionsHelper helper = getHelper();
      UpdateStrategyFactory updateStrategyFactory = new UpdateStrategyFactory(helper, (DefaultArtifactFactory) artifactFactory);
      UpdateStrategy updateStrategy = updateStrategyFactory.createFor(updateStrategyName);
      ArtifactVersion artifactVersion = updateStrategy.findNextVersion(parent, version);
//      try {
//        artifactVersion = findLatestVersion(artifact, versionRange, null, false);
//      } catch (ArtifactMetadataRetrievalException e) {
//        throw new MojoExecutionException(e.getMessage(), e);
//      }

//      if (shouldApplyUpdate(artifact, version, artifactVersion)) {
//        getLog().debug("Updating parent from " + version + " to " + artifactVersion.toString());
//        setProjectParentVersion(pom, version, artifactVersion);
//      }
    }
  }

  private void setProjectParentVersion(ModifiedPomXMLEventReader pom, String version, ArtifactVersion artifactVersion) throws XMLStreamException {
    if (PomHelper.setProjectParentVersion(pom, artifactVersion.toString())) {
      getLog().debug("Made an update from " + version + " to " + artifactVersion.toString());
    }
  }

  private boolean parentIsPartOfAReactorProject() {
    boolean contains = reactorProjects.contains(getProject().getParent());
    getLog().debug("Project's parent is part of the reactor: " + contains);
    return contains;
  }

  private boolean hasParent() {
    boolean hasParent = getProject().getParent() != null;
    getLog().debug("Project does have a parent : " + hasParent);
    return hasParent;
  }
}
