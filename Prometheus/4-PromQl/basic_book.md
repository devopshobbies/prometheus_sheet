promQL is the Prometheus Query Language. While it ends in QL, you will find that it is not an SQL-like language, as SQL languages tend to lack expressive power when it comes to the sort of calculations you would like to perform on time series. Labels are a key part of PromQL, and you can use them not only to do arbitrary aggregations but also to join different metrics together for arithmetic operations against them. There are a wide variety of functions available to you from prediction to date and math functions.
This chapter will introduce you to the basic concepts of PromQL, including aggrega‐ tion, basic types, and the HTTP API.


### Aggregation Basics

Let’s get started with some simple aggregation queries. These queries will likely cover most of your potential uses for PromQL. While PromQL is as powerful as it is possible to be,1 most of the time your needs will be reasonably simple.

## Gauge
Gauges are a snapshot of state, and usually when aggregating them you want to take a sum, average, minimum, or maximum.

Consider the metric node_filesystem_size_bytes from your Node Exporter, which reports the size of each of your mounted filesystems, and has device, fstype, and mountpoint labels. You can calculate total filesystem size on each machine with:

`sum without(device, fstype, mountpoint)(node_filesystem_size_bytes)`

This works as without tells the sum aggregator to sum everything up with the same labels, ignoring those three. So if you had the time series:
```
node_filesystem_free_bytes{device="/dev/sda1",fstype="vfat",
instance="localhost:9100",job="node",mountpoint="/boot/efi"} 70300672
node_filesystem_free_bytes{device="/dev/sda5",fstype="ext4",
instance="localhost:9100",job="node",mountpoint="/"} 30791843840
node_filesystem_free_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run"} 817094656
node_filesystem_free_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run/lock"} 5238784
node_filesystem_free_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run/user/1000"} 826912768

```

the result would be:

`{instance="localhost:9100",job="node"} 32511390720`

You will notice that the device, fstype, and mountpoint labels are now gone. The metric name is also no longer present, as this is no longer node_filesys tem_free_bytes because math has been performed on it. Since there is only one Node Exporter being scraped by Prometheus, there is only one result, but if you were scraping more, then you would have a result for each of the Node Exporters.

You could go a step further and remove the instance label with:

`sum without(device, fstype, mountpoint, instance)(node_filesystem_size_bytes)`

This as expected removes the instance label, but the value remains the same as the previous expression because there is only one Node Exporter to aggregate metrics from:

`{job="node"} 32511390720`

You can use the same approach with other aggregations. max would tell you the size of the biggest mounted filesystem on each machine:

`max without(device, fstype, mountpoint)(node_filesystem_size_bytes)`

The outputted labels are exactly the same as when you aggregated using sum:

`{instance="localhost:9100",job="node"} 30792601600`

This predictability in what labels are returned is important for vector matching with
operators, as will be discussed in Chapter 15.

You are not limited to aggregating metrics about one type of job. For example, to find
the average number of file descriptors open across all your jobs, you could use:


`avg without(instance, job)(process_open_fds)`


## Counter
Counters track the number or size of events, and the value your applications expose on their /metrics is the total since it started. But that total is of little use to you on its own; what you really want to know is how quickly the counter is increasing over time. This is usually done using the rate function, though the increase and irate functions also operate on counter values.

For example, to calculate the amount of network traffic received per second, you could use:

`rate(node_network_receive_bytes_total[5m])`

The [5m] says to provide rate with 5 minutes of data, so the returned value will be an
average over the last 5 minutes:

```
{device="lo",instance="localhost:9100",job="node"} 1859.389655172414
{device="wlan0",instance="localhost:9100",job="node"} 1314.5034482758622

```


The values here are not integers, as the 5-minute window rate is looking at does not perfectly align with the samples that Prometheus has scraped. Some estimation is used to fill in the gaps between the data points you have and the boundaries of the range.

The output of rate is a gauge, so the same aggregations apply as for gauges. The node_network_receive_bytes_total metric has a device label, so if you aggregate it away you will get the total bytes received per machine per second:


`sum without(device)(rate(node_network_receive_bytes_total[5m]))`


Running this query will give you a result like:

`{instance="localhost:9100",job="node"} 3173.8931034482762`

You can filter down which time series to request, so you could only look at eth0 and then aggregate it across all machines by aggregating away the instance label:

`sum without(instance)(rate(node_network_receive_bytes_total{device="eth0"}[5m]))`

When you run this query the instance label is gone, but the device label remains as you did not ask for it to be removed:

