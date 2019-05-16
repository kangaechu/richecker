package richecker

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/redshift"
	"log"
	"sort"
	"time"
)

type ReservedInstance struct {
	Service string
	Type    string
	Count   int64
	End     time.Time
}

type ReservedInstances []*ReservedInstance

// Sort Reserved Instance order end date (old is former).
func (ris ReservedInstances) Len() int {
	return len(ris)
}

func (ris ReservedInstances) Swap(i, j int) {
	ris[i], ris[j] = ris[j], ris[i]
}

func (ris ReservedInstances) Less(i, j int) bool {
	return ris[i].End.Before(ris[j].End)
}

func (ri ReservedInstance) Print() {
	localtime := ri.End.In(time.Local).Format("2006-01-02")
	fmt.Printf("Reserved Instance is almost expiring at %s %-11s %s x %d\n", localtime, ri.Service, ri.Type, ri.Count)
}

func Check(days int) {
	expireAt := time.Now().Add(time.Duration(days) * 24 * time.Hour)

	sess := session.Must(session.NewSession())

	ec2RIs := CheckEC2(sess, expireAt)
	for _, ri := range ec2RIs {
		ri.Print()
	}
	rdsRIs := CheckRDS(sess, expireAt)
	for _, ri := range rdsRIs {
		ri.Print()
	}
	elastiCacheRIs := CheckElastiCache(sess, expireAt)
	for _, ri := range elastiCacheRIs {
		ri.Print()
	}
	redshiftRIs := CheckRedshift(sess, expireAt)
	for _, ri := range redshiftRIs {
		ri.Print()
	}
}

func CheckEC2(sess *session.Session, expireAt time.Time) []*ReservedInstance {
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

	var RIs ReservedInstances
	for _, ri := range activeRIs.ReservedInstances {
		if ri.End.Before(expireAt) {
			RIs = append(RIs, &ReservedInstance{"EC2", *ri.InstanceType, *ri.InstanceCount, *ri.End})
		}
	}
	sort.Sort(RIs)
	return RIs
}

func CheckRDS(sess *session.Session, expireAt time.Time) []*ReservedInstance {
	svc := rds.New(sess)

	// RDS could not filter.
	// https://godoc.org/github.com/aws/aws-sdk-go/service/rds#DescribeReservedDBInstancesInput
	params := rds.DescribeReservedDBInstancesInput{}
	activeRIs, err := svc.DescribeReservedDBInstances(&params)
	if err != nil {
		log.Fatalln("cannot get RDS RI information.", err)
	}

	var RIs ReservedInstances
	for _, ri := range activeRIs.ReservedDBInstances {
		if *ri.State == "active" {
			end := ri.StartTime.Add(time.Duration(*ri.Duration) * time.Second)
			if end.Before(expireAt) {
				RIs = append(RIs, &ReservedInstance{"RDS", *ri.DBInstanceClass, *ri.DBInstanceCount, end})
			}
		}
	}
	sort.Sort(RIs)
	return RIs
}

func CheckElastiCache(sess *session.Session, expireAt time.Time) []*ReservedInstance {
	svc := elasticache.New(sess)

	params := elasticache.DescribeReservedCacheNodesInput{}
	activeRIs, err := svc.DescribeReservedCacheNodes(&params)
	if err != nil {
		log.Fatalln("cannot get ElastiCache RI information.", err)
	}

	var RIs ReservedInstances
	for _, ri := range activeRIs.ReservedCacheNodes {
		if *ri.State == "active" {
			end := ri.StartTime.Add(time.Duration(*ri.Duration) * time.Second)
			if end.Before(expireAt) {
				RIs = append(RIs, &ReservedInstance{"ElastiCache", *ri.CacheNodeType, *ri.CacheNodeCount, end})
			}
		}
	}
	sort.Sort(RIs)
	return RIs
}

func CheckRedshift(sess *session.Session, expireAt time.Time) []*ReservedInstance {
	svc := redshift.New(sess)

	params := redshift.DescribeReservedNodesInput{}
	activeRIs, err := svc.DescribeReservedNodes(&params)
	if err != nil {
		log.Fatalln("cannot get RDS RI information.", err)
	}

	var RIs ReservedInstances
	for _, ri := range activeRIs.ReservedNodes {
		if *ri.State == "active" {
			end := ri.StartTime.Add(time.Duration(*ri.Duration) * time.Second)
			if end.Before(expireAt) {
				RIs = append(RIs, &ReservedInstance{"Redshift", *ri.NodeType, *ri.NodeCount, end})
			}
		}
	}
	sort.Sort(RIs)
	return RIs
}
