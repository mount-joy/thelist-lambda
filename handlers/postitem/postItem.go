package postitem

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
func (p *postItem) Match(request events.APIGatewayV2HTTPRequest) bool {
	// POST /lists/<list_id>/items
	var re = regexp.MustCompile(`^/lists/([\w-]+)/items/?$`)
	return request.RequestContext.HTTP.Method == "POST" && re.MatchString(request.RequestContext.HTTP.Path)
}

// Handle handles this request and returns the response and status code
func (p *postItem) Handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	listID, err := getListID(request.RequestContext.HTTP.Path)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	name, err := data.GetNameFieldInJson(request.Body)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusBadRequest
	}

	item, err := p.db.CreateItem(listID, name)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	return item, http.StatusOK
}

func getListID(path string) (string, error) {
	parts := strings.SplitN(path, "/", 4)
	if len(parts) < 4 {
		return "", fmt.Errorf("Unable to match path: %s", path)
	}
	return parts[2], nil
}
