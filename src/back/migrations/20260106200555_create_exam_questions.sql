-- +goose Up
CREATE TABLE IF NOT EXISTS exam_questions (
                                              id INT AUTO_INCREMENT PRIMARY KEY,
                                              lecture_id INT NOT NULL,
                                              question_type VARCHAR(50) NOT NULL,
                                              question_text TEXT NOT NULL,
                                              correct_answer TEXT,
                                              options JSON,
                                              points INT DEFAULT 1,
                                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                              FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE,
                                              INDEX idx_lecture_id (lecture_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS exam_questions;