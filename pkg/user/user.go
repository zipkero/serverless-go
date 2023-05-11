package user

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/zipkero/serverless-go/pkg/validators"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var (
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorInvalidRequestBody      = "invalid request user data"
	ErrorInvalidEmail            = "invalid email"
	ErrorColudNotMarshalItem     = "could not marshal item"
	ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotDBPutItem       = "could not DB put item"
	ErrorUserAlreadyExists       = "user already exists"
	ErrorUserDoesNotExists       = "user does not exists"
)

func FetchUser(email string, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynamoClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func FetchUsers(tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynamoClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	items := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return items, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*User,
	error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidRequestBody)
	}
	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}
	if currentUser, _ := FetchUser(u.Email, tableName, dynamoClient); len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}
	item, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorColudNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	_, err = dynamoClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDBPutItem)
	}
	return &u, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*User,
	error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidRequestBody)
	}
	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}
	if currentUser, _ := FetchUser(u.Email, tableName, dynamoClient); len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExists)
	}
	item, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorColudNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	_, err = dynamoClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDBPutItem)
	}
	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynamoClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil
}
