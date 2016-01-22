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

	PRIMARY KEY (decision_id)
	FOREIGN KEY (person_id) REFERENCES person,
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

/*
 A decision has a criterion
*/
CREATE TABLE criterion
(
	criterion_id INTEGER NOT NULL,
	decision_id INTEGER NOT NULL,
	name TEXT,
	weight TEXT,

	PRIMARY KEY (criterion_id),
	FOREIGN KEY (decision_id) REFERENCES decision
);

CREATE TABLE vote_criterion
(
	criterion_id INTEGER NOT NULL,
	ballot_id INTEGER NOT NULL,
	weight INTEGER NOT NULL,

	/*PRIMARY KEY (criterion_id, ballot_id),*/
	/*PRIMARY KEY (criterion_id),*/

	FOREIGN KEY (criterion_id) REFERENCES criterion,
	FOREIGN KEY (ballot_id) REFERENCES ballot
);
