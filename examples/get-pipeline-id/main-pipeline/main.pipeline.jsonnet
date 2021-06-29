local pipeline = import 'spin-lib-jsonnet/pipeline.libsonnet';
local stage = import 'spin-lib-jsonnet/stage.libsonnet';
local artifact = import 'spin-lib-jsonnet/artifact.libsonnet';
local parameter = import 'spin-lib-jsonnet/parameter.libsonnet';
local kube = import 'spin-lib-jsonnet/kube.libsonnet';

function(params) (
  pipeline.Pipeline {
    limitConcurrent: false,
    application: params.application,
    name: params.pipeline_name,
    Stages:: [
      stage.Parallel {
        parallelStages: [
          stage.WaitStage {
            name: 'Stage 1',
            waitTime: 1,
            skipWaitText: '${ parameters["test123"] }',
          },
          stage.WaitStage {
            name: 'Stage 2',
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
        propertyFile: "test123",
        manifest: kube.Manifest {
          generateName:: 'test123',
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
          name: '%s-nested' % params.pipeline_name,
          Stages:: [
            stage.WaitStage {
              name: 'Stage 1',
              waitTime: 1,
              skipWaitText: '${ parameters["test123"] }',
            },
            stage.PipelineStage {
              name: 'Trigger Them Pipeline Nested',
              application: params.application,
              pipeline: 'get-pipeline-id-target-pipeline',
              pipelineParameters: {
                data: 'I Shall Invoke This Pipeline in the same application!',
              },
            },
            stage.FindArtifactFromExecutionStage {
              name: 'Get Artifact',
              application: 'test1test2test3',
              pipeline: 'get-pipeline-id-target-pipeline',
              expectedArtifact: artifact.ExpectedArtifact {
                matchArtifact: {
                  artifactAccount: 'embedded-artifact',
                  name: 'artifact1',
                  type: 'embedded/base64',
                },
              },
            },
          ],
        },
      },
      stage.PipelineStage {
        name: 'Trigger Them Pipeline Top Level',
        application: params.application,
        pipeline: 'get-pipeline-id-target-pipeline',
        pipelineParameters: {
          data: 'I Shall Invoke This Remote Application Pipeline!',
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
