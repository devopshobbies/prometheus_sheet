You already learned about aggregation in “Aggregation Basics” on page 229; however, this is only a small taste of what is possible. Aggregation is important. With applica‐ tions with thousands or even just tens of instances it’s not practical for you to sift through each instance’s metrics individually. Aggregation allows you to summarize metrics not just within one application, but across applications too. There are 12 aggregation operators in PromQL, with 2 optional clauses, without and by. In this chapter you’ll learn about the different ways you can use aggregation.


### Grouping

Before talking about the aggregation operators themselves, you need to know about how time series are grouped. Aggregation operators work only on instant vectors, and they also output instant vectors.

Let’s say you have the following time series in Prometheus:


```json

node_filesystem_size_bytes{device="/dev/sda1",fstype="vfat",
instance="localhost:9100",job="node",mountpoint="/boot/efi"} 100663296
node_filesystem_size_bytes{device="/dev/sda5",fstype="ext4",
instance="localhost:9100",job="node",mountpoint="/"} 90131324928
node_filesystem_size_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run"} 826961920
node_filesystem_size_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run/lock"} 5242880
node_filesystem_size_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run/user/1000"} 826961920
node_filesystem_size_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run/user/119"} 826961920

```


There are three instrumentation labels: device, fstype, and mountpoint. There are also two target labels: job and instance. Target and instrumentation labels are a notion that you and we have, but which PromQL knows nothing about. All labels are the same when it comes to PromQL, no matter where they originated from.

## without

Generally you will always know the instrumentation labels, as they rarely change. But you do not always know the target labels in play, as an expression you write might be used by someone else on metrics originating from different scrape configs, or Prometheus servers that might also have added in other target labels across a job, such as an env or cluster label. You might even add in such target labels yourself at some point, and it’d be nice not to have to update all your expressions. When aggregating metrics you should usually try to preserve such target labels, and thus you should use the without clause when aggregating to specify the labels you want to remove. For example, the query:

When aggregating metrics you should usually try to preserve such target labels, and thus you should use the without clause when aggregating to specify the labels you want to remove. For example, the query:

`sum without(fstype, mountpoint)(node_filesystem_size_bytes)`

will group the time series, ignoring the fstype and mountpoint labels, into three groups:

```json
# Group {device="/dev/sda1",instance="localhost:9100",job="node"}
node_filesystem_size_bytes{device="/dev/sda1",fstype="vfat",
instance="localhost:9100",job="node",mountpoint="/boot/efi"} 100663296
# Group {device="/dev/sda5",instance="localhost:9100",job="node"}
node_filesystem_size_bytes{device="/dev/sda5",fstype="ext4",
instance="localhost:9100",job="node",mountpoint="/"} 90131324928
# Group {device="tmpfs",instance="localhost:9100",job="node"}
node_filesystem_size_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run"} 826961920
node_filesystem_size_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run/lock"} 5242880
node_filesystem_size_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run/user/1000"} 826961920
node_filesystem_size_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run/user/119"} 826961920

```

and the sum aggregator will apply within each of these groups, adding up the values of the time series and returning one sample per group:


```sh
{device="/dev/sda1",instance="localhost:9100",job="node"} 100663296
{device="/dev/sda5",instance="localhost:9100",job="node"} 90131324928
{device="tmpfs",instance="localhost:9100",job="node"} 2486128640

```

Notice that the instance and job labels are preserved, as would be any other labels that had been present. This is useful because any alerts you created that included this expression somehow would have additional target labels like env or cluster preserved. This provides context for your alerts and makes them more useful (also useful when graphing). The metric name has also been removed, as this is an aggregation of the node_file system_size_bytes metric rather than the original metric. When a PromQL opera‐ tor or function could change the value or meaning of a time series, the metric name is removed.

It is valid to provide no labels to the without. For example: 
`sum without()(node_filesystem_size_bytes)` 
will give you the same result as: 
`node_filesystem_size_bytes` 
with the only difference being the metric name is removed.

## by

In addition to without there is also the by clause. Where without specifies the labels to remove, by specifies the labels to keep. Accordingly, some care is required when using by to ensure you don’t remove target labels that you would like to propagate in your alerts or use in your dashboards. You cannot use both by and without in the same aggregation.

The query:

`sum by(job, instance, device)(node_filesystem_size_bytes)`

