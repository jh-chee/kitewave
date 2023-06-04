CREATE TABLE IF NOT EXISTS messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    chat_room TEXT,
    sender TEXT,
    body TEXT,
    created TIMESTAMP DEFAULT NOW()
);
