package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	testFiles := []string{"test0.txt", "test1.txt", "test2.txt", "test3.txt", "test4.txt", "test5.txt"}

	for _, fileName := range testFiles {
		t.Run(fileName, func(t *testing.T) {
			file, err := os.Open(fileName)
			if err != nil {
				t.Fatalf("Failed to open %s: %v", fileName, err)
			}
			defer file.Close()

			chambers, numAnts, start, end := parseFile(file)
			if numAnts <= 0 {
				t.Errorf("Expected positive ants in %s, got %d", fileName, numAnts)
			}
			if start == nil {
				t.Errorf("Expected start chamber in %s, got nil", fileName)
			}
			if end == nil {
				t.Errorf("Expected end chamber in %s, got nil", fileName)
			}
			if len(chambers) == 0 {
				t.Errorf("Expected chambers map in %s, got empty", fileName)
			}
			if start != nil && end != nil && start.ID == end.ID {
				t.Errorf("Start and end chambers should be distinct in %s", fileName)
			}
		})
	}
}

func TestFindVertexDisjointPaths(t *testing.T) {
	testFiles := []string{"test0.txt", "test1.txt", "test2.txt", "test3.txt"}

	for _, fileName := range testFiles {
		t.Run(fileName, func(t *testing.T) {
			file, err := os.Open(fileName)
			if err != nil {
				t.Fatalf("Failed to open %s: %v", fileName, err)
			}
			defer file.Close()

			_, numAnts, start, end := parseFile(file)
			if start == nil || end == nil {
				t.Fatalf("Missing start or end in %s", fileName)
			}

			paths := FindVertexDisjointPaths(start, end, numAnts)
			if len(paths) == 0 {
				t.Fatalf("Expected at least one path in %s, found none", fileName)
			}

			// Verify that each path starts at start and ends at end
			for i, path := range paths {
				if len(path) < 2 {
					t.Errorf("Path %d in %s has length %d, expected >= 2", i, fileName, len(path))
				}
				if path[0].ID != start.ID {
					t.Errorf("Path %d does not start at %s", i, start.ID)
				}
				if path[len(path)-1].ID != end.ID {
					t.Errorf("Path %d does not end at %s", i, end.ID)
				}
			}

			// Verify vertex-disjoint property (no intermediate chambers shared between paths)
			seenIntermediate := make(map[string]int)
			for pathIdx, path := range paths {
				for _, chamber := range path[1 : len(path)-1] {
					if prevPath, exists := seenIntermediate[chamber.ID]; exists {
						t.Errorf("Chamber %s shared between path %d and path %d in %s (violates vertex-disjoint)",
							chamber.ID, prevPath, pathIdx, fileName)
					}
					seenIntermediate[chamber.ID] = pathIdx
				}
			}
		})
	}
}

func TestOptimalAntDistribution(t *testing.T) {
	file, err := os.Open("test1.txt")
	if err != nil {
		t.Fatalf("Failed to open test1.txt: %v", err)
	}
	defer file.Close()

	_, numAnts, start, end := parseFile(file)
	paths := FindVertexDisjointPaths(start, end, numAnts)
	if len(paths) == 0 {
		t.Fatal("Expected paths in test1.txt")
	}

	dist := OptimalAntDistribution(paths, numAnts)
	if len(dist) != len(paths) {
		t.Fatalf("Expected distribution length %d, got %d", len(paths), len(dist))
	}

	totalDistributed := 0
	for i, count := range dist {
		if count < 0 {
			t.Errorf("Negative ant count on path %d: %d", i, count)
		}
		totalDistributed += count
	}

	if totalDistributed != numAnts {
		t.Errorf("Expected total ants %d, distributed %d", numAnts, totalDistributed)
	}
}

func TestSimulateAnts(t *testing.T) {
	file, err := os.Open("test0.txt")
	if err != nil {
		t.Fatalf("Failed to open test0.txt: %v", err)
	}
	defer file.Close()

	_, numAnts, start, end := parseFile(file)
	ants := createAnts(numAnts, start)
	paths := FindVertexDisjointPaths(start, end, numAnts)
	dist := OptimalAntDistribution(paths, numAnts)

	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "sim_test.log")

	err = SimulateAntsOnPaths(ants, paths, dist, logPath)
	if err != nil {
		t.Fatalf("SimulateAntsOnPaths failed: %v", err)
	}

	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Log file was not created: %v", err)
	}
	if info.Size() == 0 {
		t.Errorf("Log file is empty")
	}
}
