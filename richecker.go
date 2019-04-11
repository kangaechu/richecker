package richecker

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/redshift"
	"log"
	"time"
)

func Check(days int) {
	fmt.Println("Days:", days)

	sess := session.Must(session.NewSession())

	var timeout time.Duration
	ctx := context.Background()

	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	CheckEC2(ctx, sess)
	CheckRDS(ctx, sess)
	CheckElastiCache(ctx, sess)
	CheckRedshift(ctx, sess)
	//CheckCloudFront(sess)
	//CheckDynamoDB(sess)
}

func CheckEC2(ctx context.Context, sess *session.Session) {
	svc := ec2.New(sess)

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

func CheckRDS(ctx context.Context, sess *session.Session) {
	svc := rds.New(sess)

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

func CheckElastiCache(ctx context.Context, sess *session.Session) {
	svc := elasticache.New(sess)

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

func CheckRedshift(ctx context.Context, sess *session.Session) {
	svc := redshift.New(sess)

	params := redshift.DescribeReservedNodesInput{}
	activeRIs, err := svc.DescribeReservedNodes(&params)
	if err != nil {
		log.Fatalln("cannot get RDS RI information.", err)
	}
	for _, ri := range activeRIs.ReservedNodes {
		if *ri.State == "active" {
			fmt.Println(*ri.NodeType)
			fmt.Println(*ri.NodeCount)
			fmt.Println(ri.StartTime.Add(time.Duration(*ri.Duration) * time.Second))
		}
	}
}

func CheckCloudFront(ctx context.Context, sess *session.Session) {
	svc := redshift.New(sess)

	params := redshift.DescribeReservedNodesInput{}
	activeRIs, err := svc.DescribeReservedNodes(&params)
	if err != nil {
		log.Fatalln("cannot get RDS RI information.", err)
	}
	for _, ri := range activeRIs.ReservedNodes {
		if *ri.State == "active" {
			fmt.Println(*ri.NodeType)
			fmt.Println(*ri.NodeCount)
			fmt.Println(ri.StartTime.Add(time.Duration(*ri.Duration) * time.Second))
		}
	}
}
