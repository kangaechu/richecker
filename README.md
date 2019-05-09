# AWS Reserved Instance Checker (richecker)

richecker is a tool for checking expiration of AWS Reserved Instance.

## Getting Started

### Installing

Download binary from releases page.

#### IAM Profile

Before you use this tool, attach following profile into your resource (EC2, Lambda, etc.).

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "RIChecker",
            "Effect": "Allow",
            "Action": [
                "rds:DescribeReservedDBInstances",
                "redshift:DescribeReservedNodes",
                "elasticache:DescribeReservedCacheNodes",
                "ec2:DescribeReservedInstances"
            ],
            "Resource": "*"
        }
    ]
}
```

### Run 

```
$ richecker -d <DAYS_BEFORE_EXPIREATION> 
```

## License

MIT


