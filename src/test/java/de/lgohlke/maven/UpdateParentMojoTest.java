package de.lgohlke.maven;

import org.apache.maven.model.io.DefaultModelWriter;
import org.apache.maven.plugin.testing.AbstractMojoTestCase;
import org.apache.maven.project.MavenProject;
import org.codehaus.plexus.PlexusTestCase;
import org.testng.annotations.BeforeTest;
import org.testng.annotations.Test;

import java.io.ByteArrayOutputStream;
import java.io.File;
import java.util.HashMap;
import java.util.Map;

/**
 * User: lars
 */
public class UpdateParentMojoTest extends AbstractMojoTestCase {

  public static final String GOAL = "update-parent";

//  @BeforeTest
//  public void setUp() throws Exception {
//    super.setUp();
//  }

  /**
   * @throws Exception
   */
  @Test
  public void testMojoGoal() throws Exception {
    File pom = new File(getBasedir() + "/pom.xml");
    System.out.println(pom.getCanonicalPath());
//    configureMojo(new UpdateParentMojo(),"x",pom);
//       setupContainer();
//    UpdateParentMojo mojo = (UpdateParentMojo) lookupMojo(GOAL, pom);
//    assertNotNull( mojo );
    MavenProject parentProject = new MavenProject();
    parentProject.setGroupId("a");
    parentProject.setArtifactId("a");
    parentProject.setVersion("1");
    MavenProject project = new MavenProject();
    project.setParent(parentProject);

    DefaultModelWriter defaultModelWriter = new DefaultModelWriter();
    ByteArrayOutputStream byteOutputStream = new ByteArrayOutputStream();
    Map<String, Object> options = new HashMap<String, Object>();
    defaultModelWriter.write(byteOutputStream, options, project.getModel());
    System.err.println(byteOutputStream.toString());
//    project.writeModel();
//    UpdateParentMojo mojo = (UpdateParentMojo) lookupMojo(GOAL, testPom);
//
//    assertNotNull(mojo);
  }
}
