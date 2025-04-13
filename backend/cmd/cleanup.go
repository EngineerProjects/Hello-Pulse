package main

import (
    "context"
    "log"
    "time"
    
    fileservice "hello-pulse.fr/internal/services/file"
)

// StartFileCleanupTask starts a background goroutine that periodically cleans up expired files
func StartFileCleanupTask(fileService *fileservice.Service) {
    go func() {
        // Run cleanup every 24 hours
        ticker := time.NewTicker(24 * time.Hour)
        defer ticker.Stop()
        
        log.Println("File cleanup task started")
        
        for range ticker.C {
            // Calculate threshold date (30 days ago)
            threshold := time.Now().AddDate(0, 0, -30)
            
            log.Printf("Running file cleanup for files deleted before %s", threshold.Format(time.RFC3339))
            
            // Run cleanup
            if err := fileService.CleanupExpiredFiles(context.Background(), threshold); err != nil {
                log.Printf("Error during file cleanup: %v", err)
            } else {
                log.Println("File cleanup completed successfully")
            }
        }
    }()
}