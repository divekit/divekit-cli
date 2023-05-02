package config

/**
 * This file an "object-oriented lookalike" implementation for the structure of the origin repository.
 */

import (
	"divekit-cli/cmd"
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
	Dir                         string
	RepositoryConfigFile        *RepositoryConfigFileType
	IndividualizationConfigFile string
}

// This method is similar to a constructor in OOP
func OriginRepo(originRepoName string) *OriginRepoType {
	log.Debug("config.InitOriginRepoPaths()")
	originRepo := &OriginRepoType{}
	originRepo.RepoDir = filepath.Join(DivekitHomeDir, cmd.OriginRepoNameFlag)
	originRepo.initDistributions()
	originRepo.ARSConfig.Dir = filepath.Join(originRepo.RepoDir, "arsConfig")
	return originRepo
}

func (originRepo *OriginRepoType) GetDistribution(distributionName string) *Distribution {
	return originRepo.DistributionMap[distributionName]
}

func (originRepo *OriginRepoType) initDistributions() {
	distributionRootDir := filepath.Join(originRepo.RepoDir, ".divekit_norepo\\distribution")
	distributionFolders, err := utils.ListSubfolders(distributionRootDir)
	utils.OutputAndAbortIfError(err)

	for _, distributionName := range distributionFolders {
		distributionFolder := filepath.Join(distributionRootDir, distributionName)
		originRepo.DistributionMap[distributionName] = &Distribution{
			Dir: distributionFolder,
		}
		originRepo.initIndividualRepositoriesFile(distributionName, distributionFolder)
		originRepo.initRepositorConfigFile(distributionName, distributionFolder)
	}
}

func (originRepo *OriginRepoType) initIndividualRepositoriesFile(distributionName string, distributionFolder string) {
	individualRepositoriesFilePath, err :=
		utils.FindUniqueFileWithPrefix(distributionFolder, "individual_repositories")
	utils.OutputAndAbortIfError(err)
	distribution, ok := originRepo.DistributionMap[distributionName]
	if !ok {
		// Create a new Distribution if it doesn't exist
		distribution = &Distribution{}
		originRepo.DistributionMap[distributionName] = distribution
	}
	distribution.IndividualizationConfigFile = individualRepositoriesFilePath
}

func (originRepo *OriginRepoType) initRepositorConfigFile(distributionName string, distributionFolder string) {
	// filename for RepositoryConfigFile is fixed, must be "repositoryConfig.json"
	repositoryConfigFile := RepositoryConfigFile(filepath.Join(distributionFolder, "repositoryConfig.json"))
	distribution, ok := originRepo.DistributionMap[distributionName]
	if !ok {
		// Create a new Distribution if it doesn't exist
		distribution = &Distribution{}
		originRepo.DistributionMap[distributionName] = distribution
	}
	distribution.RepositoryConfigFile = repositoryConfigFile
}
