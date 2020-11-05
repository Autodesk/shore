local pipeline = import 'sponnet/pipeline.libsonnet';

function(params) (
  pipeline.Pipeline {
    "application": params.application,
    "name": params.custom_name.a.a.a,
    "stages": [
      pipeline.stages.WaitStage {}
    ],
  }
)