will produce the same result as the query in the preceding section using without:

```json
{device="/dev/sda1",instance="localhost:9100",job="node"} 100663296
{device="/dev/sda5",instance="localhost:9100",job="node"} 90131324928
{device="tmpfs",instance="localhost:9100",job="node"} 2486128640

```

However, if instance or job had not been specified, then they wouldn’t have defined the group and would not be in the output. Generally, you should prefer to use without rather than by for this reason.


There are two cases where you might find by more useful. The first is that unlike without, by does keep the __name__ label if told explicitly. This allows you to use expressions like:


`sort_desc(count by(__name__)({__name__=~".+"}))`

to investigate how many time series have the same metric names.

The second is cases where you do want to remove any labels you do not know about. For example, info metrics, as discussed in “Info” on page 96, are expected to add more labels over time. To count how many machines were running each kernel version, you could use:

`count by(release)(node_uname_info)`

which on our single machine test setup returns:

`sum by()(node_filesystem_size_bytes)`

and:

`sum(node_filesystem_size_bytes)`

are exactly equivalent and will give a result like:

`{} 92718116864`

This is a single time series, and that time series has no labels.
If you executed the expression:

`sum(non_existent_metric)`

the result would be an instant vector with no time series, which will show up in the expression browser’s Console tab as “no data.”

If the input to an aggregation operator is an empty instant vector, it will output an empty instant vector. Thus, count by(foo)(non_existent_metric) will be empty rather than 0, as count and other aggregators don’t have any labels to work with. count(non_existent_metric) is consistent with this, and also returns an empty instant vector.


## Operators 

All 11 aggregation operators use the same grouping logic. You can control this with one of without or by. What differs between aggregation operators is what they do with the grouped data.


### sum

sum is the most common aggregator; it adds up all the values in a group and returns that as the value for the group. For example:

`sum without(fstype, mountpoint, device)(node_filesystem_size_bytes)`

would return the total size of the filesystems of each of your machines. When dealing with counters,2 it is important that you take a rate before aggregating with sum:

`sum without(device)(rate(node_disk_read_bytes_total[5m]))`

If you were to take a sum across counters directly, the result would be meaningless, as different counters could have been initialized at different times depending on when the exporter started, restarted, or any particular children were first used.

### count

The count aggregator counts the number of time series in a group, and returns it as the value for the group. For example:

`count without(device)(node_disk_read_bytes_total)`

would return the number of disk devices a machine has. Our machine only has one disk, so we get:

`{instance="localhost:9100",job="node"} 1`

Here it is OK not to use rate with a counter, as you care about the existence of the time series rather than its value.

### Unique label values

You can also use count to count how many unique values a label has. For example, to count the number of CPUs in each of your machines, you could use:

`count without(cpu)(count without (mode)(node_cpu_seconds_total))`


The inner count3 removes the other instrumentation label, mode, returning one time series per CPU per instance:


```json
{cpu="0",instance="localhost:9100",job="node"} 8
{cpu="1",instance="localhost:9100",job="node"} 8
{cpu="2",instance="localhost:9100",job="node"} 8
{cpu="3",instance="localhost:9100",job="node"} 8

```

The outer count then returns the number of CPUs that each instance has:

`{instance="localhost:9100",job="node"} 4`


If you didn’t want a per-machine breakdown, such as if you were investigating whether certain labels had high cardinality, you could use the by modifier to look at only one label:

`count(count by(cpu)(node_cpu_seconds_total))`

which would produce a single sample with no labels, such as:

`{} 4`

## avg

The avg aggregator returns the average of the values4 of the time series in the group as the value for the group. For example:

would give you the average usage of each CPU mode for each Node Exporter instance with a result such as:

```json

{instance="localhost:9100",job="node",mode="idle"} 0.9095948275861836
{instance="localhost:9100",job="node",mode="iowait"} 0.005543103448275879
{instance="localhost:9100",job="node",mode="irq"} 0
{instance="localhost:9100",job="node",mode="nice"} 0.0013620689655172522
{instance="localhost:9100",job="node",mode="softirq"} 0.0001465517241379329
{instance="localhost:9100",job="node",mode="steal"} 0
{instance="localhost:9100",job="node",mode="system"} 0.015836206896552414
{instance="localhost:9100",job="node",mode="user"} 0.06054310344827549

```

