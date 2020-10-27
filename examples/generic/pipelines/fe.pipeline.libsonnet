local pipeline = import 'sponnet/pipeline.libsonnet';

pipeline.Pipeline {
  "application": "test1test2test3",
  "name": "test1",
  "keepWaitingPipelines": false,
  "limitConcurrent": true,
  "spelEvaluator": "v4",
  "stages": [
    {
      "isNew": true,
      "name": "Wait",
      "refId": "1",
      "requisiteStageRefIds": [],
      "type": "wait",
      "waitTime": 30
    }
  ],
  "triggers": [],
}
