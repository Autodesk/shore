local pipeline = import 'sponnet/pipeline.libsonnet';

pipeline.Pipeline {
  "application": "test1test2test3",
  "name": "test1",
  "stages": [
    pipeline.stages.WaitStage {}
  ],
}