`{device="eth0",job="node"} 3173.8931034482762`


There is no ordering or hierarchy within labels, allowing you to aggregate by as many or as few labels as you like.


## Summary

A summary metric will usually contain both a _sum and _count, and sometimes a time series with no suffix with a quantile label. The _sum and _count are both counters.
Your Prometheus exposes an http_response_size_bytes summary for the amount of data some of its HTTP APIs return.2 http_response_size_bytes_count tracks the number of requests, and as it is a counter, you must use rate before aggregating away its handler label:

`sum without(handler)(rate(http_response_size_bytes_count[5m]))`

This gives you the total per-second HTTP request rate, and as the Node Exporter also
returns this metric, you will see both jobs in the result:

```
{instance="localhost:9090",job="prometheus"} 0.26868836781609196
{instance="localhost:9100",job="node"} 0.1

```


Similarly, http_response_size_bytes_sum is a counter with the number of bytes
each handle has returned, so the same pattern applies:


`sum without(handler)(rate(http_response_size_bytes_sum[5m]))`

This will return results with the same labels as the previous query, but the values are
larger as responses tend to return many bytes:

```
{instance="localhost:9090",job="prometheus"} 796.0015958275862
{instance="localhost:9100",job="node"} 1581.6103448275862

```

The power of a summary is that it allows you to calculate the average size of an event, in this case the average amount of bytes that are being returned in each response. If you had three responses of size 1, 4, and 7, then the average would be their sum divided by their count, which is to say 12 divided by 3. The same applies to the summary. You divide the _sum by the _count (after taking a rate) to get an average over a time period:


```
sum without(handler)(rate(http_response_size_bytes_sum[5m]))
/
sum without(handler)(rate(http_response_size_bytes_count[5m]))

```

The division operator matches the time series with the same labels, and divides, giving you the same two time series out but with the average response size over the past 5 minutes as a value:

```
{instance="localhost:9090",job="prometheus"} 2962.54580091246150133317
{instance="localhost:9100",job="node"} 15816.10344827586200000000

```

When calculating an average, it is important that you first aggregate up the sum and count, and only as the last step perform the division. Otherwise, you could end up averaging averages, which is not statistically valid.

For example, if you wanted to get the average response size across all instances of a
job, you could do:

```
sum without(instance)(
sum without(handler)(rate(http_response_size_bytes_sum[5m]))
)
/
sum without(instance)(
sum without(handler)(rate(http_response_size_bytes_count[5m]))
)

```

However, it’d be incorrect to do:


```
avg without(instance)(
sum without(handler)(rate(http_response_size_bytes_sum[5m]))
/
sum without(handler)(rate(http_response_size_bytes_count[5m]))
)

```

It is incorrect to average an average, and both the division and avg would be calculat‐
ing averages.


It is not possible for you to aggregate the quantiles of a summary
(the time series with the quantile label) from a statistical stand‐
point.


### Histogram


Histogram metrics allow you to track the distribution of the size of events, allowing you to calculate quantiles from them. For example, you can use histograms to calcu‐ late the 0.9 quantile (which is also known as the 90th percentile) latency. Prometheus 2.37.1 exposes a histogram metric called prometheus_tsdb_compac tion_duration_seconds that tracks how many seconds compaction takes for the time series database. This histogram metric has time series with a _bucket suffix called prometheus_tsdb_compaction_duration_seconds_bucket. Each bucket has a le label, which is a counter of how many events have a size less than or equal to the bucket boundary. This is an implementation detail you largely need not worry about as the histogram_quantile function takes care of this when calculating quantiles. For example, the 0.90 quantile would be:


```
histogram_quantile(
0.90,
rate(prometheus_tsdb_compaction_duration_seconds_bucket[1d]))

```

As prometheus_tsdb_compaction_duration_seconds_bucket is a counter, you must first take a rate. Compaction usually only happens every two hours, so a one-day time range is used here and you will see a result in the expression browser such as:


`{instance="localhost:9090",job="prometheus"} 7.720000000000001`


This indicates that the 90th percentile latency of compactions is around 7.72 seconds. As there will usually only be 12 compactions in a day, the 90th percentile says that 10% of compactions take longer than this, which is to say one or two compactions. This is something to be aware of when using quantiles. For example, if you want to calculate a 0.999 quantile, you should have several thousand data points to work with in order to produce a reasonably accurate answer. If you have fewer than that, single outliers could greatly affect the result, and you should consider using lower quantiles to avoid making statements about your system for which you have insufficient data to back up.
Usually you would use a 5- or 10-minute rate with histograms. All the bucket time series combined with any labels, and a long range on the rate, can make for a lot of samples that need to be processed. Be wary of PromQL expressions using ranges that are hours or days, as they can be relatively expensive to calculate.


