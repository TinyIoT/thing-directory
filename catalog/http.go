package catalog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tinyiot/thing-directory/wot"
)

const (
	// query parameters
	QueryParamOffset      = "offset"
	QueryParamLimit       = "limit"
	QueryParamJSONPath    = "jsonpath"
	QueryParamSearchQuery = "query"
)

type ValidationResult struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
}

type HTTPAPI struct {
	controller CatalogController
}

func NewHTTPAPI(controller CatalogController, version string) *HTTPAPI {
	return &HTTPAPI{
		controller: controller,
	}
}

// Post handler creates one item
func (a *HTTPAPI) Post(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var td ThingDescription
	if err := json.Unmarshal(body, &td); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Error processing the request:", err.Error())
		return
	}

	if td[wot.KeyThingID] != nil {
		id, ok := td[wot.KeyThingID].(string)
		if !ok || id != "" {
			ErrorResponse(w, http.StatusBadRequest, "Registering with user-defined id is not possible using a POST request.")
			return
		}
	}

	id, err := a.controller.add(td)
	if err != nil {
		switch err.(type) {
		case *ConflictError:
			ErrorResponse(w, http.StatusConflict, "Error creating the resource:", err.Error())
			return
		case *BadRequestError:
			ErrorResponse(w, http.StatusBadRequest, "Invalid registration:", err.Error())
			return
		case *ValidationError:
			ValidationErrorResponse(w, err.(*ValidationError).ValidationErrors)
			return
		default:
			ErrorResponse(w, http.StatusInternalServerError, "Error creating the registration:", err.Error())
			return
		}
	}

	w.Header().Set("Location", id)
	w.WriteHeader(http.StatusCreated)
}

