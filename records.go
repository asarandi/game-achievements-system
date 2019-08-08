package main

import (
	"errors"
)

/*
	get all records where .. condition a b c
*/
func getAllRecordsWhereABC(model, a, b, c interface{}) error {
	var count int
	if err := db.Where(a, b, c).Find(model).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("record not found")
	}
	return nil
}

func getRecordWhereABC(model, a, b, c interface{}) error {
	return db.Where(a, b, c).First(model).Error
}

/*
	update record where .. condition a b c
*/
/*
func updateRecordWhereABC(oldRecord, newRecord, a, b, c interface{}) error {
	if err := db.Where(a, b, c).First(oldRecord).Error; err != nil {
		return err
	}
	return db.Model(oldRecord).Omit("ID","CreatedAt","UpdatedAt","DeletedAt").Updates(newRecord).Error
}
*/

func createFromModel(model interface{}) error {
	return db.Create(model).Error
}

func getAllRecords(model interface{}) error {
	return db.Find(model).Error
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
	if err := db.First(modelA, idA).Error; err != nil {
		return err
	}
	return db.Model(modelA).Association(assoc.(string)).Find(modelB).Error
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
