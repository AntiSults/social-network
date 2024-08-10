CREATE TABLE Users (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Email TEXT NOT NULL UNIQUE,
    Password TEXT NOT NULL,
    FirstName TEXT NOT NULL,
    LastName TEXT NOT NULL,
    DOB DATE,
    NickName TEXT,
    AboutMe TEXT,
    AvatarPath TEXT
);