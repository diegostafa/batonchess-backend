CREATE TABLE IF NOT EXISTS users (
    u_id TEXT PRIMARY KEY,
    u_name CHECK(LENGTH(u_name) <= 100) NOT NULL
);

CREATE TABLE IF NOT EXISTS games (
    g_id INTEGER PRIMARY KEY AUTOINCREMENT,
    g_state TEXT NOT NULL DEFAULT 'NORMAL',
    creator_id TEXT NOT NULL,
    fen TEXT NOT NULL,
    max_players INTEGER NOT NULL,
    created_at INTEGER(4) DEFAULT (cast(strftime('%s', 'now') AS INT)),
    CHECK(
        g_state IN (
            'NORMAL',
            'CHECKMATE',
            'STALEMATE',
            'REPETITION',
            'INSUFF'
        )
    ),
    FOREIGN KEY (creator_id) REFERENCES users(u_id)
);