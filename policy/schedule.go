package policy

import (
	"log"

	"github.com/chidi150c/database/gorm"
	"github.com/chidi150c/database/model"
	"github.com/robfig/cron/v3"
)

// Define the maximum number of recent records to retain
const maxRecentRecords = 500  // Adjust this value as needed

// Function to enforce the retention policy
func enforceRetentionPolicy(dbs *gorm.DBServices) error {
    // Count the total number of records in the database
    var totalRecords int
    if err := dbs.DB.Model(&model.TradingSystem{}).Count(&totalRecords).Error; err != nil {
        return err
    }

    // Calculate the number of excess records to remove
    excessRecords := totalRecords - maxRecentRecords

    // If there are excess records, delete or archive them
    if excessRecords > 0 {
        log.Printf("\nRetention policy exceeded by %d records.\n", excessRecords)
        log.Printf("\nRetention policy exceeded by %d records, Total records: %d, MaxLimitrecords: %d.\n", excessRecords, totalRecords, maxRecentRecords)
        // Modify this query to select the excess records based on your criteria (e.g., timestamp)
        excessRecordsQuery := dbs.DB.Order("timestamp_column ASC").Limit(excessRecords)

        // Delete or archive the excess records
        if err := excessRecordsQuery.Delete(&model.TradingSystem{}).Error; err != nil {
            return err
        }
        log.Printf("\nRetention policy excess deleted successfully!!!.\n")
    }

    return nil
}

// Scheduled task to enforce the retention policy at regular intervals
func ScheduleRetentionTask(dbs *gorm.DBServices) {
    // Create a new cron scheduler
    c := cron.New()

    // Define the schedule for running the retention policy...
    _, err := c.AddFunc("@midnight", func() {
        if err := enforceRetentionPolicy(dbs); err != nil {
            log.Printf("Retention policy enforcement error: %v", err)
            // Optionally, you can send alerts or take specific actions on error
        } else {
            log.Printf("Retention policy enforcement successful.")
        }
    })

    if err != nil {
        log.Printf("Failed to add retention policy task: %v", err)
        return
    }

    // Start the cron scheduler
    c.Start()

    // Optionally, stop the scheduler when your application exits
    defer c.Stop()

    // Keep the application running (you may have other code here)
    select {}
}
