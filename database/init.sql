-- ground-turn 本地 MySQL 初始化脚本
-- 表结构由后端 GORM AutoMigrate 维护；此脚本保证容器首次启动时库表存在、字符集正确。
CREATE DATABASE IF NOT EXISTS app_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE app_db;

CREATE TABLE IF NOT EXISTS flight_turnaround (
  id                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  flight_no          VARCHAR(32)  NOT NULL,
  aircraft_reg       VARCHAR(32),
  stand_no           VARCHAR(16),
  arrival_time       DATETIME(3),
  departure_time     DATETIME(3),
  turnaround_status  VARCHAR(32)  NOT NULL,
  delay_reason       VARCHAR(255),
  INDEX idx_flight_status (turnaround_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ground_task (
  id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  turnaround_id  BIGINT UNSIGNED NOT NULL,
  task_type      VARCHAR(32) NOT NULL,
  team_id        VARCHAR(32),
  planned_start  DATETIME(3) NULL,
  deadline       DATETIME(3) NOT NULL COMMENT '原计划截止时间，登记延误永不改写',
  accepted_at    DATETIME(3) NULL COMMENT '签收时间',
  actual_finish  DATETIME(3) NULL COMMENT '实际完成时间',
  status         VARCHAR(32) NOT NULL,
  blocker_note   VARCHAR(255),
  INDEX idx_task_turnaround (turnaround_id),
  INDEX idx_task_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ground_resource (
  id                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  resource_code       VARCHAR(32) NOT NULL,
  resource_type       VARCHAR(32) NOT NULL,
  location            VARCHAR(64),
  availability_status VARCHAR(32) NOT NULL,
  maintenance_due_at  DATETIME(3) NULL,
  owner_team          VARCHAR(32),
  INDEX idx_resource_status (availability_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS resource_booking (
  id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  resource_id     BIGINT UNSIGNED NOT NULL,
  turnaround_id   BIGINT UNSIGNED NOT NULL,
  task_id         BIGINT UNSIGNED NULL,
  start_time      DATETIME(3) NOT NULL,
  end_time        DATETIME(3) NOT NULL,
  booking_status  VARCHAR(32) NOT NULL,
  conflict_reason VARCHAR(255),
  INDEX idx_booking_resource (resource_id),
  INDEX idx_booking_turnaround (turnaround_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS delay_event (
  id                   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  turnaround_id        BIGINT UNSIGNED NOT NULL,
  delay_type           VARCHAR(32) NOT NULL,
  minutes              INT NOT NULL COMMENT '延误分钟，resolved_at 为空时才顺延任务截止',
  root_cause           VARCHAR(255),
  responsibility_team  VARCHAR(32),
  created_at           DATETIME(3),
  resolved_at          DATETIME(3) NULL COMMENT '关闭时间；NULL=未关闭（仍顺延）',
  INDEX idx_delay_turnaround (turnaround_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS app_user (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  username      VARCHAR(64) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role          VARCHAR(32) NOT NULL,
  team_code     VARCHAR(32),
  display_name  VARCHAR(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS audit_log (
  id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  actor       VARCHAR(64),
  action      VARCHAR(64) NOT NULL,
  target_type VARCHAR(32) NOT NULL,
  target_id   VARCHAR(64),
  detail      VARCHAR(512),
  created_at  DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
