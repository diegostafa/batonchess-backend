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
    created_at DATETIME NOT NULL DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
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

CREATE TABLE IF NOT EXISTS users_in_games (
    game_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    is_playing BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (game_id) REFERENCES games(g_id),
    FOREIGN KEY (user_id) REFERENCES users(u_id)
);