Similar to when taking averages, using histogram_quantile should be the last step in a query expression. Quantiles cannot be aggregated, or have arithmetic performed upon them, from a statistical standpoint. Accordingly, when you want to take a histo‐ gram of an aggregate, first aggregate up with sum and then use histogram_quantile:

```
histogram_quantile(
0.90,
sum without(instance)(rate(prometheus_tsdb_compaction_duration_bucket[1d])))

```

This calculates the 0.9 quantile compaction duration across all of your Prometheus servers, and will produce a result without an instance label:


`{job="prometheus"} 7.720000000000001`


Histogram metrics also include _sum and _count metrics, which work exactly the same as for the summary metric. You can use these to calculate average event sizes, such as the average compaction duration:

```
sum without(instance)(rate(prometheus_tsdb_compaction_duration_sum[1d]))
/
sum without(instance)(rate(prometheus_tsdb_compaction_duration_count[1d]))
```

This would produce a result like:


`{job="prometheus"} 3.1766430400714287`

### Selectors


Working with all the different time series with different label values for a metric can be a bit overwhelming, and potentially confusing if a metric is coming from multiple different types of servers.5 Usually you will want to narrow down which time series you are working on. You almost always will want to limit by job label, and depending on what you are up to, you might want to only look at one instance or one handler, for example.

This limiting by labels is done using selectors. You have seen selectors in every example thus far, and now we are going to explain them to you in detail. For example:


`process_resident_memory_bytes{job="node"}`

is a selector that will return all time series with the name process_resident_ memory_bytes and a job label of node. This particular selector is most properly called an instant vector selector, as it returns the values of the given time series at a given instant. Vector here basically means a one-dimensional list, as a selector can return zero or more time series, and each time series will have one sample.

The job="node" is called a matcher, and you can have many matchers in one selector that are ANDed together.

## Matchers

There are four matchers (you have already seen the equality matcher, which is also the most commonly used):

=
This is the equality matcher; for example, job="node". With this you can specify that the returned time series has a label name with exactly the given label value. As an empty label, value is the same as not having that label, so you could use foo="" to specify that the foo label not be present.

!=
This is the negative equality matcher; for example, job!="node". With this you can specify that the returned time series do not have a label name with exactly the given label value.
=~
This is the regular expression matcher; for example, job=~"n.*". With this you specify that for the returned time series, the given label’s value will be matched by the regular expression. The regular expression is fully anchored, which is to say that the regular expression a will only match the string a, and not xa or ax. You can prepend or suffix your regular expression with .* if you do not want this behavior.6 As with relabeling, the RE2 regular expression engine is used, as covered in “Regular Expressions” on page 152.

!~
This is the negative regular expression matcher. RE2 does not support negative lookahead expressions, so this provides you with an alternative way to exclude label values based on a regular expression. You can have multiple matchers with the same label name in a selector, which can be a substitute for negative lookahead expressions. For example, to find the size of all filesystems mounted under /run but not /run/user, you could use:
`node_filesystem_size_bytes{job="node",mountpoint=~"/run/.*",mountpoint!~"/run/user/.*"}`

Internally, the metric name is stored in a label called __name__ (as discussed in “Reserved Labels and __name__” on page 90), so process_resident_ memory_bytes{job="node"} is syntactic sugar for {name="process_resident_ memory_bytes",job="node"}. You can even do regular expressions on the metric name, but this is unwise outside of when you are debugging the performance of the Prometheus server.


Having to use regular expression matchers is a little bit of a smell. If you find yourself using them a lot on a given label, consider if you should instead combine the matched label values into one. For example, for HTTP status codes instead of doing code~="4.." to catch 401s, 404s, 405s, etc., you might combine them into a label value 4xx and use the equality matcher code="4xx".

The selector {} returns an error, which is a safety measure to avoid accidentally returning all the time series inside the Prometheus server as that could be expen‐ sive. To be more precise, at least one of the matchers in a selector must not match the empty string. So {foo=""} and {foo=~".*"} will return an error, while {foo="",bar="x"}, {foo!=""}, or {foo=~".+"} are permitted.

## Instant Vector

