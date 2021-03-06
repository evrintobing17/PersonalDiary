package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/evrintobing17/PersonalDiary/domain"
	"github.com/evrintobing17/PersonalDiary/entry/controller/format"
	"github.com/evrintobing17/PersonalDiary/entry/usecase"
	"github.com/evrintobing17/PersonalDiary/util/encode"
	"github.com/evrintobing17/PersonalDiary/util/middleware"
	"github.com/evrintobing17/PersonalDiary/util/token"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

// EntryController represents all the http request of an entry made by the user
type EntryController struct {
	entryUC usecase.EntryUseCase
	pool    *redis.Pool
}

// NewEntryController creates an object of UserController
func NewEntryController(router *mux.Router, pool *redis.Pool, es usecase.EntryUseCase) {
	controller := &EntryController{
		entryUC: es,
		pool:    pool,
	}

	router.HandleFunc("/entries", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(controller.CreateEntry))).Methods("POST")
	router.HandleFunc("/entries/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(controller.UpdateEntry))).Methods("PUT")
	router.HandleFunc("/entries/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(controller.DeleteEntry))).Methods("DELETE")
	router.HandleFunc("/entries/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(controller.GetEntryOfUserByID))).Methods("GET")
	router.HandleFunc("/entries", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(controller.GetAllEntriesOfUser))).Methods("GET")

}

// CreateEntry endpoint is used to create an entry
func (ec *EntryController) CreateEntry(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	entry := domain.Entry{}
	err = json.Unmarshal(body, &entry)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	entry.Initialize()
	err = format.Validate(&entry)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Is the user authenticated?
	authDetails, err := token.ExtractTokenMetaData(r)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	uid, err := token.FetchAuthDetails(authDetails, ec.pool)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if authenticated user is the owner of this entry
	if uid != entry.OwnerID {
		encode.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	createdEntry, err := ec.entryUC.CreateEntry(&entry)
	if err != nil {
		encode.ERROR(w, http.StatusInternalServerError, errors.New("Error while creating an entry"))
		fmt.Println("CreateEntry", err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, createdEntry.ID))
	encode.JSON(w, http.StatusCreated, createdEntry)
}

// UpdateEntry endpoint
func (ec *EntryController) UpdateEntry(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check for valid entry id
	eid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		encode.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is the user authenticated?
	authDetails, err := token.ExtractTokenMetaData(r)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	uid, err := token.FetchAuthDetails(authDetails, ec.pool)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the entry exist
	entry, err := ec.entryUC.GetEntryOfUserByID(eid, uid)
	if err != nil {
		encode.ERROR(w, http.StatusNotFound, errors.New("Entry doesn't exist"))
		return
	}

	// If a user attempt to update an entry not belonging to him
	if uid != entry.OwnerID {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	entryUpdate := domain.Entry{}
	err = json.Unmarshal(body, &entryUpdate)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Check if authenticated user is the owner of this entry
	if uid != entryUpdate.OwnerID {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	entryUpdate.Initialize()
	err = format.Validate(&entryUpdate)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// set the entry ID to the updated objects ID
	entryUpdate.ID = entry.ID

	entryUpdated, err := ec.entryUC.UpdateEntry(eid, &entryUpdate)
	if err != nil {
		encode.ERROR(w, http.StatusInternalServerError, errors.New("Error while updating the entry"))
		fmt.Println("UpdateEntry", err)
		return
	}
	encode.JSON(w, http.StatusOK, entryUpdated)
}

// DeleteEntry endpoint
func (ec *EntryController) DeleteEntry(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check for valid entry id
	eid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		encode.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is the user authenticated?
	authDetails, err := token.ExtractTokenMetaData(r)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	uid, err := token.FetchAuthDetails(authDetails, ec.pool)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the entry exist
	entry, err := ec.entryUC.GetEntryOfUserByID(eid, uid)
	if err != nil {
		encode.ERROR(w, http.StatusNotFound, errors.New("Entry doesn't exist"))
		return
	}

	// Check if authenticated user is the owner of this entry
	if uid != entry.OwnerID {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	_, err = ec.entryUC.DeleteEntry(eid, uid)
	if err != nil {
		encode.ERROR(w, http.StatusBadRequest, err)
		fmt.Println("DeleteEntry", err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", eid))
	encode.JSON(w, http.StatusNoContent, "")
}

// GetEntryOfUserByID endpoint
func (ec *EntryController) GetEntryOfUserByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	eid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		encode.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is the user authenticated?
	authDetails, err := token.ExtractTokenMetaData(r)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	uid, err := token.FetchAuthDetails(authDetails, ec.pool)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	entry, err := ec.entryUC.GetEntryOfUserByID(eid, uid)
	if err != nil {
		encode.ERROR(w, http.StatusInternalServerError, err)
		fmt.Println("GetEntryOfUserByID", err)
		return
	}

	encode.JSON(w, http.StatusOK, entry)
}

// GetAllEntriesOfUser endpoint
func (ec *EntryController) GetAllEntriesOfUser(w http.ResponseWriter, r *http.Request) {

	// Is the user authenticated?
	authDetails, err := token.ExtractTokenMetaData(r)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	uid, err := token.FetchAuthDetails(authDetails, ec.pool)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	params := r.URL.Query()

	// Check for page limit
	limit, err := strconv.ParseUint(params.Get("limit"), 10, 32) // (string, base, bitSize)
	if err != nil {
		limit = 10 // default limit
	}

	// Check for page number
	pageNumber, err := strconv.ParseUint(params.Get("page"), 10, 32)
	if err != nil {
		pageNumber = 1 // default page number
	}
	if pageNumber < 1 {
		pageNumber = 1 // Incorrect page parameter or just return an error
	}

	yearFilter1, err := strconv.ParseUint(params.Get("year[gte]"), 10, 32)
	if err != nil {
		yearFilter1 = 0
	}

	yearFilter2, err := strconv.ParseUint(params.Get("year[lte]"), 10, 32)
	if err != nil {
		yearFilter2 = 0
	}

	year, err := strconv.ParseUint(params.Get("year"), 10, 32)
	if err == nil { // if only 1 year is present
		yearFilter1 = year
	}

	sort := params.Get("sort")

	entries, err := ec.entryUC.GetAllEntriesOfUser(uid, uint32(limit), uint32(pageNumber), uint32(yearFilter1), uint32(yearFilter2), sort)
	if err != nil {
		encode.ERROR(w, http.StatusInternalServerError, err)
		fmt.Println("GetAllEntriesOfUser", err)
		return
	}
	encode.JSON(w, http.StatusOK, entries)
}
