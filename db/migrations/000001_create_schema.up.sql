CREATE TABLE `projects` (
    `id` text,
    `created_at` datetime,
    `updated_at` datetime,
    `name` text,
    `color` text DEFAULT "Zinc",
    `icon` text DEFAULT "Unknown",
    PRIMARY KEY (`id`),
    CONSTRAINT `uni_projects_name` UNIQUE (`name`)
);

CREATE TABLE `tasks` (
    `id` text,
    `created_at` datetime,
    `updated_at` datetime,
    `title` text,
    `description` text,
    `start_time` datetime,
    `duration` integer DEFAULT 0,
    `completed` numeric,
    `rank` integer,
    `project_id` text,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_project-task` FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`)
)