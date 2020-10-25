package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/mount-joy/thelist-lambda/db"
)

type patchItem struct {
	db db.DB
}

func newPatchItem() RouteHandler {
	return &patchItem{
		db: db.DynamoDB(),
	}
}

func (g *patchItem) match(request events.APIGatewayV2HTTPRequest) bool {
	var re = regexp.MustCompile(`^/lists/([\w-]+)/items/([\w-]+)/?$`)
	return request.RequestContext.HTTP.Method == "PATCH" && re.MatchString(request.RequestContext.HTTP.Path)
}

func (g *patchItem) handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	listID, itemID, err := getIDs(request.RequestContext.HTTP.Path)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusBadRequest
	}

	newName, err := getName(request.Body)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusBadRequest
	}

	item, err := g.db.UpdateItem(listID, itemID, newName)
	if err != nil {
		if aerr, ok := err.(db.Error); ok {
			if aerr.ErrorType == db.ErrorNotFound {
				return nil, http.StatusNotFound
			}
		}
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	return item, http.StatusOK
}

func getName(body string) (*string, error) {
	var item data.Item
	err := json.Unmarshal([]byte(body), &item)

	return &item.Item, err
}

func getIDs(path string) (*string, *string, error) {
	parts := strings.SplitN(path, "/", 6)
	if len(parts) < 5 {
		return nil, nil, fmt.Errorf("Unable to match path: %s", path)
	}
	return &parts[2], &parts[4], nil
}
