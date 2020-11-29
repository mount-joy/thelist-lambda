package getlist

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

type getList struct {
	db db.DB
}

// New returns an instance of deleteItem satisfying the RouteHandler interface
func New() iface.RouteHandler {
	return &getList{
		db: db.DynamoDB(),
	}
}

// Match returns true if this RouteHandler should handle this request
func (g *getList) Match(request events.APIGatewayV2HTTPRequest) bool {
	// GET /lists/<list_id>
	var re = regexp.MustCompile(`^/lists/([\w-]+)/?$`)
	return request.RequestContext.HTTP.Method == "GET" && re.MatchString(request.RequestContext.HTTP.Path)
}

// Handle handles this request and returns the response and status code
func (g *getList) Handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	item, err := g.getList(request.RequestContext.HTTP.Path)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	return item, http.StatusOK
}

func (g *getList) getList(path string) (*data.List, error) {
	listID, err := getID(path)
	if err != nil {
		return nil, err
	}

	return g.db.GetList(listID)
}

func getID(path string) (string, error) {
	parts := strings.SplitN(path, "/", 4)
	if len(parts) < 3 {
		return "", fmt.Errorf("Unable to match path: %s", path)
	}
	return parts[2], nil
}
