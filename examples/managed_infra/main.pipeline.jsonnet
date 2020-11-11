local pipeline = import 'sponnet/pipeline.libsonnet';

pipeline.Pipeline {
  "application": "devx",
  "name": "nested",
  "stages": [
    pipeline.stages.PipelineStage {
      "application": "devx",
      name:  "First Stage",
      pipeline: pipeline.Pipeline {
        "application": "devx",
        "name": "first",
        "stages": [
          pipeline.stages.WaitStage {}
        ],
      },
    },
    pipeline.stages.PipelineStage {
      "application": "devx",
      name:  "Second Stage",
      pipeline: pipeline.Pipeline {
        "application": "devx",
        "name": "second",
        "stages": [
          pipeline.stages.PipelineStage {
            "application": "devx",
            name:  "recursion layer 1 stage",
            pipeline: pipeline.Pipeline {
              "application": "devx",
              "name": "Recursion layer 1 pypline",
              "stages": [

                pipeline.stages.PipelineStage {
                  "application": "devx",
                  name:  "recursion layer 2 stage",
                  pipeline: pipeline.Pipeline {
                    "application": "devx",
                    "name": "Recursion layer 2 pipeline",
                    "stages": [
                      pipeline.stages.WaitStage {}
                    ],
                  },
                }


              ],
            },
          }
        ],
      },
    }
  ],
}