package app

import (
	"fmt"
	"github.com/docaura/docaura-cli/internal/fileutils"
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
	"time"
)

// Watcher watches for file changes and triggers documentation regeneration.
type Watcher struct {
	config   Config
	watcher  *fsnotify.Watcher
	debounce time.Duration
}

// NewWatcher creates a new file watcher.
func NewWatcher(config Config) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create fsnotify watcher: %w", err)
	}

	return &Watcher{
		config:   config,
		watcher:  watcher,
		debounce: time.Duration(config.WatchInterval) * time.Second,
	}, nil
}

// Watch starts watching for file changes and calls the regenerate function.
func (w *Watcher) Watch(regenerate func() error) error {
	defer w.watcher.Close()

	// Add directories to watch
	if err := w.addWatchPaths(); err != nil {
		return fmt.Errorf("add watch paths: %w", err)
	}

	if w.config.Verbose {
		log.Printf("Watching %s for changes (checking every %v)...", w.config.ProjectDir, w.debounce)
	}

	// Create a timer for debouncing
	var timer *time.Timer

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			if w.shouldProcessEvent(event) {
				// Reset or create timer for debouncing
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(w.debounce, func() {
					if err := regenerate(); err != nil {
						log.Printf("Regeneration failed: %v", err)
					}
				})
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

// addWatchPaths adds all Go package directories to the watcher.
func (w *Watcher) addWatchPaths() error {
	packages, err := fileutils.FindGoPackages(w.config.ProjectDir, w.config.ExcludeDirs)
	if err != nil {
		return fmt.Errorf("find packages to watch: %w", err)
	}

	for _, packagePath := range packages {
		if err := w.watcher.Add(packagePath); err != nil {
			return fmt.Errorf("add watch path %q: %w", packagePath, err)
		}
	}

	return nil
}

// shouldProcessEvent determines if a file system event should trigger regeneration.
func (w *Watcher) shouldProcessEvent(event fsnotify.Event) bool {
	// Only process write and create events
	if !(event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) {
		return false
	}

	// Only process .go files
	return filepath.Ext(event.Name) == ".go"
}
