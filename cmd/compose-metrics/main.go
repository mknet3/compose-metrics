package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"log"
	"math"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	cpuUtilization := getCPUUtilization()
	messagesQueued := getMessagesQueued()
	value := 0
	if math.Round(cpuUtilization) < 80 {
		value = messagesQueued
	}

	json.NewEncoder(w).Encode(Metrics{
		Value: value,
	})
}

func handleRequests() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}

type Metrics struct {
	Value int `json:"value"`
}

func getCPUUtilization() float64 {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("REGION"))
	if err != nil {
		panic("failed to load config, " + err.Error())
	}

	ctx := context.Background()

	svc := cloudwatch.NewFromConfig(cfg)

	result, err := svc.GetMetricData(ctx, &cloudwatch.GetMetricDataInput{
		EndTime:   aws.Time(time.Unix(time.Now().Unix(), 0)),
		StartTime: aws.Time(time.Unix(time.Now().Add(time.Duration(-60)*time.Minute).Unix(), 0)),
		MetricDataQueries: []types.MetricDataQuery{
			{
				Id: aws.String("q1"),
				MetricStat: &types.MetricStat{
					Metric: &types.Metric{
						Namespace:  aws.String("AWS/RDS"),
						MetricName: aws.String("CPUUtilization"),
						Dimensions: []types.Dimension{
							types.Dimension{
								Name:  aws.String("DBInstanceIdentifier"),
								Value: aws.String("RDS SERVER NAME"),
							},
						},
					},
					Period: aws.Int32(60),
					Stat:   aws.String("Average"),
				},
			},
		},
	})

	if err != nil {
		fmt.Printf("error, %v", err)
		return 0
	}

	for _, dataResult := range result.MetricDataResults {
		for _, value := range dataResult.Values {
			fmt.Printf("CPU Utilization: %v", value)
			return value
		}
	}

	fmt.Println("There are not metric results")
	return 0
}

func getMessagesQueued() int {
	return 2000
}
