package getitems

import (
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

type getItems struct {
	db db.DB
}

// New returns an instance of getItems satisfying the RouteHandler interface
func New() iface.RouteHandler {
	return &getItems{
		db: db.DynamoDB(),
	}
}

// Match returns true if this RouteHandler should handle this request
func (g *getItems) Match(request events.APIGatewayV2HTTPRequest) bool {
	var re = regexp.MustCompile(`^/lists/([\w-]+)/items/?$`)
	return request.RequestContext.HTTP.Method == "GET" && re.MatchString(request.RequestContext.HTTP.Path)
}

// Handle handles this request and returns the response and status code
func (g *getItems) Handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	items, err := g.getItems(request.RequestContext.HTTP.Path)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	return items, http.StatusOK
}

func (g *getItems) getItems(path string) (*[]data.Item, error) {
	listID, err := getListID(path)
	if err != nil {
		return nil, err
	}

	return g.db.GetItemsOnList(listID)
}

func getListID(path string) (string, error) {
	parts := strings.SplitN(path, "/", 4)
	if len(parts) < 4 {
		return "", fmt.Errorf("Unable to match path: %s", path)
	}
	return parts[2], nil
}
