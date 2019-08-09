package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Achievement struct {
	gorm.Model
	Slug    string   `gorm:"unique;not null" json:"slug"`
	Name    string   `gorm:"unique;not null" json:"name"`
	Desc    string   `json:"desc"`
	Img     string   `json:"img"`
	Members []Member `gorm:"many2many:member_achievements;" json:"members,omitempty"`
}

type Member struct {
	gorm.Model
	Name         string        `gorm:"unique;not null" json:"name"`
	Img          string        `json:"img"`
	Achievements []Achievement `gorm:"many2many:member_achievements;" json:"achievements,omitempty"`
	Teams        []Team        `gorm:"many2many:team_members;" json:"teams,omitempty"`
	Games        []Game        `gorm:"many2many game_members;" json:"games,omitempty"`
	Stats        []Stat        `json:"stats,omitempty"`
}

type Team struct {
	gorm.Model
	Name    string   `gorm:"unique;not null" json:"name"`
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
	IsWinner     bool `json:"is_winner"`
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

/*
	get all records where .. condition a b c
*/
func getAllRecordsWhereABC(model, a, b, c interface{}) error {
	var count int
	if err := db.Where(a, b, c).Find(model).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errorRecordNotFound
	}
	return nil
}

func getRecordWhereABC(model, a, b, c interface{}) error {
	return db.Where(a, b, c).First(model).Error
}

func createFromModel(model interface{}) error {
	return db.Create(model).Error
}

func getAllRecords(model interface{}) error {
	var count int
	if err := db.Find(model).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errorRecordNotFound
	}
	return nil
}

func getRecordByID(model, id interface{}) error {
	return db.First(model, id).Error
}

func updateRecordByID(oldRecord, newRecord, id interface{}) error {
	if err := db.First(oldRecord, id).Error; err != nil {
		return err
	}
	return db.Model(oldRecord).Omit("ID","CreatedAt","UpdatedAt","DeletedAt").Updates(newRecord).Error
}

func deleteRecordByID(model, id interface{}) error {
	return db.Delete(model, id).Error
}

func findAssociationRecords(modelA, idA, assoc, modelB interface{}) error {
	var count int
	if err := db.First(modelA, idA).Error; err != nil {
		return err
	}
	if err := db.Model(modelA).Association(assoc.(string)).Find(modelB).Error; err != nil {
		return err
	}
	if count == 0 {
		return errorRecordNotFound
	}
	return nil
}

func appendAssociationRecord(modelA, idA, assoc, modelB, idB interface{}) error {
	if err := db.First(modelA, idA).Error; err != nil {
		return err
	}
	if err := db.First(modelB, idB).Error; err != nil {
		return err
	}
	return db.Model(modelA).Association(assoc.(string)).Append(modelB).Error
}

func deleteAssociationRecord(modelA, idA, assoc, modelB, idB interface{}) error {
	if err := db.First(modelA, idA).Error; err != nil {
		return err
	}
	if err := db.First(modelB, idB).Error; err != nil {
		return err
	}
	return db.Model(modelA).Association(assoc.(string)).Delete(modelB).Error
}
