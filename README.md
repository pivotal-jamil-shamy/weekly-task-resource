# Weekly Task Resource

Implements a resource that reports new versions when the current time matches the hour of the day and day of the week o fthe config.

---
## Update your pipeline

Update your pipeline to include this new declaration of resource types. See the example pipeline yml snippet below or the Concourse docs for more details [here](https://concourse.ci/configuring-resource-types.html).
```
---
resource_types:
- name: weekly-task-resource
  type: registry-image
  source:
    repository: <image repo name>

resources:
  - name: thursdays-at-10-pm
    type: weekly-task-resource
    source:
      day_to_fire: Thursday
      hour_to_fire: 22
      location: America/Toronto
```

## Source Configuration

* `location`: *Optional.* Defaults to UTC. Accepts any timezone that
  can be parsed by https://godoc.org/time#LoadLocation

  e.g.
  
  `America/Toronto`

  `America/Vancouver`
  
  `America/New_York`
  
* `day_to_fire`: The weekday you want this resource to fire. It accepts : `Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday`

  e.g.
  
  `Tuesday`

  `Thursday`
  
* `hour_to_fire`: The hour of the day you want this resource to fire. It accepts whole numbers within the [0-23] range. This resource will only fire onnce during this hour. If the pipeline is paused during that hour, and then unpaused after it, it will not fire.

  e.g.
  
  `21`

  `19`

## Behavior

### `check`: Report the current time.

Returns `time.Now()` as the version only if we are within the hour of the `hour_to_fire` and on the day of `day_to_fire`. The first time the script runs it will fire if we happen to be in that interval.

#### Parameters

*None.*

### `in`: Report the given time

If triggered by `check`, returns the original version as the resulting
version.

#### Parameters

1. *Output directory.* The directory where the in script will store
   the requested version

### `out`: Not supported.

## Developer Notes

You can test the behavior by simulating Concourse's invocations. For example:

```
$ echo '{"version": {"time":"2020-05-08T09:10:57.725589-04:00"}, "source":{"hour_to_fire": 9 , "day_to_fire": "Thursday","location":"America/Toronto"}}' | go run ./check

$ echo '{"source":{"hour_to_fire": 9 , "day_to_fire": "Thursday","location":"America/Toronto"}}' | go run ./check

$ echo '{"source":{"hour_to_fire": 15 , "day_to_fire": "Monday","location":"America/Toronto"}}' | go run ./check
```