This gives you the exact same result as:

```json

sum without(cpu)(rate(node_cpu_seconds_total[5m]))
/
count without(cpu)(rate(node_cpu_seconds_total[5m]))

```



but it is both more succinct and more efficient to use avg. When using avg, sometimes you may find that a NaN in the input is causing the entire result to become NaN. This is because any floating-point arithmetic that involves NaN will have NaN as a result. You may wonder how to filter out these NaNs in the input, but that is the wrong question to ask. Usually this is due to attempting to average averages, and one of the denominators of the first averages was 0.5 Averaging averages is not statistically valid, so what you should do instead is aggregate using sum and then finally divide, as shown in “Summary” on page 232.


- Technically it is called an arithmetic mean. In the unlikely event you need a geometric mean, the ln and exp functions combined with the avg aggregator can be used to calculate that. 
- This is as 1 / 0 = NaN.

### group

The group aggregator returns 1 for each of the time series in the group as the value for the group. For example:

`count by (instance)( group by (fstype,instance) (node_filesystem_files) )`

That query would return the number of different filesystem types for each instance. In this case, any aggregation could have worked (sum, count) in place of group. How‐ ever, using group makes it clear for anyone reading the query that we are interested in the grouping and the resulting labels themselves rather than the value produced by the inner aggregation operator.

### stddev and stdvar

The standard deviation is a statistical measure of how spread out a set of numbers is. For example, if you had the numbers [2,4,6], then the standard deviation would be 1.633.6 The numbers [3,4,5] have the same average of 4, but a standard deviation of 0.816. The main use of the standard deviation in monitoring is to detect outliers. In nor‐ mally distributed data you would expect that about 68% of samples would be within one standard deviation of the mean, and 95% within two standard deviations.7 If one instance in a job has a metric several standard deviations away from the average, that’s a good indication that something is wrong with it. For example, you could find all instances that were at least two standard deviations above the average using an expression such as:

```
some_gauge
> ignoring (instance) group_left()
(
avg without(instance)(some_gauge)
+
2 * stddev without(instance)(some_gauge)
)
```


This uses one-to-many vector matching, which will be discussed in “Many-to-One and group_left” on page 268. If your values are all tightly bunched, then this may return some time series that are more than two standard deviations away, but still operating normally and close to the average. You could add an additional filter that the value has to be at least, say, 20% higher than the average to protect against this. This is also a rare case where it is OK to take an average of an average, such as if you applied this to average latency. The standard variance is the standard deviation squared8 and has statistical uses.



- Prometheus uses the population standard deviation rather than the sample standard deviation, as you will usually be looking at all the values you are interested in rather than a random subset.


- For nonnormally distributed data, Chebyshev’s inequality provides a weaker bound.


### min and max


The min and max aggregators return the minimum or maximum value within a group as the value of the group, respectively. The same grouping rules apply as elsewhere, so the output time series will have the labels of the group.9 For example:
`max without(device, fstype, mountpoint)(node_filesystem_size_bytes)`
will return the size of the biggest filesystem on each instance, which for us returns:
`{instance="localhost:9100",job="node"} 90131324928` 
The max and min aggregators will only return NaN if all values in a group are NaN.

### topk and bottomk


topk and bottomk are different from the other aggregators discussed so far in three ways. First, the labels of time series they return for a group are not the labels of the group; second, they can return more than one time series per group; and third, they take an additional parameter. topk returns the k time series with the biggest values, so for example:

`topk without(device, fstype, mountpoint)(2, node_filesystem_size_bytes)`

would return up to two11 time series per group, such as:

```json
node_filesystem_size_bytes{device="/dev/sda5",fstype="ext4",
instance="localhost:9100",job="node",mountpoint="/"} 90131324928
node_filesystem_size_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run"} 826961920

```


As you can see, topk returns input time series with all their labels, including the __name__ label, which holds the metric name. The result is also sorted.

- If the exponentiation operator had existed at the time we were adding stdvar and stddev, then stdvar would probably not have been added.
- If you want the input time series returned, use topk or bottomk.
- In floating-point math, any comparison with NaN always returns false. Aside from causing oddities such as NaN != NaN returning false, a naive implementation of min and max would (and once did) get stuck on a NaN if it was the first value examined.
- The k is 2 in this case.

