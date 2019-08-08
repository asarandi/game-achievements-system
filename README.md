#### For your coding challenge, you will create an achievements system for a fictional online game.
Our fictional game is played online between two teams. Each team can have 3-5 players (but must be evenly matched - i.e., if one team has 4 players, the other must have 4 players). Statistics are compiled for each player during/after each game:

- Number of attempted attacks
- Number of hits
- Total amount of damage done
- Number of kills
- Number of "first hit" kills
- Number of assists
- Total number of spells cast
- Total spell damage done
- Total time played

In addition, historical player statistics are maintained as well. For instance, each player is credited with a win/loss based upon the performance of his/her team during the game. Other historical player-level stats can be tracked as well - such as total number of games played, total duration of games played, total number of kills, etc. Achievements are awarded based upon any combination of these statistics (historical OR game specific). Players collect achievements over time, and these will be displayed on their profile page (the display is outside the scope of this exercise). This profile page is visible to all other players, so achievements act as a kind of "badge of honor". After each game has ended, the achievements logic is called and achievements are awarded to each of the players if applicable.

Your mission, if you choose to accept it, is to create this achievements system. So, here's what's expected from this coding exercise:
1. Code should be written in Go, Java, C#, or C++.
2. The achievements system that you are creating should correctly assign achievements to each player involved in the game that just ended, if applicable.
3. Persistence, messaging, and other potential concerns are considered ancillary to this exercise, so don't worry about implementing those.
4. Create all domain objects suggested in the above description, using the proper data structures where applicable.
5. Create the achievements system. At a minimum, it should have an entry point that would be called at the end of a game. You must include the following achievements:
   - "Sharpshooter" Award - a user receives this for landing 75% of their attacks, assuming they have at least attacked once.
   -  "Bruiser" Award - a user receives this for doing more than 500 points of damage in one game
   - "Veteran" Award - a user receives this for playing more than 1000 games in their lifetime.
   - "Big Winner" Award - a user receives this for having 200 wins
6. Your system should be extensible in the following ways:
   - The system must easily handle the tracking of additional statistics at any level listed above (historical by player or game-specific by player)
   - The system must easily handle the addition of new achievements
7. Add a new statistic, and create a new achievement utilizing this statistic in conjunction to one of the other statistics defined above.
8. Create unit tests to demonstrate that your code is functional using (e.g. using JUnit or equivalent for C#/C++)
9. Create a driver class that sets up sample data and calls the achievements system, printing out achievements to the console. Don't worry about persistence, and use any method of creating/loading sample data that you are familiar with. Also, you can combine/integrate this into your unit tests if you prefer - we're just looking for you to demonstrate that the code runs.

---
##### endpoint reference


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


##### dependencies
`go get github.com/gorilla/mux`</br>
`go get github.com/jinzhu/gorm`</br>
`go get github.com/jinzhu/gorm/dialects/sqlite`</br>

##### references
- https://en.wikipedia.org/wiki/Create,_read,_update_and_delete
- https://en.wikipedia.org/wiki/Representational_state_transfer
- https://medium.com/studioarmix/learn-restful-api-design-ideals-c5ec915a430f
- http://www.golangprograms.com/golang-restful-api-using-grom-and-gorilla-mux.html
- http://jinzhu.me/gorm/crud.html
- https://pragmacoders.com/blog/building-a-json-api-in-go
- https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b
- https://medium.com/@shazow/how-i-design-json-api-responses-71900f00f2db
- https://github.com/mousey/fifa-soccer-12-ultimate-team-data

