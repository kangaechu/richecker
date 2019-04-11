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

type ReservedInstance struct {
	Service string
	Type    string
	Count   int64
	End     time.Time
}

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

	ec2RIs := CheckEC2(ctx, sess)
	fmt.Println(&ec2RIs)
	rdsRIs := CheckRDS(ctx, sess)
	fmt.Println(&rdsRIs)
	elastiCacheRIs := CheckElastiCache(ctx, sess)
	fmt.Println(&elastiCacheRIs)
	redshiftRIs := CheckRedshift(ctx, sess)
	fmt.Println(&redshiftRIs)

}

func CheckEC2(ctx context.Context, sess *session.Session) []*ReservedInstance {
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

	var RIs []*ReservedInstance
	for _, ri := range activeRIs.ReservedInstances {
		RIs = append(RIs, &ReservedInstance{"EC2", *ri.InstanceType, *ri.InstanceCount, *ri.End})
	}
	return RIs
}

func CheckRDS(ctx context.Context, sess *session.Session) []*ReservedInstance {
	svc := rds.New(sess)

	// RDS could not filter.
	// https://godoc.org/github.com/aws/aws-sdk-go/service/rds#DescribeReservedDBInstancesInput
	params := rds.DescribeReservedDBInstancesInput{}
	activeRIs, err := svc.DescribeReservedDBInstances(&params)
	if err != nil {
		log.Fatalln("cannot get RDS RI information.", err)
	}

	var RIs []*ReservedInstance
	for _, ri := range activeRIs.ReservedDBInstances {
		if *ri.State == "active" {
			end := ri.StartTime.Add(time.Duration(*ri.Duration) * time.Second)
			RIs = append(RIs, &ReservedInstance{"RDS", *ri.DBInstanceClass, *ri.DBInstanceCount, end})
		}
	}
	return RIs
}

func CheckElastiCache(ctx context.Context, sess *session.Session) []*ReservedInstance {
	svc := elasticache.New(sess)

	params := elasticache.DescribeReservedCacheNodesInput{}
	activeRIs, err := svc.DescribeReservedCacheNodes(&params)
	if err != nil {
		log.Fatalln("cannot get ElastiCache RI information.", err)
	}

	var RIs []*ReservedInstance
	for _, ri := range activeRIs.ReservedCacheNodes {
		if *ri.State == "active" {
			end := ri.StartTime.Add(time.Duration(*ri.Duration) * time.Second)
			RIs = append(RIs, &ReservedInstance{"RDS", *ri.OfferingType, *ri.CacheNodeCount, end})
		}
	}
	return RIs
}

func CheckRedshift(ctx context.Context, sess *session.Session) []*ReservedInstance {
	svc := redshift.New(sess)

	params := redshift.DescribeReservedNodesInput{}
	activeRIs, err := svc.DescribeReservedNodes(&params)
	if err != nil {
		log.Fatalln("cannot get RDS RI information.", err)
	}

	var RIs []*ReservedInstance
	for _, ri := range activeRIs.ReservedNodes {
		if *ri.State == "active" {
			end := ri.StartTime.Add(time.Duration(*ri.Duration) * time.Second)
			RIs = append(RIs, &ReservedInstance{"RDS", *ri.NodeType, *ri.NodeCount, end})
		}
	}
	return RIs
}
