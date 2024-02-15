# Broomstick

Turns the node pools of a kube cluster ON/OFF based on the scheduled times specified in the input configuration.

## Help

### -config <required>

Specify YAML file of the schema below.

```YAML
apiVersion: v1
clusters:
- id: sikka
	service_account: <base_64_encoded_service_account>
	project_id: sikka-dev
	location: asia-south1-b
	cluster: sikka-1
	schedule:
		monday:
			start_time: 9:00
			end_time: 22:00
```

### -schedule 

Specifies if the scheduler must run.

### -instance <required if scheduler is false> 

Specifies the id of the cluster to perform the action on.

### -action <required if scheduler is false> 

Specifies if the cluster must be turned on or off. Possible values ON/OFF.