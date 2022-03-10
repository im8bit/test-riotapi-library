package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/im8bit/test-riotapi-library/riot"
)

var tablename = "val_leaderboards"

type LeaderboardDynamoDBItem struct {
	ActId    string `json:"actid"`
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
	Rank     int    `json:"leaderboardRank"`
	Rating   int    `json:"rankedRating"`
	Wins     int    `json:"numberOfWins"`
}

func DropTable(svc dynamodbiface.DynamoDBAPI) error {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(tablename),
	}
	_, err := svc.DeleteTable(params)
	if err != nil {
		return err
	}

	describeParams := &dynamodb.DescribeTableInput{
		TableName: aws.String(tablename),
	}

	if err := svc.WaitUntilTableNotExists(describeParams); err != nil {
		return err
	}

	return nil
}

func CreateTable(svc dynamodbiface.DynamoDBAPI) error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("actid"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("puuid"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("actid"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("puuid"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(tablename),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	describeParams := &dynamodb.DescribeTableInput{
		TableName: aws.String(tablename),
	}

	if err := svc.WaitUntilTableExists(describeParams); err != nil {
		return err
	}

	return nil
}

func FindAll(svc dynamodbiface.DynamoDBAPI, actid string) []LeaderboardDynamoDBItem {
	filt := expression.Name("actid").Equal(expression.Value(actid))
	proj := expression.NamesList(expression.Name("gameName"), expression.Name("leaderboardRank"), expression.Name("numberOfWins"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tablename),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}

	itemList := []LeaderboardDynamoDBItem{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &itemList)
	if err != nil {
		log.Fatalf("Got error unmarshalling: %s", err)
	}

	return itemList
}

func AddLeaderboardItem(svc dynamodbiface.DynamoDBAPI, actid string, playerDtoData riot.PlayerDto) (string, error) {
	var finalPuuid string = "NOT-AVAILABLE"

	if playerDtoData.Puuid != "" {
		finalPuuid = playerDtoData.Puuid
	}

	item := LeaderboardDynamoDBItem{
		ActId:    actid,
		Puuid:    finalPuuid,
		GameName: playerDtoData.GameName,
		TagLine:  playerDtoData.TagLine,
		Rank:     playerDtoData.LeaderboardRank,
		Rating:   playerDtoData.RankedRating,
		Wins:     playerDtoData.NumberOfWins,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println(err.Error())
		return finalPuuid, err
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: &tablename,
	})

	if err != nil {
		fmt.Println(err.Error())
		return finalPuuid, err
	}

	return finalPuuid, nil
}
