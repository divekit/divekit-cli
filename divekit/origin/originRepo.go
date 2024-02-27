package origin

import (
	"divekit-cli/divekit"
	"divekit-cli/divekit/ars"
	"divekit-cli/utils/fileUtils"
	"fmt"
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

func NewOriginRepo(originRepoName string) (*OriginRepoType, error) {
	log.Debug("origin.InitOriginRepoPaths()")
	originRepo := &OriginRepoType{}
	originRepo.RepoDir = filepath.Join(divekit.DivekitHomeDir, originRepoName)
	if err := fileUtils.ValidateAllDirPaths(originRepo.RepoDir); err != nil {
		return nil, err
	}

	if err := originRepo.initDistributions(); err != nil {
		return nil, err
	}

	originRepo.ARSConfig.Dir = filepath.Join(originRepo.RepoDir, "ars-config_norepo")

	if err := fileUtils.ValidateAllDirPaths(originRepo.ARSConfig.Dir); err != nil {
		return nil, err
	}

	return originRepo, nil
}

func InitOriginRepo(originRepoNameFlag string) error {
	if originRepoNameFlag == "" {
		return &OriginRepoError{"The origin repo name flag is not defined"}
	}

	var err error
	if OriginRepo, err = NewOriginRepo(originRepoNameFlag); err != nil {
		return err
	}

	return nil
}

func (originRepo *OriginRepoType) GetDistribution(distributionName string) (*Distribution, error) {
	log.Debug("origin.GetDistribution()")
	if distribution := originRepo.DistributionMap[distributionName]; distribution != nil {
		return distribution, nil
	}

	return nil, &OriginRepoError{fmt.Sprintf("The distribution '%s' does not exist", distributionName)}
}

func (originRepo *OriginRepoType) initDistributions() error {
	log.Debug("origin.initDistributions()")
	distributionRootDir := filepath.Join(originRepo.RepoDir, ".divekit_norepo/distributions")
	originRepo.DistributionMap = make(map[string]*Distribution)
	distributionFolders, err := fileUtils.ListSubFolderNames(distributionRootDir)
	if err != nil {
		return fmt.Errorf("The path to the distribution root dir is invalid: %w", err)
	}

	for _, distributionName := range distributionFolders {
		distributionFolder := filepath.Join(distributionRootDir, distributionName)
		newDistribution := Distribution{
			Dir: distributionFolder,
		}
		originRepo.DistributionMap[distributionName] = &newDistribution
		if err := originRepo.initIndividualRepositoriesFile(distributionName, distributionFolder); err != nil {
			return err
		}
		if err := originRepo.initRepositoryConfigFile(distributionName, distributionFolder); err != nil {
			return err
		}
	}

	return nil
}

func (originRepo *OriginRepoType) initIndividualRepositoriesFile(distributionName string, distributionFolder string) error {
	log.Debug("origin.initIndividualRepositoriesFile()")
	individualRepositoriesFilePath, err := fileUtils.FindUniqueFileWithPrefix(distributionFolder, "individual_repositories")
	if err != nil {
		return fmt.Errorf("The path to the individual_repositories file is invalid: %w", err)
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

func (originRepo *OriginRepoType) initRepositoryConfigFile(distributionName string, distributionFolder string) error {
	log.Debug("origin.initRepositorConfigFile()")
	// filename for NewRepositoryConfigFile is fixed, must be "repositoryConfig.json"
	repositoryConfigFile, err := ars.NewRepositoryConfigFile(filepath.Join(distributionFolder, "repositoryConfig.json"))
	if err != nil {
		return err
	}

	distribution, ok := originRepo.DistributionMap[distributionName]
	if !ok {
		// Create a new Distribution if it doesn't exist
		distribution = &Distribution{}
		originRepo.DistributionMap[distributionName] = distribution
	}
	distribution.RepositoryConfigFile = repositoryConfigFile

	return nil
}

type OriginRepoError struct {
	Msg string
}

func (e *OriginRepoError) Error() string {
	return e.Msg
}
