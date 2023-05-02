package origin

/**
 * This file an "object-oriented lookalike" implementation for the structure of the origin repository.
 */

import (
	"divekit-cli/divekit"
	"divekit-cli/divekit/ars"
	"divekit-cli/utils"
	"github.com/apex/log"
	"path/filepath"
)

// all the relevant paths in the origin repository (all as full paths)
type OriginRepoType struct {
	RepoDir         string
	DistributionMap map[string]*Distribution
	ARSConfig       struct {
		Dir string
	}
}

type Distribution struct {
	Dir                             string
	RepositoryConfigFile            *ars.RepositoryConfigFileType
	IndividualizationConfigFileName string
}

// This method is similar to a constructor in OOP
func OriginRepo(originRepoName string) *OriginRepoType {
	log.Debug("origin.InitOriginRepoPaths()")
	originRepo := &OriginRepoType{}
	originRepo.RepoDir = filepath.Join(divekit.DivekitHomeDir, originRepoName)
	originRepo.initDistributions()
	originRepo.ARSConfig.Dir = filepath.Join(originRepo.RepoDir, "arsConfig")
	return originRepo
}

func (originRepo *OriginRepoType) GetDistribution(distributionName string) *Distribution {
	log.Debug("origin.GetDistribution()")
	return originRepo.DistributionMap[distributionName]
}

func (originRepo *OriginRepoType) initDistributions() {
	log.Debug("origin.initDistributions()")
	distributionRootDir := filepath.Join(originRepo.RepoDir, ".divekit_norepo\\distributions")
	originRepo.DistributionMap = make(map[string]*Distribution)
	distributionFolders, err := utils.ListSubfolderNames(distributionRootDir)
	utils.OutputAndAbortIfError(err)

	for _, distributionName := range distributionFolders {
		distributionFolder := filepath.Join(distributionRootDir, distributionName)
		newDistribution := Distribution{
			Dir: distributionFolder,
		}
		originRepo.DistributionMap[distributionName] = &newDistribution
		originRepo.initIndividualRepositoriesFile(distributionName, distributionFolder)
		originRepo.initRepositorConfigFile(distributionName, distributionFolder)
	}
}

func (originRepo *OriginRepoType) initIndividualRepositoriesFile(distributionName string, distributionFolder string) {
	log.Debug("origin.initIndividualRepositoriesFile()")
	individualRepositoriesFilePath, err :=
		utils.FindUniqueFileWithPrefix(distributionFolder, "individual_repositories")
	utils.OutputAndAbortIfError(err)
	distribution, ok := originRepo.DistributionMap[distributionName]
	if !ok {
		// Create a new Distribution if it doesn't exist
		distribution = &Distribution{}
		originRepo.DistributionMap[distributionName] = distribution
	}
	distribution.IndividualizationConfigFileName = individualRepositoriesFilePath
}

func (originRepo *OriginRepoType) initRepositorConfigFile(distributionName string, distributionFolder string) {
	log.Debug("origin.initRepositorConfigFile()")
	// filename for RepositoryConfigFile is fixed, must be "repositoryConfig.json"
	repositoryConfigFile := ars.RepositoryConfigFile(filepath.Join(distributionFolder, "repositoryConfig.json"))
	distribution, ok := originRepo.DistributionMap[distributionName]
	if !ok {
		// Create a new Distribution if it doesn't exist
		distribution = &Distribution{}
		originRepo.DistributionMap[distributionName] = distribution
	}
	distribution.RepositoryConfigFile = repositoryConfigFile
}
