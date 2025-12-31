package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ansoncodes/workshot/pkg/types"
)

const (
	defaultDirName = ".workshot"
	shotsSubdir    = "shots"
	indexFile      = "index.json"
)

// metadata stores small snapshot info for fast listing
type Metadata struct {
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	WorkingDir string    `json:"working_dir"`
	GitBranch  string    `json:"git_branch,omitempty"`
}

// index stores all snapshot metadata for fast access
type Index struct {
	Version   int                 `json:"version"`
	Snapshots map[string]Metadata `json:"snapshots"`
}

// storage handles saving and loading snapshots
type Storage struct {
	basePath  string
	indexPath string
}

// new creates and initializes storage
func New() (*Storage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	basePath := filepath.Join(home, defaultDirName, shotsSubdir)
	indexPath := filepath.Join(home, defaultDirName, indexFile)

	// create storage directory if missing
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &Storage{
		basePath:  basePath,
		indexPath: indexPath,
	}, nil
}

// save writes a snapshot to disk and updates index
func (s *Storage) Save(snap *types.Snapshot) error {
	// check snapshot schema version
	if snap.SchemaVersion != types.SchemaVersion {
		return fmt.Errorf("schema version mismatch: got %d, expected %d",
			snap.SchemaVersion, types.SchemaVersion)
	}

	filePath := filepath.Join(s.basePath, snap.Name+".json")

	// convert snapshot to json
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	// write file safely using temp file
	tempPath := filePath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	if err := os.Rename(tempPath, filePath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to finalize snapshot file: %w", err)
	}

	// update index file
	if err := s.updateIndex(snap); err != nil {
		// index can be rebuilt later
		fmt.Fprintf(os.Stderr, "Warning: failed to update index: %v\n", err)
	}

	return nil
}

// load reads a snapshot from disk
func (s *Storage) Load(name string) (*types.Snapshot, error) {
	filePath := filepath.Join(s.basePath, name+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("workshot '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}

	var snap types.Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	// migrate snapshot if version is old
	if snap.SchemaVersion < types.SchemaVersion {
		if err := s.migrateSnapshot(&snap); err != nil {
			return nil, fmt.Errorf("failed to migrate snapshot: %w", err)
		}
	}

	return &snap, nil
}

// list returns all snapshots using index
func (s *Storage) List() ([]Metadata, error) {
	index, err := s.loadIndex()
	if err != nil {
		// rebuild index if missing or broken
		return s.rebuildIndex()
	}

	// convert map to slice
	metadataList := make([]Metadata, 0, len(index.Snapshots))
	for _, meta := range index.Snapshots {
		metadataList = append(metadataList, meta)
	}

	// sort by newest first
	sort.Slice(metadataList, func(i, j int) bool {
		return metadataList[i].CreatedAt.After(metadataList[j].CreatedAt)
	})

	return metadataList, nil
}

// delete removes a snapshot file
func (s *Storage) Delete(name string) error {
	filePath := filepath.Join(s.basePath, name+".json")

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("workshot '%s' not found", name)
		}
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	// remove snapshot from index
	index, _ := s.loadIndex()
	if index != nil {
		delete(index.Snapshots, name)
		s.saveIndex(index)
	}

	return nil
}

// exists checks if snapshot file is present
func (s *Storage) Exists(name string) bool {
	filePath := filepath.Join(s.basePath, name+".json")
	_, err := os.Stat(filePath)
	return err == nil
}

// updateindex adds snapshot metadata to index
func (s *Storage) updateIndex(snap *types.Snapshot) error {
	index, err := s.loadIndex()
	if err != nil {
		index = &Index{
			Version:   1,
			Snapshots: make(map[string]Metadata),
		}
	}

	index.Snapshots[snap.Name] = Metadata{
		Name:       snap.Name,
		CreatedAt:  snap.CreatedAt,
		WorkingDir: snap.WorkingDir,
		GitBranch:  snap.GitBranch,
	}

	return s.saveIndex(index)
}

// loadindex reads index file from disk
func (s *Storage) loadIndex() (*Index, error) {
	data, err := os.ReadFile(s.indexPath)
	if err != nil {
		return nil, err
	}

	var index Index
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	return &index, nil
}

// saveindex writes index file to disk
func (s *Storage) saveIndex(index *Index) error {
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.indexPath, data, 0644)
}

// rebuildindex recreates index from snapshot files
func (s *Storage) rebuildIndex() ([]Metadata, error) {
	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read shots directory: %w", err)
	}

	index := &Index{
		Version:   1,
		Snapshots: make(map[string]Metadata),
	}

	var metadataList []Metadata

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := entry.Name()[:len(entry.Name())-5]

		// load snapshot file
		snap, err := s.Load(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping corrupted snapshot '%s': %v\n", name, err)
			continue
		}

		meta := Metadata{
			Name:       snap.Name,
			CreatedAt:  snap.CreatedAt,
			WorkingDir: snap.WorkingDir,
			GitBranch:  snap.GitBranch,
		}

		index.Snapshots[name] = meta
		metadataList = append(metadataList, meta)
	}

	// save rebuilt index
	if err := s.saveIndex(index); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save rebuilt index: %v\n", err)
	}

	// sort snapshots by time
	sort.Slice(metadataList, func(i, j int) bool {
		return metadataList[i].CreatedAt.After(metadataList[j].CreatedAt)
	})

	return metadataList, nil
}

// migratesnapshot updates snapshot version
func (s *Storage) migrateSnapshot(snap *types.Snapshot) error {
	// future migrations go here
	snap.SchemaVersion = types.SchemaVersion
	return nil
}
