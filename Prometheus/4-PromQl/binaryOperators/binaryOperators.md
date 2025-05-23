You will want to do more with your metrics than simply aggregate them, which is where the binary operators come in. Binary operators are operators that take two operands,1 such as the addition and equality operators. Binary operators in allow for more than simple arithmetic on instant vectors; you can also apply a binary operator to two instant vectors with grouping based on labels. This is where the real power of PromQL comes out, allowing classes of analysis that few other metrics systems offer. PromQL has three sets of binary operators: arithmetic operators, comparison opera‐ tions, and logical operators. This chapter will show you how to use them.


### Working with Scalars

In addition to instant vectors and range vectors, there is another type of value known as a scalar.2 Scalars are single numbers with no dimensionality. For example, 0 is a scalar with the value zero, while {} 0 is an instant vector containing a single sample with no labels and the value zero.

## Arithmetic operators

You can use scalars in arithmetic with an instant vector to change the values in the instant vector. For example:










`process_resident_memory_bytes / 1024`
would return:
`{instance="localhost:9090",job="prometheus"} 21376`
`{instance="localhost:9100",job="node"} 13316`

which is the process memory usage, in kilobytes.4 You will note that the division operator was applied to all time series in the instant vector returned by the process_resident_memory_bytes selector, and that the metric name was removed as it is no longer the metric process_resident_memory_bytes.

All six arithmetic operations work similarly, with the semantics you’d expect from other programming languages. They are:


• + addition
• - subtraction
• * multiplication
• / division
• % modulo
• ^ exponentiation


The modulo operator is a floating-point modulo and can return noninteger results if you provide it with noninteger input. For example:

5 % 1.5

will return:

0.5


As this example demonstrates, you can also use binary arithmetic operators when both operands are scalars. The result will be a scalar. This is mostly useful for readability, as it is much easier to understand the intent of (1024 * 1024 * 1024) than it is 1073741824.




In addition, you can put the scalar operand on the left side of the operator and an instant vector on the right, so for example: 1e9 - process_resident_memory_bytes would subtract the process memory from a billion. You can also use arithmetic operators with instant vectors on both sides, which is covered in “Vector Matching” on page 265.

### Trigonometric Operator
The atan2 operator returns the arc tangent of the division of two vectors, using the
signs of the two to determine the quadrant of the return value:
`x atan2 y`
This operator allows you to execute atan2 on two vectors using vector matching, which isn’t available with normal functions. It acts in the same manner as arithmetic operators (+, -, *, …).

### Comparison Operators
The comparison operators are as follows, with the usual meanings:

```
• == equals
• != not equals
• > greater than
• < less than
• >= greater than or equal to
• <= less than or equal to

```
What is a little different is that the comparison operators in PromQL are filtering.
That is to say that if you had the samples:
`process_open_fds{instance="localhost:9090",job="prometheus"} 14`
`process_open_fds{instance="localhost:9100",job="node"} 7`
and used an instant vector in a comparison with a scalar, such as in the expression:
`process_open_fds > 10`
then you would get the result:
`process_open_fds{instance="localhost:9090",job="prometheus"} 14`





What you need is some way to do the comparison but not have it filter. This is what the bool modifier does; for each comparison it returns a 0 for false or a 1 for true. For example:

`process_open_fds > bool 10`

will return:
`{instance="localhost:9090",job="prometheus"} 1`
`{instance="localhost:9100",job="node"} 0`

which as expected has one output sample per sample in the input instant vector.


From there you can sum up to get the number of processes for each job that have more than 10 open file descriptors:


`sum without(instance)(process_open_fds > bool 10)`


which produces the result you originally wanted:

```
{job="prometheus"} 1
{job="node"} 0

```


You could use a similar approach to find the proportion of machines with more than four disk devices:


```
avg without(instance)(
count without(device)(node_disk_io_now) > bool 4
)

```


This works by first using a count aggregation to find the number of disks reported by each Node Exporter, then seeing how many have more than four, and finally averaging across machines to get the proportion. The trick here is that the values returned by the bool modifier are all 0 and 1, so the count is the total number of machines, and the sum is the number of machines meeting the criteria. The avg is the count divided by the sum, giving you a ratio or proportion. The bool modifier is the only way you can compare scalars, as:


