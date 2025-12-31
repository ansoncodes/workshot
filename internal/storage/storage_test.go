package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ansoncodes/workshot/pkg/types"
)

func TestStorageSaveLoad(t *testing.T) {
	tempDir := t.TempDir()
	
	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", tempDir)
	defer os.Setenv("USERPROFILE", oldHome)

	store, err := New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	snap := types.NewSnapshot("test-snapshot")
	snap.WorkingDir = "C:\\test\\dir"
	snap.GitBranch = "main"
	snap.GitRemote = "https://github.com/test/repo.git"

	if err := store.Save(snap); err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	expectedPath := filepath.Join(tempDir, ".workshot", "shots", "test-snapshot.json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("Snapshot file was not created at %s", expectedPath)
	}

	loaded, err := store.Load("test-snapshot")
	if err != nil {
		t.Fatalf("Failed to load snapshot: %v", err)
	}

	if loaded.Name != snap.Name {
		t.Errorf("Name mismatch: got %s, want %s", loaded.Name, snap.Name)
	}
	if loaded.WorkingDir != snap.WorkingDir {
		t.Errorf("WorkingDir mismatch: got %s, want %s", loaded.WorkingDir, snap.WorkingDir)
	}
	if loaded.GitBranch != snap.GitBranch {
		t.Errorf("GitBranch mismatch: got %s, want %s", loaded.GitBranch, snap.GitBranch)
	}
}

func TestStorageList(t *testing.T) {
	tempDir := t.TempDir()
	
	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", tempDir)
	defer os.Setenv("USERPROFILE", oldHome)

	store, err := New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	names := []string{"snap1", "snap2", "snap3"}
	for _, name := range names {
		snap := types.NewSnapshot(name)
		snap.WorkingDir = "C:\\test"
		if err := store.Save(snap); err != nil {
			t.Fatalf("Failed to save snapshot %s: %v", name, err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	metadataList, err := store.List()
	if err != nil {
		t.Fatalf("Failed to list snapshots: %v", err)
	}

	if len(metadataList) != len(names) {
		t.Errorf("Expected %d snapshots, got %d", len(names), len(metadataList))
	}

	for i := 0; i < len(metadataList)-1; i++ {
		if metadataList[i].CreatedAt.Before(metadataList[i+1].CreatedAt) {
			t.Error("Snapshots are not sorted by creation time (newest first)")
			break
		}
	}
}

func TestStorageDelete(t *testing.T) {
	tempDir := t.TempDir()
	
	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", tempDir)
	defer os.Setenv("USERPROFILE", oldHome)

	store, err := New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	snap := types.NewSnapshot("to-delete")
	snap.WorkingDir = "C:\\test"
	if err := store.Save(snap); err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	if !store.Exists("to-delete") {
		t.Fatal("Snapshot should exist")
	}

	if err := store.Delete("to-delete"); err != nil {
		t.Fatalf("Failed to delete snapshot: %v", err)
	}

	if store.Exists("to-delete") {
		t.Fatal("Snapshot should not exist after deletion")
	}
}

func TestStorageSchemaVersion(t *testing.T) {
	tempDir := t.TempDir()
	
	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", tempDir)
	defer os.Setenv("USERPROFILE", oldHome)

	store, err := New()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	snap := types.NewSnapshot("version-test")
	snap.WorkingDir = "C:\\test"

	if err := store.Save(snap); err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	loaded, err := store.Load("version-test")
	if err != nil {
		t.Fatalf("Failed to load snapshot: %v", err)
	}

	if loaded.SchemaVersion != types.SchemaVersion {
		t.Errorf("Schema version mismatch: got %d, want %d",
			loaded.SchemaVersion, types.SchemaVersion)
	}
}