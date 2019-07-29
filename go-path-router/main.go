package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const tpl = `<html>
	<head>
		<meta name="go-import"
		      content="{{.Domain}}{{.Path}}
                   git https://{{.CodePath}}{{.Path}}">
		<meta name="go-source"
		      content="{{.Domain}}{{.Path}}
                   https://{{.CodePath}}{{.Path}}
                   https://{{.CodePath}}{{.Path}}/tree/master{/dir}
                   https://{{.CodePath}}{{.Path}}/blob/master{/dir}/{file}#L{line}">
		<meta http-equiv="refresh" content="0; url=https://godoc.org/{{.Domain}}{{.Path}}/">
	</head>
	<body>
		Nothing to see here; <a href="https://godoc.org/{{.Domain}}{{.Path}}/">see the package on godoc</a>.
	</body>
</html>`

type goHTMLData struct {
	Domain   string
	CodePath string
	Path     string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	data := goHTMLData{
		Domain:   os.Getenv("DOMAIN"),
		CodePath: os.Getenv("CODEPATH"),
		Path:     request.Path,
	}

	t, err := template.New("index").Parse(tpl)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("Couldn't construct template")
	}

	var resp bytes.Buffer
	err = t.Execute(&resp, data)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("Couldn't construct gopath %s", err.Error())
	}

	return events.APIGatewayProxyResponse{
		Body:       resp.String(),
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":  "text/html",
			"Cache-Control": "public, max-age=86400",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