An instant vector selector returns an instant vector of the most recent samples before the query evaluation time, which is to say a list of zero or more time series. Each of these time series will have one sample, and a sample contains both a value and a timestamp. While the instant vector returned by an instant vector selector has the timestamp of the original data,9 any instant vectors returned by other operations or functions will have the timestamp of the query evaluation time for all of their values.

When you ask for current memory usage, you do not want samples from an instance that was turned down days ago to be included, a concept known as staleness. In Prometheus 1.x this was handled by returning time series that had a sample no more than 5 minutes before the query evaluation time. This largely worked but had downsides such as double counting if an instance restarted with a new instance label within that 5-minute window.

Prometheus 2.x has a more sophisticated approach. If a time series disappears from one scrape to the next, or if a target is no longer returned from service discovery, a special type of sample called a stale marker10 is appended to the time series. When evaluating an instant vector selector, all time series satisfying all the matchers are first found, and the most recent sample in the 5 minutes before the query evaluation time is still considered. If the sample is a normal sample, then it is returned in the instant vector, but if it is a stale marker, then that time series will not be included in that instant vector. The outcome of all of this is that when you use an instant vector selector, time series that have gone stale are not returned.

If you have an exporter exposing timestamps, as described in “Timestamps” on page 82, then stale markers and the Prometheus 2.x staleness logic will not apply. The affected time series will work instead with the older logic that looks back 5 minutes.
### Range Vector

There is a second type of selector you have already seen, called the range vector selector. Unlike an instant vector selector, which returns one sample per time series, a range vector selector can return many samples for each time series.11 Range vectors are always used with the rate function, for example:

`rate(process_cpu_seconds_total[1m])`

The [1m] turns the instant vector selector into a range vector selector, and instructs PromQL to return for all time series matching the selector all samples for the minute up to the query evaluation time. If you execute just process_cpu_seconds_ total[1m] in the Console tab of the expression browser, you will see something like Figure 13-1.

In this case, each time series happens to have six samples in the past minute. You will notice that while the samples for each time series happen to be perfectly 10 seconds apart12 in line with the scrape interval you configured, the two time series timestamps are not aligned with each other. One time series has a sample with a timestamp of 1517925155.087 and the other 1517925156.245.

Figure 13-1.

This is because range vectors preserve the actual timestamps of the samples, and the scrapes for different targets are distributed in order to spread load more evenly. While you can control the frequency of scrapes and rule evaluations, you cannot con‐ trol their phase or alignment. If you have a 10-second scrape interval and hundreds of targets, then all those targets will be scraped at different points in a given 10-second window. Put another way, your time series all have slightly different ages. This generally won’t matter to you in practice, but can lead to artifacts as fundamentally metrics-based monitoring systems like Prometheus produce (quite good) estimates rather than exact answers.

You will very rarely look at range vectors directly. It only comes up when you need to see raw samples when debugging. Almost always you will use a range vector with a function such as rate or avg_over_time that takes a range vector as an argument. Staleness and stale markers have no impact on range vectors; you will get all the normal samples in a given range. Any stale markers also in that range are not returned by a range vector selector.


## Durations
Durations in Prometheus as used in PromQL and the configuration file support several units. You have already seen m for minute.
#### fig table


### Subqueries

While range vectors act on time series, they cannot be used in combination with functions. If you want to combine max_over_time with rate, you can either use recording rules, which would record the result of the rate function and pass it to the vector function, or you can use a subquery. A subquery is a part of a query that allows you to do a range query within a query. The syntax for a subquery uses square brackets, like range selectors. But it takes two different durations: the range and the resolution. The range is the range returned by the subquery, and the resolution acts as a step:

`max_over_time( rate(http_requests_total[5m])[30m:1m] )`

The preceding query runs rate(http_requests_total[5m]) every minute (1m) for the last 30 minutes (30m), then feeds the result in a max_over_time() function. The resolution can be omitted, such as in [30m:]. In this case, the global evaluation interval is used as resolution.

### Offset

There is a modifier you can use with either type of vector selector called offset. offset allows you to take the evaluation time for a query, and on a per-selector basis put it further back in time. For example:

`process_resident_memory_bytes{job="node"} offset 1h`
would get memory usage an hour before the query evaluation time. offset is not used much in simple queries like this, as it would be easier to change the evaluation time for the whole query instead. Where this can be useful is when you only want to adjust one selector in a query expression. For example:

`process_resident_memory_bytes{job="node"}`