```
42 <= bool 13

```
will return:
`0`
where the 0 indicates false.

### Vector Matching

Using operators between scalars and instant vectors will cover many of your needs, but using operators between two instant vectors is where PromQL’s power really starts to shine.


When you have a scalar and an instant vector, it is obvious that the scalar can be applied to each sample in the vector. With two instant vectors, which samples should apply to which other samples? This matching of the instant vectors is known as vector matching.


### One-to-One

In the simplest cases there will be a one-to-one mapping between your two vectors. Say that you had the following samples:

```
process_open_fds{instance="localhost:9090",job="prometheus"} 14
process_open_fds{instance="localhost:9100",job="node"} 7
process_max_fds{instance="localhost:9090",job="prometheus"} 1024
process_max_fds{instance="localhost:9100",job="node"} 1024
```

Then when you evaluated the expression:
```
process_open_fds
/
process_max_fds

```
you would get the result:
```
{instance="localhost:9090",job="prometheus"} 0.013671875
{instance="localhost:9100",job="node"} 0.0068359375

```
What has happened here is that samples with exactly the same labels, except for the metric name in the label __name__, were matched together. That is to say that the two samples with the labels {instance="localhost:9090",job="prometheus"} got matched together, and the two samples with the labels {instance="local host:9100",job="node"} got matched together.

In this case there was a perfect match, with each sample on both sides of the operator being matched. If a sample on one side had no match on the other side, then it would not be present in the result, as binary operators need two operands. If a binary operator returns an empty instant vector when you were expecting a result, it is probably because the labels of the samples in the operands don’t match. This is often due to a label that is present on one side of the operator but not the other.

Sometimes you will want to match two instant vectors whose labels do not quite match. Similar to how aggregation allows you to specify which labels matter, as discussed in “Grouping” on page 249, vector matching also has clauses controlling which labels are considered.

You can use the ignoring clause to ignore certain labels when matching, similar to how without works for aggregation. Say you were working with node_cpu_ seconds_total, which has cpu and mode as instrumentation labels, and wanted to know what proportion of time was being spent in the idle mode for each instance. You could use the expression:

```
sum without(cpu)(rate(node_cpu_seconds_total{mode="idle"}[5m]))
/ ignoring(mode)
sum without(mode, cpu)(rate(node_cpu_seconds_total[5m]))

```

This will give you a result such as:

`{instance="localhost:9100",job="node"} 0.8423353718871361`

Here the first sum produces an instant vector with a mode="idle" label, whereas the second sum produces an instant vector with no mode label. Usually vector matching will fail to match the samples, but with ignoring(mode) the mode label is discarded when the vectors are being grouped, and matching succeeds. As the mode label was not in the match group, it is not in the output.


You can tell the preceding expression is correct in terms of vector matching by inspection, without having to know anything about the underlying time series. The removal of cpu is balanced on both sides, and ignoring(mode) handles one side having a mode and the other not. This can be trickier when there are different time series with differ‐ ent labels in play, but looking at expressions in terms of how the labels flow is a handy way for you to spot errors.

The on clause allows you to consider only the labels you provide, similar to how by works for aggregation. The expression:

```
sum by(instance, job)(rate(node_cpu_seconds_total{mode="idle"}[5m]))
/ on(instance, job)
sum by(instance, job)(rate(node_cpu_seconds_total[5m]))

```

will produce the same result as the previous expression,7 but as with by, the on clause has the disadvantage that you need to know all labels that are currently on the time series or that may be present in the future in other contexts.

The value that is returned for the arithmetic operators is the result of the calculation, but you may be wondering what happens for the comparison operators when there are two instant vectors. The answer is that the value from the lefthand side is returned. For example, the expression:

```
process_open_fds
>
(process_max_fds * .5)

```

will return for you the value of process_open_fds for all instances whose open file descriptors are more than halfway to the maximum.

If you had instead used:

```
(process_max_fds * .5)
<
process_open_fds

```

you would get half the maximum file descriptors as the return value. While the result will have the same labels, this value might be semantically less useful when alerting9 or when used in a dashboard! In general, a current value is more informative than the limit, so you should try to structure your math so that the most interesting number is on the lefthand side of a comparison.

