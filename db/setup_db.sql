CREATE TABLE users (
    u_id TEXT PRIMARY KEY,
    u_name CHECK(LENGTH(u_name) <= 100) NOT NULL DEFAULT 'anon'
);

CREATE TABLE games (
    g_id TEXT PRIMARY KEY AUTOINCREMENT,
    fen TEXT NOT NULL,
    -- minutes_per_side INTEGER NOT NULL,
    -- seconds_increment_per_move INTEGER NOT NULL
    max_players_per_side INTEGER NOT NULL,
    g_status TEXT CHECK(
        outcome IN (
            'NORMAL',
            'CHECKMATE',
            'STALEMATE',
            'REPETITION',
            'INSUFF'
        )
    ) NOT NULL DEFAULT 'NORMAL',
);

CREATE TABLE game_players (
    user_id TEXT NOT NULL,
    game_id INTEGER NOT NULL,
    FOREIGN KEY (game_id) REFERENCES games(g_id),
    FOREIGN KEY (user_id) REFERENCES users(u_id)
);

/*
 user_id
 game_props
 
 games.add(auto_id, creator, props...)
 id = select.where creator=user_id and active
 
 return id
 
 ....
 
 game_id, user_id, join
 game_players.add(user_id, game_id, side)
 return select fen where game id = gid
 
 ...
 
 
 */