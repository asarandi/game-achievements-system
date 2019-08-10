### intro
this is my solution to a coding challenge i received when i applied for a job as a backend engineer.

the project can be summarized as _**"an achievements system for a fictional online game"**_. please see `SPEC.md` for challenge details.

i decided to implement a restful api in go, so this is a web server connected to a database that exposes several endpoints which accept `GET`, `POST`, `PUT` and `DELETE` requests in order to **create**, **read**, **update** and **delete** records

project components:
- server: https://golang.org/pkg/net/http/
- routes: https://github.com/gorilla/mux
- database: https://github.com/jinzhu/gorm


### api

endpoint|GET|POST|PUT|DELETE|notes
-|-|-|-|-|-
**/achievements**                   |&check;|&check;|       |       |**get** all achievements</br>**create** a new achievement</br>_fields:_ `slug` `name` `desc` `img`
/achievements/{id:}                 |&check;|       |&check;|&check;|**get**, **update** or **delete** an achievement
/achievements/{id:}/members         |&check;|       |       |       |**get** members who have an achievement
**/members**                        |&check;|&check;|       |       |**get** all members</br>**create** a new member</br>_fields:_ `name` `img`
/members/{id:}                      |&check;|       |&check;|&check;|**get**, **update** or **delete** a member
/members/{id:}/achievements         |&check;|       |       |       |**get** all achievements of a member
/members/{id:}/teams                |&check;|       |       |       |**get** all teams of a member
/members/{id:}/teams/{id:}          |       |&check;|       |&check;|a member **joins** - or - **leaves** a team
**/teams**                          |&check;|&check;|       |       |**get** all teams</br>**create** a new team</br>_fields:_ `name` `img`
/teams/{id:}                        |&check;|       |&check;|&check;|**get**, **update** or **delete** a team
/teams/{id:}/members                |&check;|       |       |       |**get** all members of a team
/teams/{id:}/members/{id:}          |       |&check;|       |&check;|**add** a member **to** a team</br>**remove** a member **from** a team
**/games**                          |&check;|&check;|       |       |**get** all games</br>**create** a new game
/games/{id:}                        |       |       |       |&check;|**close** a game and run achievement logic
/games/{id:}/stats                  |&check;|       |       |       | **get** all stats for a game
/games/{id:}/teams                  |&check;|       |       |       |**get** all teams that joined a game
/games/{id:}/teams/{id:}            |       |&check;|       |       |**add** a team to a game;</br>team must contain 3-5 members;
/games/{id:}/teams/{id:}/stats      |&check;|       |       |       |**get** game stats for a team
/games/{id:}/members                |&check;|       |       |       |**get** all members of a game
/games/{id:}/winners                |&check;|       |       |       |**get** all winning members of a game
/games/{id:}/members/{id:}/stats    |&check;|       |&check;|       |**get** or **update** game stats for a member;</br>_fields:_ `num_attacks` `num_hits` `amount_damage` `num_kills` `instant_kills` `num_assists` `num_spells` `spells_damage`


### config
by default the server listens on `0.0.0.0:4242` and uses a file `database.sqlite` to store data. these settings are in `config.go` file

### project files

##### achievements.go
- individual functions that are run at the end of each game to check if a player qualifies for an achievement
##### config.go
- program configuration settings: server listening address, database file name, etc
##### game.go
- contains game related logic: adding teams to games, setting winners at the end of a game, etc
##### handlers.go
- functions responsible for each endpoint and request type
##### main.go
- initializes other program components and launches server
##### main_test.go
- unit tests and sample data driver
##### model.go
- data structures and generalized functions for manipulating database records
##### response.go
- errors and json response encoding
##### routes.go
- endpoints, http request methods



### running
before running the program, make sure to install the dependencies, this only has to be done once:

`go get github.com/gorilla/mux`

`go get github.com/jinzhu/gorm`

`go get github.com/jinzhu/gorm/dialects/sqlite`

**to build and run:**

`go build && ./game-achievements-system`

**or just run without building:**

`go run .`

program should print `server ready at: 0.0.0.0:4242`


### testing

file `main_test.go` contains unit tests and also acts as a driver.

the test scenario is simple:
- create the four default achievements
- create 10 sample members (players)
- create 4 sample teams
- start 2 games
- update member stats
- close games
- verify that server awarded correct wins and achievements to qualifying members