### Many-to-One and group_left

If you were to remove the matcher on mode from the preceding section and try to evaluate:

```
sum without(cpu)(rate(node_cpu_seconds_total[5m]))
/ ignoring(mode)
sum without(mode, cpu)(rate(node_cpu_seconds_total[5m]))

```

you would get the error:

```
multiple matches for labels:
many-to-one matching must be explicit (group_left/group_right)
```


This is because the samples no longer match one-to-one, as there are multiple sam‐ ples with different mode labels on the lefthand side for each sample on the righthand side. This can be a subtle failure mode, as a time series may appear later on that breaks your expression. You can see that this is a potential issue, as looking at the label flow there’s nothing restricting the mode label to one potential value10 on the lefthand side.

Errors like this are usually due to incorrectly written expressions, so PromQL does not attempt to do anything smart by default. Instead, you must specifically request that you want to do many-to-one matching using the group_left modifier. group_left lets you specify that there can be multiple matching samples in the group of the lefthand operand.11 For example:


```

sum without(cpu)(rate(node_cpu_seconds_total[5m]))
/ ignoring(mode) group_left
sum without(mode, cpu)(rate(node_cpu_seconds_total[5m]))

```

will produce one output sample for each different mode label within each group on the lefthand side:


```
{instance="localhost:9100",job="node",mode="irq"} 0
{instance="localhost:9100",job="node",mode="nice"} 0
{instance="localhost:9100",job="node",mode="softirq"} 0.00005226389784152013
{instance="localhost:9100",job="node",mode="steal"} 0
{instance="localhost:9100",job="node",mode="system"} 0.01720353303949279
{instance="localhost:9100",job="node",mode="user"} 0.10345203045243238
{instance="localhost:9100",job="node",mode="idle"} 0.8608691486211044
{instance="localhost:9100",job="node",mode="iowait"} 0.01842302398912871

```


group_left always takes all of its labels from samples of your operand on the lefthand side. This ensures that the extra labels that are on the left side that require this to be many-to-one vector matching are preserved.

This is much easier than having to run a one-to-one expression with a matcher for each potential mode label: group_left does it all for you in one expression. You can use this approach to determine the proportion each label value within a metric represents of the whole, as shown in the preceding example, or to compare a metric from a leader of a cluster against the replicas.

There is another use for group_left—adding labels from info metrics to other metrics from a target. Instrumentation with info metrics was covered in “Info” on page 96. The role of info metrics is to allow you to provide labels that would be useful for a target or metric to have but that would clutter up the metric if you were to use it as a normal label. The prometheus_build_info metric, for example, provides you with build informa‐ tion from Prometheus:

```
prometheus_build_info{branch="HEAD",goversion="go1.10",
instance="localhost:9090",job="prometheus",
revision="bc6058c81272a8d938c05e75607371284236aadc",version="2.2.1"}

```

You can join this with metrics such as up:


```
up
* on(instance) group_left(version)
prometheus_build_info

```

which will produce a result like:

`{instance="localhost:9090",job="prometheus",version="2.2.1"} 1`

You can see that the version label has been copied over from the righthand operand to the lefthand operand as was requested by group_left(version), in addition to returning all the labels from the lefthand operand as group_left usually does. You can specify as many labels as you like to group_left, but usually it’s only one or two.13 This approach works no matter how many instrumentation labels the lefthand side has, as the vector matching is many-to-one.

The preceding expression used on(instance), which relies on each instance label only being used for one target within your Prometheus. While this is often the case, it isn’t always, so you may also need to add other labels such as job to the on clause. prometheus_build_info applies to a whole target. There are also info-style14 metrics such as node_hwmon_sensor_label mentioned in “Hwmon Collector” on page 130 that apply to children of a different metric:

