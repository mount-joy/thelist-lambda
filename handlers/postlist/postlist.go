package postlist

import (
	"log"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/mount-joy/thelist-lambda/db"
	"github.com/mount-joy/thelist-lambda/handlers/iface"
)

type postList struct {
	db db.DB
}

// New returns an instance of postList satisfying the RouteHandler interface
func New() iface.RouteHandler {
	return &postList{
		db: db.DynamoDB(),
	}
}

// Match returns true if this RouteHandler should handle this request
func (p *postList) Match(request events.APIGatewayV2HTTPRequest) bool {
	// POST /lists AND /lists/
	var re = regexp.MustCompile(`^/lists\/?$`)
	return request.RequestContext.HTTP.Method == "POST" && re.MatchString(request.RequestContext.HTTP.Path)
}

// Handle handles creat list requests and returns the response body and status code
func (p *postList) Handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	name, err := data.GetNameFieldInJson(request.Body)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusBadRequest
	}

	list, err := p.db.CreateList(name)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	return list, http.StatusOK
}
