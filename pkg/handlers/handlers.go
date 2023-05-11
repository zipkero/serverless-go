package handlers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/zipkero/serverless-go/pkg/user"
	"github.com/zipkero/serverless-go/pkg/validators"
	"net/http"
)

var ErrorMethodNotAllowed = "method not allow"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	if validators.IsEmailValid(email) {
		result, err := user.FetchUser(email, tableName, dynamoClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	} else {
		result, err := user.FetchUsers(tableName, dynamoClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.CreateUser(req, tableName, dynamoClient)
	if err != nil {
		return apiResponse(http.StatusInternalServerError, ErrorBody{ErrorMsg: aws.String(err.Error())})
	}
	return apiResponse(http.StatusCreated, result)
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.UpdateUser(req, tableName, dynamoClient)
	if err != nil {
		return apiResponse(http.StatusInternalServerError, ErrorBody{ErrorMsg: aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	err := user.DeleteUser(req, tableName, dynamoClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, nil)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
