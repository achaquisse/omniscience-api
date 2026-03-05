CREATE TABLE `CertificateLog` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `template` varchar(50) NOT NULL,
    `student_name` varchar(500) NOT NULL,
    `cert_description` text NOT NULL,
    `created_by` varchar(500) NOT NULL,
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