-
`process_resident_memory_bytes{job="node"} offset 1h`
would give the change in memory usage in the Node Exporter over the past hour.13 The same approach works with range vectors:
`rate(process_cpu_seconds_total{job="node"}[5m])`
-
`rate(process_cpu_seconds_total{job="node"}[5m] offset 1h)`
offset allows you to look further back into the past, but also in the future, using a negative offset. This can be used when doing prediction or when the sample of the metrics is unaligned with the reality:
`rate(process_cpu_seconds_total{job="node"}[5m]) offset -1h`
-
`rate(process_cpu_seconds_total{job="node"}[5m])`
Note that this query will likely not return anything for the last hour. Grafana has a feature to shift in time a panel to a different time range than the rest of the dashboard it is a part of. In Grafana 5.0.0 you can find this in the Time range tab of the panel editor.


### At Modifier

Similar to the offset modifier, PromQL supports an @ modifier that lets you change the evaluation of vector selectors, range selectors, and subqueries to a fixed revaluation time. The @ modifier can be used with a Unix timestamp. The query http_requests_total @ 1667491200 returns the value of http_requests_total at 2022-11-03T16:00:00+00:00. The query rate(http_requests_total[5m] @1667491200) returns the 5-minute rate of http_requests_total at the same time. Additionally, start() and end() can be used as values for the @ modifier. For a range query, they resolve respectively with the start and the end of the range query. For an instant query, they both resolve to the evaluation time. In practice, it is possible to use the @ modifier to graph the evolution of the http_request_total that has a high rate at the end of the evaluation interval:


```
rate(http_requests_total[1m])
and
topk(5, rate(http_requests_total[1h] @ end()))
```



The topk(5, rate(http_requests_total[1h] @ end())) acts as a ranking function, filtering only the higher values at the end of the evaluation interval.

### HTTP API

Prometheus offers a number of HTTP APIs. The ones you will mostly interact with are query and query_range, which give you access to PromQL and can be used by dashboarding tools or custom reporting scripts. All the endpoints of interest are under /api/v1/, and beyond executing PromQL you can also look up time series metadata and perform administrative actions, such as taking snapshots and deleting time series. These other APIs are mainly of interest to dashboarding tools such as Grafana, which can use metadata to enhance its UI, and to those administering Prometheus, but are not relevant to PromQL execution.

### query

The query endpoint, or more formally /api/v1/query, executes a PromQL expression at a given time and returns the result. For example, http://localhost:9090/api/v1/query? query=process_resident_memory_bytes will return results like:

```
{
"status": "success",
"data": {
"resultType": "vector",
"result": [
{
"metric": {
"__name__": "process_resident_memory_bytes",
"instance": "localhost:9090",
"job": "prometheus"
},
"value": [1517929228.782, "91656192"]
},
{
"metric": {
"__name__": "process_resident_memory_bytes",
"instance": "localhost:9100",
"job": "node"
},
"value": [1517929228.782, "15507456"]
}
]
}
}

```


The status is success, meaning that the query worked. If it had failed, the status would be error, and an error field would provide more details. This particular result is an instant vector, which you can tell from "resultType": "vector". For each of the samples in the result, the labels are in the metric map, and the sample value is in the value list. The first number in the value list is the timestamp of the sample, in seconds, and the second is the actual value of the sample. The value is inside a string, as JSON cannot represent nonreal values such as NaN and +Inf. The time of all the samples will be the query evaluation time, even if the expression consisted of only an instant vector selector. Here the query evalua‐ tion time defaulted to the current time, but you can specify a time with the time URL parameter, which can be a Unix time, in seconds, or an RFC 3339 time. For example, http://localhost:9090/api/v1/query?query=process_resident_mem‐ ory_bytes&time=1514764800 would evaluate the query at midnight of January 1st, 2018. You can also use range vectors with the query endpoint. For example, http://local‐ host:9090/api/v1/query?query=prometheus_tsdb_head_samples_appended_total[1m] will return results like:

```

{
"status": "success",
"data": {
"resultType": "matrix",
"result": [
{
"metric": {
"__name__": "process_resident_memory_bytes",
"instance": "localhost:9090",
"job": "prometheus"
},
"values": [
[1518008453.662, "87318528"],
[1518008463.662, "87318528"],
[1518008473.662, "87318528"]
]
},
{
"metric": {
"__name__": "process_resident_memory_bytes",
"instance": "localhost:9100",
"job": "node"
},
"values": [
[1518008444.819, "17043456"],
[1518008454.819, "17043456"],
[1518008464.819, "17043456"]
]
}
]
}
}


```


