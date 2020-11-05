package patchitem

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/mount-joy/thelist-lambda/db"
	"github.com/mount-joy/thelist-lambda/handlers/iface"
)

type patchItem struct {
	db db.DB
}

// New returns an instance of patchItem satisfying the RouteHandler interface
func New() iface.RouteHandler {
	return &patchItem{
		db: db.DynamoDB(),
	}
}

// Match returns true if this RouteHandler should handle this request
func (p *patchItem) Match(request events.APIGatewayV2HTTPRequest) bool {
	// PATCH /lists/<list_id>/items/<item_id>
	var re = regexp.MustCompile(`^/lists/([\w-]+)/items/([\w-]+)/?$`)
	return request.RequestContext.HTTP.Method == "PATCH" && re.MatchString(request.RequestContext.HTTP.Path)
}

// Handle handles this request and returns the response and status code
func (p *patchItem) Handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
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

	item, err := p.db.UpdateItem(listID, itemID, newName)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, http.StatusNotFound
		}
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	return item, http.StatusOK
}

func getName(body string) (string, error) {
	var item data.Item
	err := json.Unmarshal([]byte(body), &item)

	return item.Name, err
}

func getIDs(path string) (string, string, error) {
	parts := strings.SplitN(path, "/", 6)
	if len(parts) < 5 {
		return "", "", fmt.Errorf("Unable to match path: %s", path)
	}
	return parts[2], parts[4], nil
}
