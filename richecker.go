package richecker

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/rds"
	"log"
	"time"
)

func Check(days int) {
	fmt.Println("Days     :", days)

	sess := session.Must(session.NewSession())

	CheckEC2(sess)
	CheckRDS(sess)
	CheckElastiCache(sess)
	//CheckRedshift(sess)
	//CheckCloudFront(sess)
	//CheckDynamoDB(sess)
}

func CheckEC2(sess *session.Session) {

	var timeout time.Duration
	svc := ec2.New(sess)

	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	// filter state=active
	params := ec2.DescribeReservedInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("state"),
				Values: []*string{aws.String("active")},
			},
		},
	}
	activeRIs, err := svc.DescribeReservedInstances(&params)
	if err != nil {
		log.Fatalln("cannot get EC2 RI information.", err)
	}

	for _, ri := range activeRIs.ReservedInstances {
		fmt.Println(*ri.InstanceType)
		fmt.Println(*ri.InstanceCount)
		fmt.Println(ri.End)
	}
}

func CheckRDS(sess *session.Session) {
	var timeout time.Duration

	svc := rds.New(sess)

	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	// RDS could not filter.
	// https://godoc.org/github.com/aws/aws-sdk-go/service/rds#DescribeReservedDBInstancesInput
	params := rds.DescribeReservedDBInstancesInput{}
	activeRIs, err := svc.DescribeReservedDBInstances(&params)
	if err != nil {
		log.Fatalln("cannot get RDS RI information.", err)
	}
	for _, ri := range activeRIs.ReservedDBInstances {
		if *ri.State == "active" {
			fmt.Println(*ri.DBInstanceClass)
			fmt.Println(*ri.DBInstanceCount)
			fmt.Println(ri.StartTime.Add(time.Duration(*ri.Duration) * time.Second))
		}
	}
}

func CheckElastiCache(sess *session.Session) {
	var timeout time.Duration

	svc := elasticache.New(sess)

	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	params := elasticache.DescribeReservedCacheNodesInput{}
	activeRIs, err := svc.DescribeReservedCacheNodes(&params)
	if err != nil {
		log.Fatalln("cannot get ElastiCache RI information.", err)
	}
	for _, ri := range activeRIs.ReservedCacheNodes {
		if *ri.State == "active" {
			fmt.Println(ri)
			fmt.Println(*ri.OfferingType)
			fmt.Println(*ri.CacheNodeCount)
			fmt.Println(ri.StartTime.Add(time.Duration(*ri.Duration) * time.Second))
		}
	}
}