bottomk is the same as topk, except that it returns the k time series with the smallest values rather than the k biggest values. Both aggregators will, where possible, avoid returning time series with NaN values. There is a gotcha when using these aggregators with the query_range HTTP API endpoint. As was discussed in “query_range” on page 245, the evaluation of each step is independent. If you use topk, it is possible that the top time series will change from step to step. So a topk(5, some_gauge) for a query_range with 1,000 steps could in the worst case return 5,000 different time series. The way to handle this is to use the at (@) modifier, as discussed in “At Modifier” on page 242.


### quantile


The quantile aggregator returns the specified quantile of the values of the group as the group’s return value. As with topk, quantile takes a parameter. So, for example, if we wanted to know across the different CPUs in each of our machines what the 90th percentile of the system mode CPU usage is, we could use:

`quantile without(cpu)(0.9, rate(node_cpu_seconds_total{mode="system"}[5m]))`

which produces a result like:

`{instance="localhost:9100",job="node",mode="system"} 0.024558620689654007`

This means that 90% of our CPUs are spending at least 0.02 seconds per second in the system mode. This would be a more useful query if we had tens of CPUs in our machine, rather than the four it actually has. In addition to the mean, you could use quantile to show the median, 25th, and 75th percentiles12 on your graphs. For example, for process CPU usage the expressions would be:

```promql
# average, arithmetic mean
avg without(instance)(rate(process_cpu_seconds_total[5m]))
# 0.25 quantile, 25th percentile, 1st or lower quartile
quantile without(instance)(0.25, rate(process_cpu_seconds_total[5m]))
# 0.5 quantile, 50th percentile, 2nd quartile, median
quantile without(instance)(0.5, rate(process_cpu_seconds_total[5m]))
# 0.75 quantile, 75th percentile, 3rd or upper quartile
quantile without(instance)(0.75, rate(process_cpu_seconds_total[5m]))
```

This would give you a sense of how your different instances for a job are behaving, without having to graph each instance individually. This allows you to keep your dashboards readable as the number of underlying instances grows. Personally we find that per-instance graphs break down somewhere around three to five instances.

# quantile, histogram_quantile, and quantile_over_time


As you may have noticed by now, there is more than one PromQL function or operator with quantile in the name. The quantile aggregator works across an instant vector in an aggregation group. The quantile_over_time function works across a single time series at a time in a range vector. The histogram_quantile function works across the buckets of one histogram metric child at a time in an instant vector.

### count_values

The final aggregation operator is count_values. Like topk it takes a parameter and can return more than one time series from a group. What it does is build a frequency histogram of the values of the time series in the group, with the count of each value as the value of the output time series and the original value as a new label. That’s a bit of a mouthful, so we will show you an example. Say you had a time series called software_version with the following values:

```json
software_version{instance="a",job="j"} 7
software_version{instance="b",job="j"} 4
software_version{instance="c",job="j"} 8
software_version{instance="d",job="j"} 4
software_version{instance="e",job="j"} 7
software_version{instance="f",job="j"} 4
```
If you evaluated the query:


`count_values without(instance)("version", software_version)`

on these time series, you would get the result:


```
{job="j",version="7"} 2
{job="j",version="8"} 1
{job="j",version="4"} 3

```


There were two time series in the group with a value of 7, so a time series with a version="7" plus the group labels was returned with the value 2. The result is similar for the other time series. There is no bucketing involved when the frequency histogram is created; the exact values of the time series are used. Thus this is only really useful with integer values and where there will not be too many unique values. This is most useful with version numbers,13 or with the number of objects of some type that each instance of your application sees. If you have too many versions deployed at once, or different applications are continuing to see different numbers of objects, something might be stuck somewhere. count_values can be combined with count to calculate the number of unique values for a given aggregation group. For example, the number of versions of software that are deployed can be calculated with:

`count without(version)( count_values without(instance)("version", software_version) )`


which in this case would return:

`{job="j"} 3`

You could also combine count_values with count in the other direction; for exam‐ ple, to see how many of your machines had how many disk devices:


`count_values without(instance)( "devices", count without(device) (node_disk_io_now) )`

In our case we have one machine with five disk devices:


`{devices="5",job="node"} 1`

Now that you understand aggregators, we will look at binary operators, like addition and subtraction, and how vector matching works.