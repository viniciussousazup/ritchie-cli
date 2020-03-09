@Library('ZupSharedLibs@marte') _
node {

  try {

    def projectName = "ritchie-cli"

    buildDockerBuilderAWS {
      dockerRepositoryName = projectName
      dockerFileLocation = "."
      team = "Marte"
      dockerRegistryGroup = "Marte"
      dockerBuilderImage = "golang:1.13-alpine"
    }

    buildWithMakefileAWS {
      dockerRepositoryName = projectName
      dockerFileLocation = "."
      team = "Marte"
      dockerRegistryGroup = "Marte"
      dockerBuildingImage = "${projectName}:builder"
      dockerECRRegion = "sa-east-1"
      notPublishable = true
    }

    syncWithGithubRepo {
      githubDestinationOrg = "martetech"
      githubDestinationRepo = projectName
      githubDestinationBranch = "marte"
    }

  } catch (e) {

      notifyBuildStatus {
        buildStatus = "FAILED"
      }
      throw e

  }

}