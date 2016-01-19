/* TODO : person_id should be autoincrement
 change that and fix Person and the tests
 */
CREATE TABLE person (
	person_id INTEGER NOT NULL,
	email TEXT NOT NULL UNIQUE,
    pw_hash BYTEA NOT NULL,
	name_first TEXT,
	name_last TEXT,
    
	PRIMARY KEY (person_id)
);

CREATE TABLE decision
(
	decision_id INTEGER NOT NULL,
	person_id INTEGER NOT NULL, /* what should we do */
	name TEXT,
	description TEXT,
	owner_id INTEGER NOT NULL, /* remove this ? */
	stage INTEGER NOT NULL,
    criterion_vote_style TEXT NOT NULL,
    alternative_vote_style TEXT NOT NULL,
	client_settings TEXT,

	FOREIGN KEY (person_id) REFERENCES person, /* what should we do */
	PRIMARY KEY (decision_id)
);

CREATE TABLE facilitates
(
    decision_id INTEGER NOT NULL,
    person_id INTEGER NOT NULL,
    
    FOREIGN KEY (decision_id) REFERENCES decision,
    FOREIGN KEY (person_id) REFERENCES person,
    PRIMARY KEY (decision_id, person_id)
);

CREATE TABLE criterion
(
	criterion_id INTEGER NOT NULL,
	decision_id INTEGER NOT NULL,
	name TEXT,
	weight TEXT,

	PRIMARY KEY (criterion_id),
	FOREIGN KEY (decision_id) REFERENCES decision
);

CREATE TABLE alternative
(
	alternative_id INTEGER NOT NULL,
	decision_id INTEGER NOT NULL,
	name TEXT,
	rating REAL,

	FOREIGN KEY (decision_id) REFERENCES decision,
	PRIMARY KEY (alternative_id)
);

CREATE TABLE criterion_alternative
(
	criterion_id INTEGER NOT NULL,
	alternative_id INTEGER NOT NULL,
	rating REAL,

	PRIMARY KEY (criterion_id, alternative_id),
	FOREIGN KEY (criterion_id) REFERENCES criterion,
	FOREIGN KEY (alternative_id) REFERENCES alternative
);

CREATE TABLE ballot
(
	ballot_id INTEGER NOT NULL,
	decision_id INTEGER NOT NULL,
	secret INTEGER NOT NULL,
	name TEXT,
	email TEXT,

	PRIMARY KEY (ballot_id),
	FOREIGN KEY (decision_id) REFERENCES decision
);

CREATE TABLE vote_criterion
(
	criterion_id INTEGER NOT NULL,
	ballot_id INTEGER NOT NULL,
	weight INTEGER NOT NULL,

	PRIMARY KEY (criterion_id, ballot_id),

	FOREIGN KEY (criterion_id) REFERENCES criterion,
	FOREIGN KEY (ballot_id) REFERENCES ballot
);

CREATE TABLE vote_alternative
(
	criterion_id INTEGER NOT NULL,
	alternative_id INTEGER NOT NULL,
	ballot_id INTEGER NOT NULL,
	rating INTEGER NOT NULL,

	PRIMARY KEY (criterion_id, alternative_id, ballot_id),

	FOREIGN KEY (criterion_id) REFERENCES criterion,
	FOREIGN KEY (alternative_id) REFERENCES alternative,
	FOREIGN KEY (ballot_id) REFERENCES ballot
);
