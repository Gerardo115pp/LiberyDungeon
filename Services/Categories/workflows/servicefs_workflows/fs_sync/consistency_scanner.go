package fs_sync

import (
	"fmt"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

type unregisteredFile struct {
	FilePath     string
	CategoryUUID string
}

type unregisteredCategory struct {
	DirectoryPath string
	ParentUUID    string
}

type stateSyncErrors struct {
	ScanRootPath                string
	StateMap                    map[string][]dungeon_models.MediaWeakIdentity
	UnregisteredFiles           []unregisteredFile     // Supported files found in the fs but not in the db
	UnregisteredCategoriesPaths []unregisteredCategory // Directories found in the fs but that do not represent a category in the db
	GhostIdentities             []dungeon_models.MediaWeakIdentity
	SeenPaths                   map[string]struct{} // Used to keep track of the paths already scanned to find ghost medias and categories. its of type map[string]struct{} because struct{} is of size 0.
}

func (sync_errors *stateSyncErrors) addUnregisteredFile(unregistered_file_path string) error {
	var parent_path string = dungeon_helpers.GetParentDirectory(unregistered_file_path)
	parent_path = dungeon_helpers.NormalizePath(parent_path)

	parent_medias, exists := sync_errors.StateMap[parent_path]
	if !exists {
		return fmt.Errorf("Parent directory not found in the state map")
	}

	if len(parent_medias) == 0 {
		return fmt.Errorf("Parent directory has no registered categories")
	}

	var parent_uuid string = parent_medias[0].CategoryUUID
	var new_unregistered_file unregisteredFile = unregisteredFile{
		FilePath:     unregistered_file_path,
		CategoryUUID: parent_uuid,
	}

	sync_errors.UnregisteredFiles = append(sync_errors.UnregisteredFiles, new_unregistered_file)

	echo.EchoDebug(fmt.Sprintf("Unregistered file: '%s'", unregistered_file_path))

	return nil
}

// Adds a category to the list of unregistered categories. The parent of the given path must be a registered category.
// It will retrieve the parent uuid from the media weak identity list mapped to the parent directory of the given path.
// that list must have at least one element.
func (sync_errors *stateSyncErrors) addUnregisteredCategory(unregistered_path string) error {
	var parent_path string = dungeon_helpers.GetParentDirectory(unregistered_path)
	parent_path = dungeon_helpers.NormalizePath(parent_path)

	parent_medias, exists := sync_errors.StateMap[parent_path]
	if !exists {
		return fmt.Errorf("Parent directory not found in the state map")
	}

	if len(parent_medias) == 0 {
		return fmt.Errorf("Parent directory has no registered categories")
	}

	var parent_uuid string = parent_medias[0].CategoryUUID
	var new_unregistered_category unregisteredCategory = unregisteredCategory{
		DirectoryPath: unregistered_path,
		ParentUUID:    parent_uuid,
	}

	sync_errors.UnregisteredCategoriesPaths = append(sync_errors.UnregisteredCategoriesPaths, new_unregistered_category)

	echo.EchoDebug(fmt.Sprintf("Unregistered category: '%s'", unregistered_path))

	return nil

}

func (sync_errors *stateSyncErrors) fileInState(fs_path string) bool {
	var file_exists bool = false
	directory_path := dungeon_helpers.GetParentDirectory(fs_path)
	directory_path = dungeon_helpers.NormalizePath(directory_path)
	filename := filepath.Base(fs_path)

	media_identities, exists := sync_errors.StateMap[directory_path]
	if !exists {
		return file_exists
	}

	for _, media_identity := range media_identities {
		if filename == media_identity.MediaName {
			file_exists = true
			break
		}
	}

	echo.EchoDebug(fmt.Sprintf("Scanning file '%s' in directory '%s': %t", filename, directory_path, file_exists))

	return file_exists
}

func (sync_errors *stateSyncErrors) directoryInState(fs_path string) bool {
	var directory_exists bool = false
	directory_path := dungeon_helpers.NormalizePath(fs_path)

	_, exists := sync_errors.StateMap[directory_path]
	if exists {
		directory_exists = true
	}

	echo.EchoDebug(fmt.Sprintf("Scanning directory '%s': %t", directory_path, directory_exists))

	return directory_exists
}

func (sync_errors stateSyncErrors) String() string {
	var str_format string = "SYNC ERRORS\n"

	str_format += fmt.Sprintf("Scan root path: '%s'\n", sync_errors.ScanRootPath) + strings.Repeat("-", 80) + "\n"

	if len(sync_errors.StateMap) > 0 {
		str_format += "State map:\n"
		for key, value := range sync_errors.StateMap {
			identity_format := fmt.Sprintf("\t%s:", key)
			for _, media_identity := range value {
				item_label := media_identity.MediaName
				if item_label == "" {
					item_label = media_identity.CategoryPath
				}

				identity_format += fmt.Sprintf("%s->%s '%s'", echo.OrangeFG, echo.WhiteFG, item_label)
			}
			str_format += identity_format + "\n\n"
		}
		str_format += "\n\n" + strings.Repeat("-", 80) + "\n"

	} else {
		str_format += "NO STATE MAP\n" + strings.Repeat("-", 80) + "\n"
	}

	if len(sync_errors.UnregisteredFiles) > 0 {
		str_format += "Unregistered files:\n"
		for _, unregistered_file := range sync_errors.UnregisteredFiles {
			str_format += fmt.Sprintf("\t'%s' should be on %s\n", unregistered_file.FilePath, unregistered_file.CategoryUUID)
		}
		str_format += "\n" + strings.Repeat("-", 80) + "\n"
	} else {
		str_format += "NO UNREGISTERED FILES\n" + strings.Repeat("-", 80) + "\n"
	}

	if len(sync_errors.UnregisteredCategoriesPaths) > 0 {
		str_format += "Unregistered categories:\n"
		for _, unregistered_category := range sync_errors.UnregisteredCategoriesPaths {
			str_format += fmt.Sprintf("\t'%s' should be on %s\n", unregistered_category.DirectoryPath, unregistered_category.ParentUUID)
		}
		str_format += "\n" + strings.Repeat("-", 80) + "\n"
	} else {
		str_format += "NO UNREGISTERED CATEGORIES\n" + strings.Repeat("-", 80) + "\n"
	}

	if len(sync_errors.SeenPaths) > 0 {
		str_format += "Seen paths:\n"
		for seen_path := range sync_errors.SeenPaths {
			str_format += fmt.Sprintf("\t%s\n", seen_path)
		}
		str_format += "\n" + strings.Repeat("-", 80) + "\n"
	} else {
		str_format += "NO SEEN PATHS, THIS IS VERY WEIRD\n" + strings.Repeat("-", 80) + "\n"
	}

	if len(sync_errors.GhostIdentities) > 0 {
		str_format += "Ghost identities:\n"
		for _, ghost_identity := range sync_errors.GhostIdentities {
			str_format += fmt.Sprintf("\t%s\n", ghost_identity.String())
		}
		str_format += "\n" + strings.Repeat("-", 80) + "\n"
	} else {
		str_format += "NO GHOST IDENTITIES\n" + strings.Repeat("-", 80) + "\n"
	}

	return str_format
}

func (sync_errors *stateSyncErrors) reportGhostFiles(category_identity *dungeon_models.CategoryIdentity, all_media_identities []dungeon_models.MediaWeakIdentity) {
	var base_path string = category_identity.ClusterPath
	echo.EchoDebug(fmt.Sprintf("Getting ghost identities for base path: %s", base_path))

	for _, media_identity := range all_media_identities {
		media_identity_path := filepath.Join(base_path, media_identity.CategoryPath)
		media_identity_path = dungeon_helpers.NormalizePath(media_identity_path)
		if media_identity.MediaName != "" {
			media_identity_path = filepath.Join(media_identity_path, media_identity.MediaName)
		}

		if _, exists := sync_errors.SeenPaths[media_identity_path]; !exists {
			echo.EchoDebug(fmt.Sprintf("File '%s' not seen", media_identity_path))
			sync_errors.GhostIdentities = append(sync_errors.GhostIdentities, media_identity)
		}
	}
}

// Callback function for filepath.WalkDir. It will scan the given path for consistency with the database state.
func (sync_errors *stateSyncErrors) scanClusterPath(path string, entry os.DirEntry, err error) error {
	echo.EchoDebug(fmt.Sprintf("%s\nScanning path: %s", strings.Repeat("-", 80), path))
	is_directory := entry.IsDir()

	if is_directory {
		std_path := dungeon_helpers.NormalizePath(path)
		sync_errors.SeenPaths[std_path] = struct{}{}

		if !sync_errors.directoryInState(std_path) {
			err := sync_errors.addUnregisteredCategory(std_path)
			if err != nil {
				return fmt.Errorf("Error adding unregistered category: %s", err.Error())
			}

			return filepath.SkipDir
		}
	} else {
		if !dungeon_helpers.IsSupportedFileExtension(path) {
			echo.EchoDebug(fmt.Sprintf("Skipping unsupported file: '%s'", path))
			return nil
		}

		sync_errors.SeenPaths[path] = struct{}{}

		if !sync_errors.fileInState(path) {
			err := sync_errors.addUnregisteredFile(path)
			if err != nil {
				return fmt.Errorf("Error adding unregistered file: %s", err.Error())
			}
		}
	}

	return nil
}

func newStateSyncErrors(state_map map[string][]dungeon_models.MediaWeakIdentity, scan_root_path string) *stateSyncErrors {
	var sync_errors *stateSyncErrors = new(stateSyncErrors)

	sync_errors.UnregisteredFiles = make([]unregisteredFile, 0)
	sync_errors.UnregisteredCategoriesPaths = make([]unregisteredCategory, 0)
	sync_errors.GhostIdentities = make([]dungeon_models.MediaWeakIdentity, 0)
	sync_errors.StateMap = state_map
	sync_errors.ScanRootPath = scan_root_path
	sync_errors.SeenPaths = make(map[string]struct{})

	return sync_errors
}

func scanSyncErrors(root_path string, db_state_map map[string][]dungeon_models.MediaWeakIdentity) (*stateSyncErrors, *dungeon_models.LabeledError) {
	echo.EchoDebug(fmt.Sprintf("%s\nScanning cluster path: %s", strings.Repeat("-", 80), root_path))
	var sync_errors *stateSyncErrors = newStateSyncErrors(db_state_map, root_path)

	err := filepath.WalkDir(root_path, sync_errors.scanClusterPath)
	if err != nil {
		fmt.Println(sync_errors.String())
		return nil, dungeon_models.NewLabeledError(err, "in scanSyncErrors, while calling filepath.WalkDir", dungeon_models.ErrProcessError)
	}

	echo.EchoDebug("Finished scanning cluster path")
	return sync_errors, nil
}

// Returns a look up table of fs_paths -> MediaWeakIdentity. Used to quickly check if a file or directory is already registered in the db.
// It will also add the category directories to the branch_content list so it's consistent with the lookup table.
func buildFsPathLookupTable(branch_content *[]dungeon_models.MediaWeakIdentity, branch_path string) map[string][]dungeon_models.MediaWeakIdentity {
	echo.EchoDebug(fmt.Sprintf("Building lookup table for branch: %s", branch_path))
	var lookup_table map[string][]dungeon_models.MediaWeakIdentity = make(map[string][]dungeon_models.MediaWeakIdentity)
	var extra_media_identities []dungeon_models.MediaWeakIdentity = make([]dungeon_models.MediaWeakIdentity, 0)

	for _, media_identity := range *branch_content {
		fs_path := filepath.Join(branch_path, media_identity.CategoryPath)
		fs_path = dungeon_helpers.NormalizePath(fs_path)
		echo.EchoDebug(fmt.Sprintf("%s -> '%s'", media_identity.CategoryPath, fs_path))

		if _, exists := lookup_table[fs_path]; !exists {
			echo.EchoDebug(fmt.Sprintf("Adding fs_path: %s", fs_path))
			lookup_table[fs_path] = make([]dungeon_models.MediaWeakIdentity, 0)
			if media_identity.MediaName != "" {
				// media_identity with an empty MediaName means it's a category directory.
				// but the database will no return a category record if it has media in it. so we need to add the category directory to the lookup table.
				category_identity := dungeon_models.MediaWeakIdentity{
					CategoryPath: media_identity.CategoryPath,
					MediaName:    "",
					MediaUUID:    "",
					CategoryUUID: media_identity.CategoryUUID,
				}

				lookup_table[fs_path] = append(lookup_table[fs_path], category_identity)
				extra_media_identities = append(extra_media_identities, category_identity)
			}
		}

		lookup_table[fs_path] = append(lookup_table[fs_path], media_identity)
	}

	*branch_content = append(*branch_content, extra_media_identities...)

	echo.EchoDebug("Finished building lookup table")
	return lookup_table
}