### adding new achievements
to address `SPEC.md` requirement #7, "_add a new statistic and create a new achievement .._":

to add a simple statistic like "number of collected coins" during a game we would simply add it to our `Stat` object which is defined in `models.go`; it might look like this:
```go
type Stat struct {
	gorm.Model
	GameID       uint `json:"game_id"`
	TeamID       uint `json:"team_id"`
	MemberID     uint `json:"member_id"`
	NumAttacks   uint `json:"num_attacks"`
	NumHits      uint `json:"num_hits"`
	AmountDamage uint `json:"amount_damage"`
	NumKills     uint `json:"num_kills"`
	InstantKills uint `json:"instant_kills"`
	NumAssists   uint `json:"num_assists"`
	NumSpells    uint `json:"num_spells"`
	SpellsDamage uint `json:"spells_damage"`
	NumCoins     uint `json:"num_coins"`		// <<--- we added a new statistic
	IsWinner     bool `json:"is_winner"`
}
```

next, we want to add the achievement to the database; here is one way of doing that:

```sh
curl -X POST 0.0.0.0:4242/achievements --data '{"name": "1UP", "slug": "oneup", "desc": "You have collected 100 or more coins during a single game! Wow!", "img": "https://www.pngfind.com/pngs/m/22-226955_mario-mushroom-png-mario-1-up-mushroom-transparent.png"}'
```

if successful, the server replies ` "success":true, "code":201 `
```sh
{"success":true,"code":201,"message":"ok","result":{"ID":5,"CreatedAt":"2019-08-09T23:44:41.040203-07:00","UpdatedAt":"2019-08-09T23:44:41.040203-07:00","DeletedAt":null,"slug":"oneup","name":"1UP","desc":"You have collected 100 or more coins during a single game! Wow!","img":"https://www.pngfind.com/pngs/m/22-226955_mario-mushroom-png-mario-1-up-mushroom-transparent.png"}}
```


then, in `achievements.go` we want to add a function which will check if a member qualifies for the achievement; function signature should be `func (Stat) bool` it might look like this:
```go
func isOneUpAward(stat Stat) bool {
	if stat.NumCoins >= 100 {
		return true
	}
	return false
}
```

we also have to "register" the achievement by adding it to the `asf` list which is in `achievements.go`;
```
    asf = []typeAsf{
		{achievement: Achievement{}, slug: "sharpshooter", function: isSharpshooterAward},
		{achievement: Achievement{}, slug: "bruiser", function: isBruiserAward},
		{achievement: Achievement{}, slug: "veteran", function: isVeteranAward},
		{achievement: Achievement{}, slug: "bigwinner", function: isBigwinnerAward},
		{achievement: Achievement{}, slug: "oneup", function: isOneUpAward},		// <<--- 1UP new achievement
	}
```
it is important that the `slug` field matches the `slug` that we previously inserted into database.

last step is to restart the server. at the end of each game the new achievement will be awarded to qualifying members.

##### note:

we can create more interesting achievements by using sql capabilities of the database,

for example to check if a member qualifies for a _"Win-Five-Games-In-A-Row"_ achievement:
```go
func isWonLastFive(stat Stat) bool {
	lastFiveStats := []Stat{}
	db.Order("CreatedAt desc").Limit(5).Joins("JOIN games ON games.ID = stats.game_id AND games.status = ? AND stats.member_id = ?", finishedGame, stat.MemberID).Find(&lastFiveStats)
	if len(lastFiveStats) == 5 {
		result := true
		for i := range lastFiveStats { result = result && lastFiveStats[i].IsWinner }
		return result
	}
	return false;
}
```


### references
- https://en.wikipedia.org/wiki/Create,_read,_update_and_delete
- https://en.wikipedia.org/wiki/Representational_state_transfer
- https://medium.com/studioarmix/learn-restful-api-design-ideals-c5ec915a430f
- http://www.golangprograms.com/golang-restful-api-using-grom-and-gorilla-mux.html
- http://jinzhu.me/gorm/crud.html
- https://pragmacoders.com/blog/building-a-json-api-in-go
- https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b
- https://medium.com/@shazow/how-i-design-json-api-responses-71900f00f2db
- https://github.com/mousey/fifa-soccer-12-ultimate-team-data