This is different than the previous instant vector result, as resultType is now matrix, and each time series has multiple values. When used with a range vector, the query endpoint returns the raw samples,16 but be wary of asking for too much data at once because one end or the other may run out of memory. There is one other type of result called a scalar. Scalars don’t have labels, they are just numbers.17 http://localhost:9090/api/v1/query?query=42 would produce:


```
{
"status": "success",
"data": {
"resultType": "scalar",
"result": [1518008879.023, "42"]
}
}

```

### query_range

The query range endpoint at /api/v1/query_range is the main HTTP endpoint of Prometheus you will use, as it is the endpoint to use for graphing. Under the covers, query_range is syntactic sugar (plus some performance optimizations) for multiple calls to the query endpoint. In addition to a query URL parameter, you provide query_range with a start time, an end time, and a step. The query is first executed at the start time. Then it is executed step seconds after the start time. Then it is executed twice step seconds after the start time and so on, stopping when the query evaluation time would exceed the end time. All the instant vector18 results from the different executions are combined into a range vector and returned. For example, if you wanted to query the number of samples Prometheus ingested in the first 15 minutes of 2018, you could run the following: http://localhost:9090/api/v1/ query_range?query=rate(prometheus_tsdb_head_samples_appended_total[5m])&start =1514764800&end=1514765700&step=60, which would produce a result like:

```
{
"status": "success",
"data": {
"resultType": "matrix",
"result": [
{
"metric": {
"instance": "localhost:9090",
"job": "prometheus"
},
"values": [
[1514764800, "85.07241379310345"],
[1514764860, "102.6793103448276"],
[1514764920, "120.30344827586208"],
[1514764980, "137.93103448275863"],
[1514765040, "146.7586206896552"],
[1514765100, "146.7793103448276"],
[1514765160, "146.8"],
[1514765220, "146.8"],
[1514765280, "146.8"],
[1514765340, "146.8"],
[1514765400, "146.8"],
[1514765460, "146.8"],
[1514765520, "146.8"],
[1514765580, "146.8"],
[1514765640, "146.8"],
[1514765700, "146.8"],
]
}
]
}
}


```


There are a few aspects of this that you should take note of. The first is that the sample timestamps align with the start time and step, as each result comes from a different instant query evaluation and instant query results always use their evalua‐ tion time as the timestamp of results. The second is that the last sample here is at the end time, which is to say that the range is inclusive and the last point will be the end time if it happens to line up with the step. The third is that we selected a range of 5 minutes for the rate function, which is larger than the step. Since query_range is doing repeated instant query evaluations, there is no state being passed between the evaluations. If the range was smaller than the step, then we would have been skipping over data. For example, a 1-minute range with a 5-minute step would have ignored 80% of the samples. To prevent this you should use ranges that are at least one or two scrape intervals larger than the step you are using.

When using range vectors with query_range, you should usually use a range that is longer than your step in order to not skip data.

The fourth is that some of the samples are not particularly round, and that any numbers are round at all is due to this being a simple setup of the sample values. When working with metrics your data is rarely perfectly clean; different targets are scraped at different times and scrapes can be delayed. When performing queries that are not perfectly aligned with the underlying data or aggregating across multiple hosts, you will rarely get round results. In addition, the nature of floating-point calculations can lead to numbers that are almost round.


Here, there is a sample for each step. If it happened that there was no result for a given time series for a step, then that sample would simply be missing in the end result. If there are more than 11,000 steps for a query_range, Prometheus will reject the query with an error. This is to prevent accidentally sending extremely large queries to Prometheus, such as a 1-second step for a week. As monitors with a horizontal resolution of over 11,000 pixels are rare, you are unlikely to run into this when graphing. If you are writing reporting scripts, you can split up query_range requests that would hit this limit. This limit allows for a minute resolution for a week, or an hour of resolution for a year, so most of the time it should not apply. Aligned data When using tools like Grafana it’s common for the alignment of query_range to be based on the current time, and so your results will not align perfectly with minutes, hours, or days. While this is fine when you are looking at dashboards, it is rarely what you want with reporting scripts. query_range does not have an option to specify alignment, instead it is up to you to specify a start parameter with the right alignment. For example, if you wanted to have samples every hour on the hour in Python, the expression (time.time() // 3600) * 3600 will return the start of the current hour,19 which you can adjust in steps of 3,600 and use as the start and end URL parameters, and then use a step parameter of 3600. Now that you know the basics of how to use PromQL and execute queries via the HTTP APIs, we will go into more detail on aggregation.