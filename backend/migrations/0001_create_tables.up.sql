-- Enums
CREATE TYPE friend_status AS ENUM ('pending', 'accepted', 'blocked');
CREATE TYPE game_status AS ENUM ('waiting', 'in_progress', 'finished', 'cancelled');
CREATE TYPE invite_status AS ENUM ('pending', 'accepted', 'declined');

-- Tables
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE friends (
    user_id INT NOT NULL,
    friend_id INT NOT NULL,
    status friend_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, friend_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_friend FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_no_self_friend CHECK (user_id <> friend_id)
);

CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    host_id INT NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    private BOOLEAN NOT NULL DEFAULT FALSE,
    status game_status NOT NULL DEFAULT 'waiting',
    CONSTRAINT fk_host FOREIGN KEY (host_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE rounds (
    id SERIAL PRIMARY KEY,
    game_id INT NOT NULL,
    status game_status NOT NULL DEFAULT 'waiting',
    word VARCHAR(127) NOT NULL,
    spy_word VARCHAR(127) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_round_game FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE
);

CREATE TABLE invites (
    id SERIAL PRIMARY KEY,
    round_id INT NOT NULL,
    inviter_id INT NOT NULL,
    invitee_id INT NOT NULL,
    status invite_status NOT NULL DEFAULT 'pending',
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_invite_round FOREIGN KEY (round_id) REFERENCES rounds(id) ON DELETE CASCADE,
    CONSTRAINT fk_inviter FOREIGN KEY (inviter_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_invitee FOREIGN KEY (invitee_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_no_self_invite CHECK (inviter_id <> invitee_id),
    UNIQUE (round_id, invitee_id)
);

CREATE TABLE game_participants (
    round_id INT NOT NULL,
    user_id INT NOT NULL,
    is_spy BOOLEAN NOT NULL DEFAULT FALSE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (round_id, user_id),
    CONSTRAINT fk_gp_round FOREIGN KEY (round_id) REFERENCES rounds(id) ON DELETE CASCADE,
    CONSTRAINT fk_gp_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
