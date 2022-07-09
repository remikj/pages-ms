package contoller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/remikj/pages-ms/src/service"
	"net/http"
	"strconv"
)

type PageController interface {
	HandlePageGet(writer http.ResponseWriter, request *http.Request)
}

type PageControllerImpl struct {
	PageService service.PageService
}

func NewPageController(pageService service.PageService) *PageControllerImpl {
	return &PageControllerImpl{pageService}
}

func (pc *PageControllerImpl) HandlePageGet(writer http.ResponseWriter, request *http.Request) {
	pageId, err := getPageIdFromRequest(request)
	if err != nil {
		fmt.Println(err)
		handleBadRequest(writer)
		return
	}

	page, err := pc.PageService.GetPage(pageId)
	if err != nil {
		fmt.Println(err)
		handleInternalServerError(writer)
		return
	}

	if page == nil {
		handleNotFoundServerError(writer)
		return
	}

	marshal, err := json.Marshal(page)
	if err != nil {
		fmt.Println(err)
		handleInternalServerError(writer)
		return
	}
	fmt.Printf("Found page with id: %v Page: %v\n", pageId, string(marshal))

	err = writeResponse(writer, marshal)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getPageIdFromRequest(request *http.Request) (int, error) {
	pageIdStr := chi.URLParam(request, "id")
	pageId, err := strconv.Atoi(pageIdStr)
	if err != nil {
		return -1, err
	}
	return pageId, nil
}

func handleNotFoundServerError(writer http.ResponseWriter) {
	writeStatusAndText(writer, http.StatusNotFound, "result not found")
}

func writeResponse(writer http.ResponseWriter, marshal []byte) error {
	writer.Header().Set("Content-Type", "application/json")
	_, err := writer.Write(marshal)
	return err
}

func handleInternalServerError(writer http.ResponseWriter) {
	writeStatusAndText(writer, http.StatusInternalServerError, "Unexpected error")
}

func handleBadRequest(writer http.ResponseWriter) {
	writeStatusAndText(writer, http.StatusBadRequest, "Expected pageId to be number")
}

func writeStatusAndText(writer http.ResponseWriter, status int, text string) {
	writer.WriteHeader(status)
	_, err := writer.Write([]byte(text))
	if err != nil {
		fmt.Println(err)
	}
}
