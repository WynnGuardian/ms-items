CREATE TABLE IF NOT EXISTS WG_AuthenticatedItems (
    Id VARCHAR(24) PRIMARY KEY NOT NULL,
    LastRanked DATETIME NOT NULL,
    ItemName VARCHAR(64) NOT NULL,
    OwnerMCUUID VARCHAR(32) NOT NULL,
    OwnerUserId VARCHAR(20) NOT NULL,
    Position INT NOT NULL,
    TrackingCode VARCHAR(32) NOT NULL,
    OwnerPublic INT NOT NULL,
    Bytes VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS WG_AuthenticatedItemStats (
    ItemId VARCHAR(64) NOT NULL,
    StatId VARCHAR(32) NOT NULL,
    Value INT NOT NULL,

    PRIMARY KEY (ItemId, StatId),
    FOREIGN KEY (ItemId) REFERENCES WG_AuthenticatedItems (Id)
);

CREATE TABLE IF NOT EXISTS WG_Criteria (
    ItemName VARCHAR(64) NOT NULL,
    StatId VARCHAR(32) NOT NULL,
    Value DECIMAL(5,3) NOT NULL,
    PRIMARY KEY (ItemName, StatId)
);

CREATE TABLE IF NOT EXISTS WG_WynnItems (
    Name VARCHAR(64) NOT NULL PRIMARY KEY,
    Sprite VARCHAR(64) NOT NULL,
    ReqLevel INT NOT NULL,
    ReqStrenght INT NOT NULL,
    ReqAgility INT NOT NULL,
    ReqDefence INT NOT NULL,
    ReqIntelligence INT NOT NULL,
    ReqDexterity INT NOT NULL
);  

CREATE TABLE IF NOT EXISTS WG_WynnItemStats (
    ItemName VARCHAR(64) NOT NULL,
    StatId VARCHAR(32) NOT NULL,
    Lower INT NOT NULL,
    Upper INT NOT NULL,

    PRIMARY KEY (ItemName, StatId),
    FOREIGN KEY (ItemName) REFERENCES WG_WynnItems (Name)
);

CREATE TABLE IF NOT EXISTS WG_Surveys (
    Id VARCHAR(24) PRIMARY KEY NOT NULL,
    ChannelID VARCHAR(20) NOT NULL,
    AnnouncementMessageID VARCHAR(20) NOT NULL,
    Status TINYINT NOT NULL,
    ItemName VARCHAR(64) NOT NULL,
    OpenedAt DATETIME NOT NULL,
    Deadline DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS WG_Votes (
    MessageId VARCHAR(20) NOT NULL,
    UserId VARCHAR(20) NOT NULL,
    SurveyId VARCHAR(24) NOT NULL,
    Token VARCHAR(32) NOT NULL,
    Status TINYINT NOT NULL,
    VotedAt DATETIME NOT NULL,

    PRIMARY KEY (UserId, SurveyId),
    FOREIGN KEY (SurveyId) REFERENCES WG_Surveys (Id)
);

CREATE TABLE IF NOT EXISTS WG_VoteEntries (
    SurveyId VARCHAR(24) NOT NULL,
    UserId VARCHAR(20) NOT NULL,
    StatId VARCHAR(32) NOT NULL,
    Value DECIMAL(5,3) NOT NULL,

    PRIMARY KEY (SurveyId, UserId, StatId),
    FOREIGN KEY (SurveyId) REFERENCES WG_Surveys (Id)
);