local pipeline = import 'sponnet/pipeline.libsonnet';

pipeline.Pipeline {
  "application": "devx",
  "name": "nested",
  "stages": [
    pipeline.stages.WaitStage {},
    pipeline.stages.NestedPipelineStage {
      name: "A nested stage level 1",
      // Top Level Parent
      Parent:: $,
      Pipeline:: {
        "name": "Child Pipeline level 1",
        "stages": [
           pipeline.stages.NestedPipelineStage {
            name: "A nested stage level 2",
           // Top Level Parent
            Parent:: $,
            Pipeline:: {
              "name": "Child Pipeline level 2",
              "stages": [
                pipeline.stages.WaitStage {},
                pipeline.stages.WaitStage {}
              ]
            }
          },
          pipeline.stages.WaitStage {}
        ]
      }
    }
  ],
}
