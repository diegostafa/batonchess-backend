CREATE TABLE IF NOT EXISTS users (
    u_id TEXT PRIMARY KEY,
    u_name CHECK(LENGTH(u_name) <= 100) NOT NULL
);

CREATE TABLE IF NOT EXISTS games (
    g_id INTEGER PRIMARY KEY AUTOINCREMENT,
    creator_id TEXT NOT NULL,
    fen TEXT NOT NULL,
    max_players_per_side INTEGER NOT NULL,
    g_state TEXT NOT NULL DEFAULT 'NORMAL',
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

CREATE TABLE IF NOT EXISTS game_players (
    user_id TEXT NOT NULL,
    game_id INTEGER NOT NULL,
    FOREIGN KEY (game_id) REFERENCES games(g_id),
    FOREIGN KEY (user_id) REFERENCES users(u_id)
);

/*
 
 create game
 insert (game id, info, creator)
 return last
 
 
 user_id
 game_props
 
 games.add(auto_id, creator, props...)
 id = select.where creator=user_id and active
 
 user.games
 return id
 
 ....
 
 game_id, user_id, join
 game_players.add(user_id, game_id, side)
 return select fen where game id = gid
 
 ...
 
 
 */