// Put handler updates an existing item (Response: StatusOK)
// If the item does not exist, a new one will be created with the given id (Response: StatusCreated)
func (a *HTTPAPI) Put(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	body, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var td ThingDescription
	if err := json.Unmarshal(body, &td); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Error processing the request:", err.Error())
		return
	}

	if id, ok := td[wot.KeyThingID].(string); !ok || id == "" {
		ErrorResponse(w, http.StatusBadRequest, "Registration without id is not possible using a PUT request.")
		return
	}
	if params["id"] != td[wot.KeyThingID] {
		ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Resource id in path (%s) does not match the id in body (%s)", params["id"], td[wot.KeyThingID]))
		return
	}

	err = a.controller.update(params["id"], td)
	if err != nil {
		switch err.(type) {
		case *NotFoundError:
			// Create a new device with the given id
			id, err := a.controller.add(td)
			if err != nil {
				switch err.(type) {
				case *ConflictError:
					ErrorResponse(w, http.StatusConflict, "Error creating the registration:", err.Error())
					return
				case *BadRequestError:
					ErrorResponse(w, http.StatusBadRequest, "Invalid registration:", err.Error())
					return
				case *ValidationError:
					ValidationErrorResponse(w, err.(*ValidationError).ValidationErrors)
					return
				default:
					ErrorResponse(w, http.StatusInternalServerError, "Error creating the registration:", err.Error())
					return
				}
			}
			w.Header().Set("Location", id)
			w.WriteHeader(http.StatusCreated)
			return
		case *BadRequestError:
			ErrorResponse(w, http.StatusBadRequest, "Invalid registration:", err.Error())
			return
		case *ValidationError:
			ValidationErrorResponse(w, err.(*ValidationError).ValidationErrors)
			return
		default:
			ErrorResponse(w, http.StatusInternalServerError, "Error updating the registration:", err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// Patch updates parts or all of an existing item (Response: StatusOK)
func (a *HTTPAPI) Patch(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	body, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var td ThingDescription
	if err := json.Unmarshal(body, &td); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Error processing the request:", err.Error())
		return
	}

	if id, ok := td[wot.KeyThingID].(string); ok && id == "" {
		if params["id"] != td[wot.KeyThingID] {
			ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Resource id in path (%s) does not match the id in body (%s)", params["id"], td[wot.KeyThingID]))
			return
		}
	}

	err = a.controller.patch(params["id"], td)
	if err != nil {
		switch err.(type) {
		case *NotFoundError:
			ErrorResponse(w, http.StatusNotFound, "Invalid registration:", err.Error())
			return
		case *BadRequestError:
			ErrorResponse(w, http.StatusBadRequest, "Invalid registration:", err.Error())
			return
		case *ValidationError:
			ValidationErrorResponse(w, err.(*ValidationError).ValidationErrors)
			return
		default:
			ErrorResponse(w, http.StatusInternalServerError, "Error updating the registration:", err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get handler get one item
func (a *HTTPAPI) Get(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	td, err := a.controller.get(params["id"])
	if err != nil {
		switch err.(type) {
		case *NotFoundError:
			ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		default:
			ErrorResponse(w, http.StatusInternalServerError, "Error retrieving the registration: ", err.Error())
			return
		}
	}

	b, err := json.Marshal(td)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", wot.MediaTypeThingDescription)
	_, err = w.Write(b)
	if err != nil {
		log.Printf("ERROR writing HTTP response: %s", err)
	}
}

// Delete removes one item
func (a *HTTPAPI) Delete(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	err := a.controller.delete(params["id"])
	if err != nil {
		switch err.(type) {
		case *NotFoundError:
			ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		default:
			ErrorResponse(w, http.StatusInternalServerError, "Error deleting the registration:", err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// List lists entries in paginated format or as a stream
func (a *HTTPAPI) List(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Error parsing the query:", err.Error())
		return
	}

	// pagination is done only when limit is set
	if req.Form.Get(QueryParamLimit) != "" {
		a.listPaginated(w, req)
		return
	} else {
		a.listStream(w, req)
		return
	}
}

func (a *HTTPAPI) listPaginated(w http.ResponseWriter, req *http.Request) {
	var err error
	var limit, offset int

	limitStr := req.Form.Get(QueryParamLimit)
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	offsetStr := req.Form.Get(QueryParamOffset)
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	items, err := a.controller.listPaginate(offset, limit)
	if err != nil {
		switch err.(type) {
		case *BadRequestError:
			ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		default:
			ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	b, err := json.Marshal(items)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", wot.MediaTypeJSONLD)
	_, err = w.Write(b)
	if err != nil {
		log.Printf("ERROR writing HTTP response: %s", err)
	}
}

func (a *HTTPAPI) listStream(w http.ResponseWriter, req *http.Request) {
	//flusher, ok := w.(http.Flusher)
	//if !ok {
	//	panic("expected http.ResponseWriter to be an http.Flusher")
	//}

	w.Header().Set("Content-Type", wot.MediaTypeJSONLD)
	w.Header().Set("X-Content-Type-Options", "nosniff") // tell clients not to infer content type from partial body

	_, err := fmt.Fprintf(w, "[")
	if err != nil {
		log.Printf("ERROR writing HTTP response: %s", err)
	}

	first := true
	for item := range a.controller.iterateBytes(req.Context()) {
		select {
		case <-req.Context().Done():
			log.Println("Cancelled by client.")
			if err := req.Context().Err(); err != nil {
				log.Printf("Client err: %s", err)
				return
			}

		default:
			if first {
				first = false
			} else {
				_, err := fmt.Fprint(w, ",")
				if err != nil {
					log.Printf("ERROR writing HTTP response: %s", err)
				}
			}

			_, err := w.Write(item)
			if err != nil {
				log.Printf("ERROR writing HTTP response: %s", err)
			}
			//time.Sleep(500 * time.Millisecond)
			//flusher.Flush()
		}

	}
	_, err = fmt.Fprintf(w, "]")
	if err != nil {
		log.Printf("ERROR writing HTTP response: %s", err)
	}
}

// SearchJSONPath returns the JSONPath query result
func (a *HTTPAPI) SearchJSONPath(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Error parsing the query: ", err.Error())
		return
	}

	query := req.Form.Get(QueryParamSearchQuery)
	if query == "" {
		ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("No value for %s argument", QueryParamSearchQuery))
		return
	}
	w.Header().Add("X-Request-Query", query)

	b, err := a.controller.filterJSONPathBytes(query)
	if err != nil {
		switch err.(type) {
		case *BadRequestError:
			ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		default:
			ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", wot.MediaTypeJSON)
	w.Header().Set("X-Request-URL", req.RequestURI)
	_, err = w.Write(b)
	if err != nil {
		log.Printf("ERROR writing HTTP response: %s", err)
		return
	}
}
