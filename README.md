# Backend

This is the backend written in [Golang](http://golang.org). This verison is
almost complete. It's almost same as our initial design. I will explain briefly
explain how the api works. There are also javascript examples written in
`webtest/` to show how to do various things.

# Installing and building

First make sure go is installed and configured then type

`go get`

to get the dependencies

```
go build
./gobackend
```

# Code documentation, and testing

Run

```
godoc -http=localhost:6969
```

open the link in a browser find gobackend package and read the documentation.
To test the code just run the following 

```
./run_tests.sh
```

This command will run all testing and create a system.html file that can be
previewed a the browser, to see which code is covered.

# API

For all API below if the header accept json then json objects are sent as mentioned below,
if the accept header wants html then the same objects are passed to static file that
runs a javascript file with that name found in `static/` folder.
eg : Requestion html from `GET /persons` will give an html reply that will run
the file `static/persons_list.js`

## person

These are the things we can do with `person`, under this table is description
of some of the things we need to send or the things we receive back.


| Description            | URL                          | Method | Wants   | Gives               |
|------------------------|------------------------------|--------|---------|---------------------|
| Create a user          | /person                      | POST   | p1      | p2                  |
| Get all persons        | /persons                     | GET    | nothing | array of p2         |
| Get person info        | /person/:person_id/info      | GET    | nothing | p2                  |
| Get person's decisions | /person/:person_id/decisions | GET    | nothing | a persons decisions |
| Delete a person        | /person/:person_id           | DELETE | nothing | s1                  |
| Update a user          | /person/:person_id           | PUT    | p1      | p2                  |

Wants/Gives

```
p1 = {"email":<str>,"pw_hash":<str>,"name_first":<str>,"name_last":<str>}
p2 = {"person_id":<int>, "email":<str>,"pw_hash":<str>,"name_first":<str>,"name_last":<str>}
s1 = {"result": "deleted"}
```
What we actually return is `{"somehting": object}` for example :
- Get person will return

```
{"person": p1 object}
```

- Get persons will return

```
{"persons": array of p1 object}
```

and so on..


## decision

These are the things we can do with `decision`. Things like creating a decision, configuring it's criterions, creating ballots, doing voting on the ballots, and gather statistics.

| Description                  | URL                                                                           | Method | Wants   | Gives       |
|------------------------------|-------------------------------------------------------------------------------|--------|---------|-------------|
| Create a decision            | /decision                                                                     | POST   | d1      | d2          |
| Get all decisions            | /decisions                                                                    | GET    | nothing | array of d2 |
| Decision update              | /decision/:decision_id                                                        | PUT    | d1      | d2          |
| Get a decision info          | /decision/:decision_id/info                                                   | GET    | nothing | d2          |
| Get a decision statistics    | /decision/:decision_id/stats                                                  | GET    | nothing | s1          |
| Delete a decision            | /decision/:decision_id                                                        | DELETE | nothing | r1          |
| Create a ballot for decision | /decision/:decision_id/ballot                                                 | POST   | b1      | b2          |
| List ballots in a decision   | /decision/:decision_id/ballots                                                | GET    | nothing | array of b2 |
| Show a ballot info           | /decision/:decision_id/ballot/:ballot_id/info                                 | GET    | nothing | b2          |
| Delete a ballot              | /decision/:decision_id/ballot/:ballot_id                                      | DELETE | nothing | r1          |
| Update a ballot              | /decision/:decision_id/ballot/:ballot_id                                      | PUT    | b1      | b2          |
| Create a decision criterion  | /decision/:decision_id/criterion                                              | POST   | c1      | c2          |
| List all criterions          | /decision/:decision_id/criterions                                             | GET    | nothing | array of c2 |
| Get a criterion info         | /decision/:decision_id/criterion/:criterion_id/info                           | GET    | nothing | c2          |
| Update a criterion           | /decision/:decision_id/criterion/:criterion_id                                | PUT    | c1      | c2          |
| Delete a criterion           | /decision/:decision_id/criterion/:criterion_id                                | DELETE | nothing | r1          |
| Ballot votes in a criterion  | /decision/:decision_id/ballot/:ballot_id/criterion/:criterion_id/vote/:weight | GET    | nothing | v2          |
| Update Ballot vote           | /decision/:decision_id/ballot/:ballot_id/criterion/:criterion_id/vote/:weight | PUT    | nothing | v2          |
| Show ballot votes            | /decision/:decision_id/ballot/:ballot_id/votes                                | GET    | nothing | array of v2 |
| Delete a vote in ballot      | /decision/:decision_id/ballot/:ballot_id/criterion/:criterion_id/vote         | DELETE | nothing | r1          |

Wants/Gives

```
d1 = {"person_id":<int>, "name":<str>, "description":<str>, "stage":<int>, "criterion_vote_style":<str>, "alternative_vote_style":<str>, "client_settings":<str>}
d2 = {"decision_id":<int>, "person_id":<int>, "name":<str>, "description":<str>, "stage":<int>, "criterion_vote_style":<str>, "alternative_vote_style":<str>, "client_settings":<str>}
s1 = "undecided yet"
b1 = {"name":<str>, "email":<str>}
b2 = {"ballot_id":<int>, "decision_id":<int>, "secret":<str>, "name":<str>, "email":<str>}
c1 = {"name":<str>, "weight":<str>}
c2 = {"criterion_id":<int>, "decision_id":<int>, "name":<str>, "weight":<str>}
v1 = {"weight":<int>}
v2 = {"criterion_id":<int>, "ballot_id":<int>, "weight":<int>}
r1 = {"result": "deleted"}
```

What we actually return is `{"somehting": object}` for example :
- Get decision will return

```
{"decision": d1 object}
```

- Get decisions will return

```
{"decisions": array of d1 object}
```

and so on..


# Authentication

Authentication is implemented as middlewares applied to routes, currently they are applied
to 0 routes to make it easier for the front-end to start writing their code. The login
you get a cookie and this cookie is checked..etc

| Description                             | URL         | Method | Wants   | Gives   |
|-----------------------------------------|-------------|--------|---------|---------|
| Login                                   | /login      | POST   | l1      | l2      |
| Logout                                  | /logout     | GET    | Nothing | Nothing |
| Get person_id of current logged in user | /whoami     | GET    | Nothing | g1      |

Wants/Gives

```
l1 = {"email":<str>, "password":<str>}
l2 = {"error":<str>} or {"status":<str>}
g1 = {"person_id": <int>} or {"error": <str> }
```

# Missing things

- Handlers errors does not respect html and json replies
- HballotLogin does not implement html and json replies
- Applying the permission middlewares to routes
- SMTP Server for ballots invites
- Ballot route to send invites to the ballot with secret link

# Features

- CORS
- Cross-platform
- Very fast
- Almost done
- Fully concurrent
- SecureCookie Authentication
- Decision Statistics
- Replies html/json

# Changed things from previous design

- Routes are a little different but makes more sense this way IMO
- Alternative voting removed
- Facilitates table removed
- Only one person can own a decision
- All primary keys are auto incremented
- Schema generated by the application
- Backend logic handles DB table restrictions instead of the DB
