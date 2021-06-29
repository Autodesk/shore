local pipeline = import 'spin-lib-jsonnet/pipeline.libsonnet';
local spel = import 'spin-lib-jsonnet/spel.libsonnet';
local evaluateArtifacts = import 'armory-lib-jsonnet/evaluate-artifacts.libsonnet';
local parameter = import 'spin-lib-jsonnet/parameter.libsonnet';

function(params) (
  pipeline.Pipeline {
    limitConcurrent: false,
    application: params.application,
    name: params.pipeline_name,
    Stages:: [
      evaluateArtifacts.EvaluateArtifactStage {
        name: "eval!",
        EvaluateArtifacts:: [
          evaluateArtifacts.EvaluateArtifact {
            name: 'artifact1',
            contents: '{"name": %s}' % [spel.expression(spel.parameter('data'))],
          },
        ],
      },
    ],
    parameterConfig: [
      parameter.Parameter {
        name: 'data',
        description: 'data',
        label: 'data'
      },
    ],
  }
)
