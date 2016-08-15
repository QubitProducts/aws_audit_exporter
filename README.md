# Prometheus AWS audit exporter

This program is intended to export various AWS statistics as prometheus
metrics. It is primarily intended to assist with billing. Currently the
following metrics are exported:

# EC2 Instance Counts

 - *aws_ec2_instances_count*: Count of istances

The following labels are exposed:

 - *az*: availability zone
 - *instance_type*: type of instance
 - *lifecycle*: spot or scheduled instance
 - *groups*: sorted comma separated list of groups.
 - *owner_id*: The owner id
 - *requester_id*: The requester id (default to owner id if none is present)
 - *aws_tag_*: Any tags passed in with the -instance-tags flag are added as labels

# EC2 Reserved Instances
Every set of instance reservations gets its own time series, this is intended to allow
the end time of reserved intances to be tracked and potentially alerted upon.

 - *aws_ec2_reserved_instances_usage_price_dollars*: cost of reserved instance usage in dollars
 - *aws_ec2_reserved_instances_fixed_price_dollars*: fixed cost of reserved instance in dollars
 - *aws_ec2_reserved_instances_price_per_hour_dollars*: hourly cost of reserved instance in dollars
 - *aws_ec2_reserved_instances_count*: Number of reserved instances in this reservation
 - *aws_ec2_reserved_instances_start_time*: Start time of this reservation
 - *aws_ec2_reserved_instances_end_time*: End time of this reservation

The following labels are exposed:

 - *id*: the reservation id
 - *az*: availability zone
 - *instance_type*: type of instance
 - *product*: The product description
 - *tenancy*:
 - *offer_type*:

# EC2 Spot Instance Request

Only fullfilled active spot instances requests are currently tracke

 - *aws\_ec2\_spot\_request\_count*: How active spot instances of a given type you have running
 - *aws\_ec2\_spot\_request\_bid\_price\_hourly\_dollars*: Your maximum bid price
 - *aws\_ec2\_spot\_request\_actual\_block\_price\_hourly\_dollars*: The price paid for limited duration spot instances

The following labels are exposed:

 - *az*: availability zone
 - *instance_type*: type of instance
 - *product*: The product description

# EC2 Spot Instance Pricing

Only prices for products that have been seen in spot instance requests are tracked.

 - *aws\_ec2\_spot\_price\_per\_hour\_dollars*: The current market price for a spot instance

The following labels are exposed:

 - *az*: availability zone
 - *instance_type*: type of instance
 - *product*: The product description

# Usage

  Your aws credentials should either be in $HOME/.aws/credentials , or set via AWS\_ACCESS\_KEY and AWS\_SECRET\_ACCESS\_KEY

  Usage of /go/bin/aws_audit_exporter:
  -addr string
        port to listen on (default ":9190")
  -duration duration
        How often to query the API (default 4m0s)
  -instance-tags string
        comma seperated list of tag keys to use as metric labels
  -region string
        the region to query (default "eu-west-1")

# TODO

 - Add optional Push gateway support
 - Make tracking full Reserved instance tracking optional and pre-aggregate

