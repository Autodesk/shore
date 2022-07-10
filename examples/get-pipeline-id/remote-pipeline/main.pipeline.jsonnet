local pipeline = import 'spin-lib-jsonnet/pipeline.libsonnet';
local stage = import 'spin-lib-jsonnet/stage.libsonnet';
local trigger = import 'spin-lib-jsonnet/trigger.libsonnet';
local artifact = import 'spin-lib-jsonnet/artifact.libsonnet';
local parameter = import 'spin-lib-jsonnet/parameter.libsonnet';

function(params) (
  pipeline.Pipeline {
    limitConcurrent: false,
    application: params.application,
    name: params.pipeline,
    Stages:: [
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
    triggers: [
      trigger.PipelineTrigger {
        application: 'test1test2test3',
        pipeline: 'get-pipeline-id-target-pipeline',
      },
      trigger.WebhookTrigger {
        source: 'abc123',
      },
    ],
  }
)
