-- +goose Up

CREATE TABLE IF NOT EXISTS courses (
                                       id INT AUTO_INCREMENT PRIMARY KEY,
                                       name VARCHAR(255) NOT NULL,
                                       semester VARCHAR(50) NOT NULL,
                                       description TEXT,
                                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                       INDEX idx_semester (semester)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lectures (
                                        id INT AUTO_INCREMENT PRIMARY KEY,
                                        course_id INT NOT NULL,
                                        week INT NOT NULL,
                                        title VARCHAR(255) NOT NULL,
                                        content LONGTEXT NOT NULL,
                                        github_url VARCHAR(500),
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                        FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
                                        INDEX idx_course_week (course_id, week)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS labs (
                                    id INT AUTO_INCREMENT PRIMARY KEY,
                                    course_id INT NOT NULL,
                                    number INT NOT NULL,
                                    title VARCHAR(255) NOT NULL,
                                    description LONGTEXT NOT NULL,
                                    deadline TIMESTAMP NULL,
                                    max_score INT DEFAULT 100,
                                    github_url VARCHAR(500),
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
                                    INDEX idx_course_number (course_id, number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS grade_sheets (
                                            id INT AUTO_INCREMENT PRIMARY KEY,
                                            course_id INT NOT NULL,
                                            sheet_url VARCHAR(500) NOT NULL,
                                            description VARCHAR(255),
                                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                            FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



-- +goose Down


DROP TABLE IF EXISTS grade_sheets;
DROP TABLE IF EXISTS labs;
DROP TABLE IF EXISTS lectures;
DROP TABLE IF EXISTS courses;

