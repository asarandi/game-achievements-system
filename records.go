package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// get record where .. condition a b c
func getRecordsWhereABC(w http.ResponseWriter, r *http.Request, model, a, b, c interface{}) {
	var count int
	if err := db.Where(a, b, c).Find(model).Count(&count).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	if count == 0 {
		jsonResponse(w, Response{false, http.StatusNotFound, "record not found", nil})
		return
	}
	jsonResponse(w, Response{true, http.StatusOK, "ok", model})
}

// update record where .. condition a b c
func updateRecordWhereABC(w http.ResponseWriter, r *http.Request, model, a, b, c interface{}) {
	if err := db.Where(a, b, c).First(model).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(model); err != nil {
		jsonResponse(w, Response{false, http.StatusBadRequest, err.Error(), nil})
		return
	}
	if err := db.Save(model).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	jsonResponse(w, Response{true, http.StatusOK, "ok", model})
}

// create new record
func createFromData(w http.ResponseWriter, r *http.Request, model interface{}) {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(model); err != nil {
		jsonResponse(w, Response{false, http.StatusBadRequest, err.Error(), nil})
		return
	}
	if err := db.Create(model).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	jsonResponse(w, Response{true, http.StatusCreated, "ok", model})
}

func createEmpty(w http.ResponseWriter, r *http.Request, model interface{}) {
	if err := db.Create(model).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	jsonResponse(w, Response{true, http.StatusCreated, "ok", model})
}

// get all records
func getAllRecords(w http.ResponseWriter, r *http.Request, model interface{}) {
	if err := db.Find(model).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	jsonResponse(w, Response{true, http.StatusOK, "ok", model})
}

// get record by id
func getRecordByID(w http.ResponseWriter, r *http.Request, model interface{}) {
	vars := mux.Vars(r)
	if err := db.First(model, vars["id0"]).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	jsonResponse(w, Response{true, http.StatusOK, "ok", model})
}

// update record by id
// logic: first find existing record to populate struct with all extras: ID, CreatedAt, UpdatedAt, etc
// load updated data from request body into same struct and then save
func updateRecordByID(w http.ResponseWriter, r *http.Request, model interface{}) {
	vars := mux.Vars(r)
	if err := db.First(model, vars["id0"]).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(model); err != nil {
		jsonResponse(w, Response{false, http.StatusBadRequest, err.Error(), nil})
		return
	}
	if err := db.Save(model).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	jsonResponse(w, Response{true, http.StatusAccepted, "ok", model})
}

// delete record by id
func deleteRecordByID(w http.ResponseWriter, r *http.Request, model interface{}) {
	vars := mux.Vars(r)
	if err := db.First(model, vars["id0"]).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	if err := db.Delete(model).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	jsonResponse(w, Response{true, http.StatusNoContent, "ok", nil})
}

// get association records
func getAssociationRecords(w http.ResponseWriter, r *http.Request, a, b, c interface{}) {
	vars := mux.Vars(r)
	if err := db.First(a, vars["id0"]).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, errorMessage, nil})
		return
	}
	db.Model(a).Association(b.(string)).Find(c)
	jsonResponse(w, Response{true, http.StatusOK, "ok", c})
}

// add association records
func addAssociationRecord(w http.ResponseWriter, r *http.Request, a, b, c interface{}) {
	vars := mux.Vars(r)
	if err := db.First(a, vars["id0"]).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id0"], errorMessage), nil})
		return
	}
	if err := db.First(c, vars["id1"]).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id1"], errorMessage), nil})
		return
	}
	db.Model(a).Association(b.(string)).Append(c)
	jsonResponse(w, Response{true, http.StatusOK, "ok", nil})
}

// remove association records
func removeAssociationRecord(w http.ResponseWriter, r *http.Request, a, b, c interface{}) {
	vars := mux.Vars(r)
	if err := db.First(a, vars["id0"]).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id0"], errorMessage), nil})
		return
	}
	if err := db.First(c, vars["id1"]).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id1"], errorMessage), nil})
		return
	}
	db.Model(a).Association(b.(string)).Delete(c)
	jsonResponse(w, Response{true, http.StatusOK, "ok", nil})
}