```

node_hwmon_sensor_label{chip="platform_coretemp_0",instance="localhost:9100", job="node",label="core_0",sensor="temp2"} 1
node_hwmon_sensor_label{chip="platform_coretemp_0",instance="localhost:9100", job="node",label="core_1",sensor="temp3"} 1
node_hwmon_temp_celsius{chip="platform_coretemp_0",instance="localhost:9100", job="node",sensor="temp1"} 42
node_hwmon_temp_celsius{chip="platform_coretemp_0",instance="localhost:9100", job="node",sensor="temp2"} 42
node_hwmon_temp_celsius{chip="platform_coretemp_0",instance="localhost:9100", job="node",sensor="temp3"} 41

```
The node_hwmon_sensor_label metric has children that match with some (but not all) of the time series in node_hwmon_temp_celsius. In this case you know that there is only one additional label (which is called label), so you can use ignoring with group_left to add this label to the node_hwmon_temp_celsius samples:

```
node_hwmon_temp_celsius
* ignoring(label) group_left(label)
node_hwmon_sensor_label

```

which will produce results such as:

```
{chip="platform_coretemp_0",instance="localhost:9100", job="node",label="core_0",sensor="temp2"} 42
{chip="platform_coretemp_0",instance="localhost:9100", job="node",label="core_1",sensor="temp3"} 41

```


Notice that there is no sample with sensor="temp1" as there was no such sample in node_hwmon_sensor_label (how to match sparse instant vectors will be covered in “or operator” on page 271).

There is also a group_right modifier that works in the same way as group_left except that the one and the many sides are switched, with the many side now being your operand on the righthand side. Any labels you specify in the group_right modifier are copied from the left to the right. For the sake of consistency, you should prefer group_left.

### Many-to-Many and Logical Operators

There are three logical or set operators you can use:
• or union
• and intersection
• unless set subtraction


There is no not operator, but the absent function discussed in “Missing Series, absent, and absent_over_time” on page 287 serves a similar role. All the logical operators work in a many-to-many fashion, and they are the only operators that work many-to-many. They are different from the arithmetic and com‐ parison operators you have already seen in that no math is performed; all that matters is whether a group contains samples.

### or operator


In the preceding section, node_hwmon_sensor_label did not have a sample to go with every node_hwmon_temp_celsius, so results were only returned for samples that were present in both instant vectors. Metrics with inconsistent children, or whose children are not always present, are tricky to work with, but you can deal with them using the or operator.


How the or operator works is that for each group where the group on the lefthand side has samples, then they are returned; otherwise, the samples in the group on the righthand side are returned. If you are familiar with SQL, this operator can be used in a similar way as the SQL COALESCE function, but with labels.

Continuing the example from the preceding section, or can be used to substitute the missing time series from node_hwmon_sensor_label. All you need is some other time series that has the labels you need, which in this case is node_hwmon_temp_celsius. node_hwmon_temp_celsius does not have the label label, but all the other labels match up so you can ignore this using ignoring:

```
node_hwmon_sensor_label
or ignoring(label)
(node_hwmon_temp_celsius * 0 + 1)

```

The vector matching produced three groups of labels. The first two groups had a sample from node_hwmon_sensor_label so that was what was returned, including the metric name as there was nothing to change it. For the third group, however, which included sensor="temp1", there was no sample in the group for the lefthand side, so the values in the group from the righthand side were used. Because arithmetic operators were used on the value, the metric name was removed.

x * 0 + 1 will change all15 the values of the x instant vector to 1. This is also useful when you want to use group_left to copy labels, as 1 is the identity element for multiplication, which is to say it does not change the value you are multiplying.

This expression can now be used in the place of node_hwmon_sensor_label:


```
node_hwmon_temp_celsius
* ignoring(label) group_left(label)
(
node_hwmon_sensor_label
or ignoring(label)
(node_hwmon_temp_celsius * 0 + 1)
)

```

which will produce:


```
{chip="platform_coretemp_0",instance="localhost:9100", job="node",sensor="temp1"} 42
{chip="platform_coretemp_0",instance="localhost:9100", job="node",label="core_0",sensor="temp2"} 42
{chip="platform_coretemp_0",instance="localhost:9100", job="node",label="core_1",sensor="temp3"} 41

```

The sample with sensor="temp1" is now present in your result. It has no label called label, which is the same as saying that that label label has the empty string as a value. In simpler cases you will be working with metrics without any instrumentation labels. For example, you might be using the textfile collector, as covered in “Textfile Collec‐ tor” on page 134, and expecting it to expose a metric called node_custom_metric. In the event that metric doesn’t exist, you would like to return 0 instead. In cases like this, you can use the up metric that is associated with every target:

