package getitem

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

// New returns an instance of deleteItem satisfying the RouteHandler interface
func New() iface.RouteHandler {
	return &getItems{
		db: db.DynamoDB(),
	}
}

// Match returns true if this RouteHandler should handle this request
func (g *getItems) Match(request events.APIGatewayV2HTTPRequest) bool {
	// GET /lists/<list_id>/items/<item_id>
	var re = regexp.MustCompile(`^/lists/([\w-]+)/items/([\w-]+)/?$`)
	return request.RequestContext.HTTP.Method == "GET" && re.MatchString(request.RequestContext.HTTP.Path)
}

// Handle handles this request and returns the response and status code
func (g *getItems) Handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	item, err := g.getItem(request.RequestContext.HTTP.Path)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	return item, http.StatusOK
}

func (g *getItems) getItem(path string) (*data.Item, error) {
	listID, itemID, err := getIDs(path)
	if err != nil {
		return nil, err
	}

	return g.db.GetItem(listID, itemID)
}

func getIDs(path string) (string, string, error) {
	parts := strings.SplitN(path, "/", 6)
	if len(parts) < 5 {
		return "", "", fmt.Errorf("Unable to match path: %s", path)
	}
	return parts[2], parts[4], nil
}
