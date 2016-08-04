# Prometheus AWS audit exporter

This program is intended to export various AWS statistics as prometheus
metrics. It is primarily intended to assist with billing. Currently the
following metrics are exported:

# EC2 Instance Counts

 - *aws_ec2_instances_count*: Count of istances

The following labels are exposed:

 - *az*: availability zone
 - *instance_type*: type of instance
 - *groups*: sorted comma separated list of groups.
 - *owner_id*: The owner id
 - *requester_id*: The requester id (default to owner id if none is present)

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
 - *tenancy*:
 - *offer_type*:
 - *product*:

