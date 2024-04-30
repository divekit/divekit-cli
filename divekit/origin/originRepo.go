package origin

/**
 * This file an "object-oriented lookalike" implementation for the structure of the origin repository.
 */

import (
	"divekit-cli/divekit"
	"divekit-cli/divekit/ars"
	"divekit-cli/utils"
	"path/filepath"

	"github.com/apex/log"
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
	utils.OutputAndAbortIfErrors(utils.ValidateAllDirPaths(originRepo.RepoDir))

	originRepo.initDistributions()
	originRepo.ARSConfig.Dir = filepath.Join(originRepo.RepoDir, "ars-config_norepo")
	utils.OutputAndAbortIfErrors(utils.ValidateAllDirPaths(originRepo.ARSConfig.Dir))
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

func (originRepo *OriginRepoType) initIndividualRepositoriesFile(distributionName string, distributionFolder string) error {
	log.Debug("origin.initIndividualRepositoriesFile()")
	individualRepositoriesFilePath, err :=
		utils.FindUniqueFileWithPrefix(distributionFolder, "individual_repositories")
	//utils.OutputAndAbortIfError(err) // TODO: aborts on commands that don't need this file (e.g. setup) - should only happen on commands that need it (validation?)
	if err != nil {
		return err
	}
	distribution, ok := originRepo.DistributionMap[distributionName]
	if !ok {
		// Create a new Distribution if it doesn't exist
		distribution = &Distribution{}
		originRepo.DistributionMap[distributionName] = distribution
	}
	distribution.IndividualizationConfigFileName = individualRepositoriesFilePath
	return nil
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
