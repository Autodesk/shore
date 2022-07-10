local pipeline = import 'spin-lib-jsonnet/pipeline.libsonnet';
local stage = import 'spin-lib-jsonnet/stage.libsonnet';
local parameter = import 'spin-lib-jsonnet/parameter.libsonnet';
local kube = import 'spin-lib-jsonnet/kube.libsonnet';

function(params) (
  pipeline.Pipeline {
    limitConcurrent: false,
    application: params.application,
    name: params.pipeline,
    Stages:: [
      stage.WaitStage {
        name: 'Stage 1',
        waitTime: 1,
        skipWaitText: '${ parameters["test123"] }',
      },
      stage.Parallel {
        parallelStages: [
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
        ],
      },
      stage.RunKubeJobStage {
        name: 'Test Output',
        application: 'testoutput',
        consumeArtifactSource: "propertyFile",
        account: 'kubernetes',
        credentials: 'kubernetes',
        propertyFile: "test123asdad",
        manifest: kube.Manifest {
          generateName:: 'test234123',
          containers:: [
            kube.Container {
              command: ['sh', '-c', 'echo SPINNAKER_PROPERTY_TEST=TEST && echo SPINNAKER_PROPERTY_OUTPUT=OUTPUT && echo SPINNAKER_PROPERTY_SOMETHING=SOMETHING'],
              image: 'alpine',
              name: 'test123',
            }
          ],
          labels:: {
            purpose: 'spinnaker-run-kube-job-stage',
          }
        }
      },
      stage.NestedPipelineStage {
        name: 'Nesting',
        Parent: $,
        Pipeline: pipeline.Pipeline {
          application: $.application,
          name: 'Nested Test',
          Stages:: [
            stage.WaitStage {
              name: 'Stage 1',
              waitTime: 1,
              skipWaitText: '${ parameters["test123"] }',
            },
          ],
        },
      },
    ],
    parameterConfig: [
      parameter.Parameter {
        name: 'test123',
        description: 'test123',
        label: 'test123'
      },
    ],
  }
)
