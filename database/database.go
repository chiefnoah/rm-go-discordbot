package database






const initTableS = "CREATE IF NOT EXISTS Users(ID INTEGER PRIMARY KEY, Nick TEXT NOT NULL, UserID TEXT NOT NULL, email TEXT, Username TEXT NOT NULL, Avatar TEXT, Discriminator TEXT NOT NULL, Token TEXT, Verified INTEGER DEFAULT 0, Bot INTEGER DEFAULT 0, Points INTEGER DEFAULT 0, PointTimestamp INTEGER DEFAULT0, TempBannedTimestamp INTEGER DEFAULT 0);"