```
node_custom_metric
or
up * 0

```

This has a small problem in that it will return a value even for a failed scrape, which is not how scraped metrics work.16 It will also return results for other jobs. You can fix this with a matcher and some filtering:



```
node_custom_metric
or
(up{job="node"} == 1) * 0

```
Another way you can use the or operator is to return the larger of two series:

```

(a >= b) or b

```

If a is larger it will be returned by the comparison, and then the or operator since the group on the lefthand side was not empty. If, on the other hand, b is larger, then the comparison will return nothing, and or will return b as the group on the lefthand side was empty.


### unless operator

The unless operator does vector matching in the same way as the or operator, working based on whether groups from the right and left operands are empty or have samples. The unless operator returns the lefthand group, unless the righthand group has members, in which case it returns no samples for that group. You can use unless to restrict what time series are returned based on an expression. For example, if you wanted to know the average CPU usage of processes except those using less than 100 MB of resident memory, you could use the expression:

```
rate(process_cpu_seconds_total[5m])
unless
process_resident_memory_bytes < 100 * 1024 * 1024

```

unless can also be used to spot when a metric is missing from a target. For example:

```
up{job="node"} == 1
unless
node_custom_metric

```

would return a sample for every instance that was missing the node_custom_metric metric, which you could use in alerting. By default, as with all binary operators, unless looks at all labels when grouping. If node_custom_metric had instrumentation labels, you could use on or ignoring to check that at least one relevant time series existed without having to know the values of the other labels:

```
up == 1
unless on (job, instance)
node_custom_metric

```

Even if there are multiple samples from the right operand in a group, this is OK as unless uses many-to-many matching.

## and operator


The and operator is the opposite of the unless operator. It returns a group from the lefthand operand only if the matching righthand group has samples; otherwise, it returns no samples for that match group. You can think of it as an if operator.17 You will use the and operator most commonly in alerting to specify more than one condition. For example, you might want to return when both latency is high and there is more than a trickle of user requests. To do this for Prometheus for handlers that were taking over a second on average and had at least one request per second, you could use:


```
(
rate(http_request_duration_seconds_sum{job="prometheus"}[5m])
/
rate(http_request_duration_seconds_count{job="prometheus"}[5m])
) > 1
and
rate(http_request_duration_seconds_count{job="prometheus"}[5m]) > 1



```

This will return a sample for every individual handler on every prometheus job, so it could get a little spammy even with the one request per second restriction. Usually you would want to aggregate across a job when alerting. You can use on and ignoring with the and operator, as you can with the other binary operators. In particular, on() can be used to have a condition that has no common labels at all between the two operands. You can use this, for example, to limit the time of day an expression will return results for:

```
(
rate(http_request_duration_microseconds_sum{job="prometheus"}[5m])
/
rate(http_request_duration_microseconds_count{job="prometheus"}[5m])
) > 1000000
and
rate(http_request_duration_microseconds_count{job="prometheus"}[5m]) > 1
and on()
hour() >= 9 < 17

```

The hour function is covered in “minute, hour, day_of_week, day_of_month, day_of_year, days_in_month, month, and year” on page 284; it returns an instant vector with one sample with no labels and the hour of the UTC day of the query evaluation time as the value.

### Operator Precedence

When evaluating an expression with multiple binary operators, PromQL does not simply go from left to right. Instead, there is an order of operators that is largely the same as the order used in other languages:

1. ^
2. * / % atan2
3. + -
4. == != > < >= <=
5. unless and
6. or
For example, a or b * c + d is the same as a or ((b * c) + d).


All operators except ^ are left-associative. That means that a / b * c is the same as
(a / b) * c, but a ^ b ^ c is a ^ (b ^ c).
You can use parentheses to change the order of evaluation. We also recommend
adding parentheses where the evaluation order may not be immediately clear for an
expression, as not everyone will have memorized the operator precedence.
Now that you understand both aggregators and operators, let’s look at the final part of
PromQL: functions.