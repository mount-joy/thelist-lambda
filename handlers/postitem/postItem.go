package postitem

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
	"github.com/mount-joy/thelist-lambda/handlers/iface"
)

type postItem struct {
	db db.DB
}

// New returns an instance of postItem satisfying the RouteHandler interface
func New() iface.RouteHandler {
	return &postItem{
		db: db.DynamoDB(),
	}
}

// Match returns true if this RouteHandler should handle this request
func (g *postItem) Match(request events.APIGatewayV2HTTPRequest) bool {
	// POSY /lists/<list_id>/items
	var re = regexp.MustCompile(`^/lists/([\w-]+)/items/?$`)
	return request.RequestContext.HTTP.Method == "POST" && re.MatchString(request.RequestContext.HTTP.Path)
}

// Handle handles this request and returns the response and status code
func (g *postItem) Handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	listID, err := getListID(request.RequestContext.HTTP.Path)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	name, err := getName(request.Body)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusBadRequest
	}

	item, err := g.db.CreateItem(listID, name)
	if err != nil {
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

func getListID(path string) (string, error) {
	parts := strings.SplitN(path, "/", 4)
	if len(parts) < 4 {
		return "", fmt.Errorf("Unable to match path: %s", path)
	}
	return parts[2], nil
}
