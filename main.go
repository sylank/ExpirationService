package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/sylank/lavender-commons-go/dynamo"
	"github.com/sylank/lavender-commons-go/properties"
	"github.com/sylank/lavender-commons-go/utils"
)

const (
	EMAIL_TEMPLATE      = "./config/deletion_template.html"
	DATABASE_PROPERTIES = "./config/database_properties.json"
)

type ReservationDynamoModel struct {
	ReservationID string
	FromDate      string
	ToDate        string
	UserId        string
}

func UnmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue, out interface{}) error {
	dbAttrMap := make(map[string]*dynamodb.AttributeValue)
	for k, v := range attribute {
		var dbAttr dynamodb.AttributeValue
		bytes, marshalErr := v.MarshalJSON()
		if marshalErr != nil {
			return marshalErr
		}
		json.Unmarshal(bytes, &dbAttr)
		dbAttrMap[k] = &dbAttr
	}
	return dynamodbattribute.UnmarshalMap(dbAttrMap, out)
}

func handler(ctx context.Context, req events.DynamoDBEvent) {
	dynamoProperties, err := properties.ReadDynamoProperties(DATABASE_PROPERTIES)
	userTableName := dynamoProperties.GetTableName("userData")

	dynamo.CreateConnection(dynamoProperties)

	if err != nil {
		log.Println("Failed to read database properties")
		panic(err)
	}

	for _, record := range req.Records {
		if record.EventName == "REMOVE" {
			fmt.Println("Old image")
			fmt.Println(record.Change.OldImage)

			reservationItem := &ReservationDynamoModel{}

			fmt.Println("Unmarshalling image data")
			err = UnmarshalStreamImage(record.Change.OldImage, reservationItem)
			if err != nil {
				fmt.Println(err)

				return
			}

			log.Println("Reservation item:")
			log.Println(reservationItem)

			proj := expression.NamesList(expression.Name("FullName"), expression.Name("Email"), expression.Name("Phone"), expression.Name("UserId"))
			result, err := dynamo.CustomQuery("UserId", reservationItem.UserId, userTableName, proj)
			if err != nil {
				log.Println("Failed to fetch user data")
				panic(err)
			}

			for _, i := range result.Items {
				item := dynamo.UserModel{}

				err = dynamodbattribute.UnmarshalMap(i, &item)
				if err != nil {
					log.Println("Failed to unmarshall user data record")
					panic(err)
				}

				log.Println("User item:")
				log.Println(item)

				log.Println("Sending transactional mail")
				templateBytes := utils.ReadBytesFromFile(EMAIL_TEMPLATE)
				tempateString := string(templateBytes)

				r := strings.NewReplacer(
					"<from_date>", reservationItem.FromDate,
					"<to_date>", reservationItem.ToDate,
					"<reservation_id>", reservationItem.ReservationID,
				)

				err = SendTransactionalMail(item.Email, "Foglalásod törlésre került", r.Replace(tempateString))
				if err != nil {
					log.Println("Failed to send transactional email")
					panic(err)
				}
			}
		}
	}
}

func main() {
	lambda.Start(handler)
}