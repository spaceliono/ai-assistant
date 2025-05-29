package database

import (
	"together-ai-assistant/models"
)

// 插入新记录
func InsertAIAssistantRecord(record *models.AIAssistantRecords) error {
	query := `INSERT INTO AI_Assistant_Records (source_platform, source_id, ai_platform, role, content, create_time) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, record.SourcePlatform, record.SourceID, record.AIPlatform, record.Role, record.Content, record.CreateTime)
	return err
}

// 查询记录
func GetAIAssistantRecordsBySource(sourcePlatform, sourceID string) ([]models.AIAssistantRecords, error) {
	rows, err := db.Query("SELECT id, source_platform, source_id, ai_platform, role, content, create_time FROM AI_Assistant_Records WHERE source_platform = ? AND source_id = ? order by id", sourcePlatform, sourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var AIrecords []models.AIAssistantRecords
	for rows.Next() {
		var record models.AIAssistantRecords
		if err := rows.Scan(&record.ID, &record.SourcePlatform, &record.SourceID, &record.AIPlatform, &record.Role, &record.Content, &record.CreateTime); err != nil {
			return nil, err
		}
		AIrecords = append(AIrecords, record)
	}
	return AIrecords, nil
}

func GetJobs() ([]string, error) {
	rows, err := db.Query("SELECT name FROM eb_user_tag")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []string
	for rows.Next() {
		var job string
		if err := rows.Scan(&job); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
