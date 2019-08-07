package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Achievement struct {
	gorm.Model
	Slug    string   `gorm:"unique" json:"slug"`
	Name    string   `gorm:"unique" json:"name"`
	Desc    string   `json:"desc"`
	Img     string   `json:"img"`
	Members []Member `gorm:"many2many:member_achievements;" json:"members,omitempty"`
}

type Member struct {
	gorm.Model
	Name         string        `gorm:"unique" json:"name"`
	Img          string        `json:"img"`
	Achievements []Achievement `gorm:"many2many:member_achievements;" json:"achievements,omitempty"`
	Teams        []Team        `gorm:"many2many:team_members;" json:"teams,omitempty"`
	Games        []Game        `gorm:"many2many game_members;" json:"games,omitempty"`
	Stats        []Stat        `json:"stats,omitempty"`
}

type Team struct {
	gorm.Model
	Name    string   `gorm:"unique" json:"name"`
	Img     string   `json:"img"`
	Members []Member `gorm:"many2many:team_members;" json:"members,omitempty"`
	Games   []Game   `gorm:"many2many game_teams;" json:"games,omitempty"`
	Stats   []Stat   `json:"stats,omitempty"`
}

type Game struct {
	gorm.Model
	Status  gameStatus `json:"status"`
	Teams   []Team     `gorm:"many2many:game_teams;" json:"teams,omitempty"`
	Members []Member   `gorm:"many2many:game_members;" json:"members,omitempty"`
	Stats   []Stat     `json:"stats,omitempty"`
}

type Stat struct {
	gorm.Model   `json:"-"`
	GameID       uint `json:"-"`
	TeamID       uint `json:"-"`
	MemberID     uint `json:"-"`
	NumAttacks   uint `json:"num_attacks"`
	NumHits      uint `json:"num_hits"`
	AmountDamage uint `json:"amount_damage"`
	NumKills     uint `json:"num_kills"`
	InstantKills uint `json:"instant_kills"`
	NumAssists   uint `json:"num_assists"`
	NumSpells    uint `json:"num_spells"`
	SpellsDamage uint `json:"spells_damage"`
	IsWinner     bool `json:"-"`
}

var db *gorm.DB

func initDatabase() {
	var err error
	db, err = gorm.Open(databaseDialect, databaseFile)
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&Achievement{}, &Member{}, &Team{}, &Game{}, &Stat{})
}
