package origin

/**
 * This file an "object-oriented lookalike" implementation for the structure of the origin repository.
 */

import (
	"divekit-cli/divekit"
	"divekit-cli/divekit/ars"
	"divekit-cli/utils/errorHandling"
	"divekit-cli/utils/fileUtils"
	"github.com/apex/log"
	"path/filepath"
)

// global variable for the origin repository
var (
	OriginRepo *OriginRepoType
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
func NewOriginRepo(originRepoName string) *OriginRepoType {
	log.Debug("origin.InitOriginRepoPaths()")
	originRepo := &OriginRepoType{}
	originRepo.RepoDir = filepath.Join(divekit.DivekitHomeDir, originRepoName)
	errorHandling.OutputAndAbortIfErrors(fileUtils.ValidateAllDirPaths(originRepo.RepoDir),
		"The path to the originRepo is ")

	originRepo.initDistributions()
	originRepo.ARSConfig.Dir = filepath.Join(originRepo.RepoDir, "ars-config_norepo")
	errorHandling.OutputAndAbortIfErrors(fileUtils.ValidateAllDirPaths(originRepo.ARSConfig.Dir),
		"The path to the ars config dir is invalid")
	return originRepo
}

func InitOriginRepo(originRepoNameFlag string) {
	if originRepoNameFlag != "" {
		OriginRepo = NewOriginRepo(originRepoNameFlag)
	}
}

func (originRepo *OriginRepoType) GetDistribution(distributionName string) *Distribution {
	log.Debug("origin.GetDistribution()")
	return originRepo.DistributionMap[distributionName]
}

func (originRepo *OriginRepoType) initDistributions() {
	log.Debug("origin.initDistributions()")
	distributionRootDir := filepath.Join(originRepo.RepoDir, ".divekit_norepo/distributions")
	originRepo.DistributionMap = make(map[string]*Distribution)
	distributionFolders, err := fileUtils.ListSubfolderNames(distributionRootDir)
	errorHandling.OutputAndAbortIfError(err, "The path to the distribution root dir is invalid")

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
		fileUtils.FindUniqueFileWithPrefix(distributionFolder, "individual_repositories")
	errorHandling.OutputAndAbortIfError(err, "The path to the individual_repositories file is invalid")
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
	// filename for NewRepositoryConfigFile is fixed, must be "repositoryConfig.json"
	repositoryConfigFile := ars.NewRepositoryConfigFile(filepath.Join(distributionFolder, "repositoryConfig.json"))
	distribution, ok := originRepo.DistributionMap[distributionName]
	if !ok {
		// Create a new Distribution if it doesn't exist
		distribution = &Distribution{}
		originRepo.DistributionMap[distributionName] = distribution
	}
	distribution.RepositoryConfigFile = repositoryConfigFile
}
