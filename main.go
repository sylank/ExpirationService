package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/sylank/lavender-commons-go/dynamo"
	"github.com/sylank/lavender-commons-go/properties"
	"github.com/sylank/lavender-commons-go/utils"
)

const (
	EMAIL_TEMPLATE          = "./config/deletion_template.html"
	DATABASE_PROPERTIES     = "./config/database_properties.json"
	EXPIRATION_TRESHOLD_SEC = 50
)

type ReservationDynamoModel struct {
	ReservationID string
	FromDate      string
	ToDate        string
	UserId        string
	Expiring      int64
}

func unmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue) *ReservationDynamoModel {
	reservationItem := &ReservationDynamoModel{}
	reservationItem.ReservationID = attribute["ReservationId"].String()
	reservationItem.FromDate = attribute["FromDate"].String()
	reservationItem.ToDate = attribute["ToDate"].String()
	reservationItem.UserId = attribute["UserId"].String()
	expiring, err := attribute["Expiring"].Integer()

	if err != nil {
		panic("Failed to convert dynamo integer to int64")
	}

	reservationItem.Expiring = expiring

	return reservationItem
}

func isExpired(expiration int64) bool {
	nowSec := time.Now().Unix()
	fmt.Println(nowSec)
	fmt.Println(expiration)
	if nowSec-EXPIRATION_TRESHOLD_SEC <= expiration && expiration <= nowSec {
		return true
	}

	return false
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

			fmt.Println("Unmarshalling image data")
			reservationItem := unmarshalStreamImage(record.Change.OldImage)
			if err != nil {
				panic(err)
			}

			log.Println("Reservation item:")
			log.Println(reservationItem)

			if isExpired(reservationItem.Expiring) {
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
			} else {
				fmt.Println("Item deleted from table but it was not a TTL event")
			}
		}
	}
}

func main() {
	lambda.Start(handler)
}
