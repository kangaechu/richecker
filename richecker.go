package richecker

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"time"
)

func Check(days int) {
	sess := session.Must(session.NewSession())

	CheckEC2(sess, days)
}

func CheckEC2(sess *session.Session, days int) {
	fmt.Println("Days     :", days)

	var timeout time.Duration
	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials. A
	// Session should be shared where possible to take advantage of
	// configuration and credential caching. See the session package for
	// more information.

	// Create a new instance of the service's client with a Session.
	// Optional aws.Config values can also be provided as variadic arguments
	// to the New function. This option allows you to provide service
	// specific configuration.
	svc := ec2.New(sess)

	// Create a context with a timeout that will abort the upload if it takes
	// more than the passed in timeout.
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	// Ensure the context is canceled to prevent leaking.
	// See context package for more information, https://golang.org/pkg/context/
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
