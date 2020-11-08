local pipeline = import 'sponnet/pipeline.libsonnet';

function(params) (
  pipeline.Pipeline {
    limitConcurrent: false,
    application: params.application,
    name: params.custom_name.a.a.a,
    stages: [
      pipeline.stages.WaitStage {
        waitTime: 1,
        skipWaitText: '${ parameters["test123"] }',
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
