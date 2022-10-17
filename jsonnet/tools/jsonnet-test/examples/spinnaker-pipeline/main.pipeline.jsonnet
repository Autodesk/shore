local pipeline = import 'spin-lib-jsonnet/pipeline.libsonnet';
local stage = import 'spin-lib-jsonnet/stage.libsonnet';
local parameter = import 'spin-lib-jsonnet/parameter.libsonnet';

function(params) (
  pipeline.Pipeline {
    limitConcurrent: false,
    application: params.application,
    name: params.pipeline_name,
    Stages:: [
      stage.WaitStage {
        name: 'Stage 1',
        waitTime: 1,
        skipWaitText: '${ parameters["test123"] }',
      },
      stage.WaitStage {
        name: 'Stage 2',
        waitTime: 1,
        skipWaitText: '${ parameters["test123"] }',
      },
      stage.WaitStage {
        name: 'Stage 3',
        waitTime: 1,
        skipWaitText: '${ parameters["test123"] }',
        failOnFailedExpressions: true,
      },
      stage.RunKubeJobStage {
        name: 'Test Output',
        application: 'kubernetes',
        account: 'test',
        credentials: 'test',
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
      },
    ],
    parameterConfig: [
      parameter.Parameter {
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
