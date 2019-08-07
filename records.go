package main

import (
	"errors"
)

// get record where .. condition a b c
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

// update record where .. condition a b c
func updateRecordWhereABC(oldRecord, newRecord, a, b, c interface{}) error {
	if err := db.Where(a, b, c).First(oldRecord).Error; err != nil {
		return err
	}
	return db.Model(oldRecord).Omit("ID","CreatedAt","UpdatedAt","DeletedAt").Updates(newRecord).Error
}

func createFromModel(model interface{}) error {
	return db.Create(model).Error
}

// get all records
func getAllRecords(model interface{}) error {
	return db.Find(model).Error
}

// get record by id
func getRecordByID(model, id interface{}) error {
	return db.First(model, id).Error
}

func updateRecordByID(oldRecord, newRecord, id interface{}) error {
	if err := db.First(oldRecord, id).Error; err != nil {
		return err
	}
	return db.Model(oldRecord).Omit("ID","CreatedAt","UpdatedAt","DeletedAt").Updates(newRecord).Error
}

// delete record by id
func deleteRecordByID(model, id interface{}) error {
	return db.Delete(model, id).Error
}

// find association records
func findAssociationRecords(modelA, idA, assoc, modelB interface{}) error {
	if err := db.First(modelA, idA).Error; err != nil {
		return err
	}
	return db.Model(modelA).Association(assoc.(string)).Find(modelB).Error
}

// append association records
func appendAssociationRecord(modelA, idA, assoc, modelB, idB interface{}) error {
	if err := db.First(modelA, idA).Error; err != nil {
		return err
	}
	if err := db.First(modelB, idB).Error; err != nil {
		return err
	}
	return db.Model(modelA).Association(assoc.(string)).Append(modelB).Error
}

// delete association records
func deleteAssociationRecord(modelA, idA, assoc, modelB, idB interface{}) error {
	if err := db.First(modelA, idA).Error; err != nil {
		return err
	}
	if err := db.First(modelB, idB).Error; err != nil {
		return err
	}
	return db.Model(modelA).Association(assoc.(string)).Delete(modelB).Error
}
