DROP TABLE IF EXISTS employees;

CREATE TABLE employees (
  id INT AUTO_INCREMENT PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255),
  email VARCHAR(255)
);

INSERT INTO employees (first_name, last_name, email)
VALUES 
  ("Ammar", "Alaa", "ammar@wave.com"),
  ("Wael", "Subhi", "wael@wave.com"),
  ("Salam", "Adil", "salam@wave.com"),
  ("Yaseen", "Taha", "yaseen@wave.com");