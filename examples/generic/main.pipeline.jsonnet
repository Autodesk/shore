local pipeline = import 'sponnet/pipeline.libsonnet';

function(params) (
  pipeline.Pipeline {
    limitConcurrent: false,
    application: params.application,
    name: params.custom_name.a.a.a,
    stages: [
      pipeline.stages.WaitStage {
        name: 'Stage 1',
        waitTime: 1,
        skipWaitText: '${ parameters["test123"] }',
        refId: '1',
        requisiteStageRefIds: [],
      },
      pipeline.stages.WaitStage {
        name: 'Stage 2',
        waitTime: 1,
        skipWaitText: '${ parameters["test123"] }',
        refId: '2',
        requisiteStageRefIds: ['1'],
      },
      pipeline.stages.WaitStage {
        name: 'Stage 3',
        waitTime: 1,
        skipWaitText: '${ parameters["test123"] }',
        failOnFailedExpressions: true,
        refId: '3',
        requisiteStageRefIds: ['2'],
      },
      pipeline.stages.RunKubeJobStage {
        name: 'Test Output',
        application: 'kubernetes',
        consumeArtifactSource: "propertyFile",
        propertyFile: "test123",
        manifest: {
          apiVersion: 'batch/v1',
          kind: 'Job',
          metadata: {
            name: 'pi',
          },
          spec: {
            template: {
              spec: {
                containers: [{
                  command: ['sh', '-c', 'echo SPINNAKER_PROPERTY_TEST=TEST && echo SPINNAKER_PROPERTY_OUTPUT=OUTPUT && echo SPINNAKER_PROPERTY_SOMETHING=SOMETHING'],
                  image: 'alpine',
                  name: 'test123',
                }],
                restartPolicy: 'Never',
              },
            },
          },
        },
        refId: '4',
        requisiteStageRefIds: ['3'],
      },
    ],
    parameterConfig: [
      {
        default: '',
        description: '',
        hasOptions: false,
        label: 'test123',
        name: 'test123',
        options: [
          {
            value: '',
          },
        ],
        pinned: false,
        required: false,
      },
    ],
  }
)
