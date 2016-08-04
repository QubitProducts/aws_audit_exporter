// Copyright 2016 Qubit Group
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/tcolgate/client_golang/prometheus"
)

var (
	region  = flag.String("region", "eu-west-1", "the region to query")
	taglist = flag.String("instance-tags", "", "comma seperated list of tag keys to use as metric labels")
	dur     = flag.Duration("duration", time.Minute*4, "How often to query the API")
	addr    = flag.String("addr", ":9190", "port to listen on")

	riLabels = []string{
		"az",
		"tenancy",
		"instance_type",
		"offer_type",
		"product",
		"id",
	}
	riUsagePrice = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_ec2_reserved_instances_usage_price_dollars",
		Help: "cost of reserved instance usage in dollars",
	},
		riLabels)
	riFixedPrice = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_ec2_reserved_instances_fixed_price_dollars",
		Help: "fixed cost of reserved instance in dollars",
	},
		riLabels)
	riHourlyPrice = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_ec2_reserved_instances_price_per_hour_dollars",
		Help: "hourly cost of reserved instance in dollars",
	},
		riLabels)
	riInstanceCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_ec2_reserved_instances_count",
		Help: "Number of reserved instances in this reservation",
	},
		riLabels)
	riStartTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_ec2_reserved_instances_start_time",
		Help: "Start time of this reservation",
	},
		riLabels)
	riEndTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_ec2_reserved_instances_end_time",
		Help: "End time of this reservation",
	},
		riLabels)

	instancesLabels = []string{
		"groups",
		"owner_id",
		"requester_id",
		"az",
		"instance_type",
	}
)

// We have to construct the set of tags for this based on the program
// args, so it is created in main
var instancesCount *prometheus.GaugeVec
var instanceTags = map[string]string{}

func main() {
	flag.Parse()

	tagl := []string{}
	for _, tstr := range strings.Split(*taglist, ",") {
		ctag := tagname(tstr)
		instanceTags[tstr] = ctag
		tagl = append(tagl, ctag)
	}
	instancesCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_ec2_instances_count",
		Help: "End time of this reservation",
	},
		append(instancesLabels, tagl...))

	prometheus.Register(instancesCount)
	prometheus.Register(riUsagePrice)
	prometheus.Register(riFixedPrice)
	prometheus.Register(riHourlyPrice)
	prometheus.Register(riInstanceCount)
	prometheus.Register(riStartTime)
	prometheus.Register(riEndTime)

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("failed to create session %v\n", err)
	}

	svc := ec2.New(sess, &aws.Config{Region: aws.String(*region)})

	go func() {
		for {
			go instances(svc, *region)
			go reservations(svc, *region)
			<-time.After(*dur)
		}
	}()

	http.Handle("/metrics", prometheus.Handler())

	log.Println(http.ListenAndServe(*addr, nil))
}

func instances(svc *ec2.EC2, awsRegion string) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-code"),
				Values: []*string{aws.String("16")},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", awsRegion, err.Error())
		log.Fatal(err.Error())
	}

	instancesCount.Reset()
	labels := prometheus.Labels{}
	for _, r := range resp.Reservations {
		groups := []string{}
		for _, g := range r.Groups {
			groups = append(groups, *g.GroupName)
		}
		sort.Strings(groups)
		labels["groups"] = strings.Join(groups, ",")
		labels["owner_id"] = *r.OwnerId
		labels["requester_id"] = *r.OwnerId
		if r.RequesterId != nil {
			labels["requester_id"] = *r.RequesterId
		}
		for _, ins := range r.Instances {
			labels["az"] = *ins.Placement.AvailabilityZone
			labels["instance_type"] = *ins.InstanceType
			for _, label := range instanceTags {
				labels[label] = ""
			}
			for _, tag := range ins.Tags {
				label, ok := instanceTags[*tag.Key]
				if ok {
					labels[label] = *tag.Value
				}
			}
			instancesCount.With(labels).Inc()
		}
	}
}

func reservations(svc *ec2.EC2, awsRegion string) {
	params := &ec2.DescribeReservedInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("state"),
				Values: []*string{aws.String("active")},
			},
		},
	}
	resp, err := svc.DescribeReservedInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", awsRegion, err.Error())
		log.Fatal(err.Error())
	}

	labels := prometheus.Labels{}
	for _, r := range resp.ReservedInstances {
		labels["az"] = *r.AvailabilityZone
		labels["instance_type"] = *r.InstanceType
		labels["tenancy"] = *r.InstanceTenancy
		labels["offer_type"] = *r.OfferingType
		labels["product"] = *r.ProductDescription
		labels["id"] = *r.ReservedInstancesId

		riUsagePrice.With(labels).Set(*r.UsagePrice)
		riFixedPrice.With(labels).Set(*r.FixedPrice)
		riHourlyPrice.With(labels).Set(0)
		for _, c := range r.RecurringCharges {
			if *c.Frequency == "Hourly" {
				riHourlyPrice.With(labels).Set(*c.Amount)
			}
		}
		riInstanceCount.With(labels).Set(float64(*r.InstanceCount))
		riStartTime.With(labels).Set(float64(r.Start.Unix()))
		riEndTime.With(labels).Set(float64(r.End.Unix()))

	}
}

var cleanre = regexp.MustCompile("[^A-Za-z0-9]")

func tagname(n string) string {
	c := cleanre.ReplaceAllString(n, "_")
	c = strings.ToLower(strings.Trim(c, "_"))
	return "aws_tag_" + c
}
