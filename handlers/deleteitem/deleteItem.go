package deleteitem

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/db"
	"github.com/mount-joy/thelist-lambda/handlers/iface"
)

type deleteItem struct {
	db db.DB
}

// New returns an instance of deleteItem satisfying the RouteHandler interface
func New() iface.RouteHandler {
	return &deleteItem{
		db: db.DynamoDB(),
	}
}

// Match returns true if this RouteHandler should handle this request
func (g *deleteItem) Match(request events.APIGatewayV2HTTPRequest) bool {
	var re = regexp.MustCompile(`^/lists/([\w-]+)/items/([\w-]+)/?$`)
	return request.RequestContext.HTTP.Method == "DELETE" && re.MatchString(request.RequestContext.HTTP.Path)
}

// Handle handles this request and returns the response and status code
func (g *deleteItem) Handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	listID, itemID, err := getIDs(request.RequestContext.HTTP.Path)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusBadRequest
	}

	err = g.db.DeleteItem(listID, itemID)
	if err != nil {
		if aerr, ok := err.(db.Error); ok {
			if aerr.ErrorType == db.ErrorNotFound {
				return nil, http.StatusNotFound
			}
		}
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func getIDs(path string) (*string, *string, error) {
	parts := strings.SplitN(path, "/", 6)
	if len(parts) < 5 {
		return nil, nil, fmt.Errorf("Unable to match path: %s", path)
	}
	return &parts[2], &parts[4], nil
}
