// local pipeline = import 'sponnet/pipeline.libsonnet';

// function(params) (
//   pipeline.NewPipeline {
//     local this = self,
//     region:: error '`region` is a required parameter for FePipeline instance',
//     envName:: error '`envName` is a required parameter for FePipeline instance',
//     envObj:: error '`envObj` is a required parameter for FePipeline instance',
//     infraRegionConfig:: getInfraRegionConfig(params.envName, params.region, params.infraAdf),
//     parameterConfig: parameters.getEnvParams(params),
//     triggers: [],
//     stages: self.connectStagesByIndex(createStages(params.region, params.envName, params.envObj, this)),
//     // notifications: notifications.newEnvNotifications(params.pipelineAdf),
//     keepWaitingPipelines: true,
//   }
// )

{
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
