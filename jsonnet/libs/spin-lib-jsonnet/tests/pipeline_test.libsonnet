local pipeline = import '../pipeline.libsonnet';

local tests = [
  pipeline.Pipeline {
    name: 'name',
    application: 'application',
  },
  pipeline.Pipeline {
    name: 'name',
    application: 'application',
    Stages:: [{ name: 1, type: 'a' }, { name: 2, type: 'a' }],
  },
];

local assertions = [
  {
    application: 'application',
    expectedArtifacts: [],
    keepWaitingPipelines: false,
    limitConcurrent: true,
    name: 'name',
    parameterConfig: [],
    spinLibJsonnetVersion: '0.0.0',
    stages: [],
    triggers: [],
  },
  {
    application: 'application',
    expectedArtifacts: [],
    keepWaitingPipelines: false,
    limitConcurrent: true,
    name: 'name',
    parameterConfig: [],
    spinLibJsonnetVersion: '0.0.0',
    stages: [
      {
        name: 1,
        refId: '1',
        requisiteStageRefIds: [],
        type: 'a',
      },
      {
        name: 2,
        refId: '2',
        requisiteStageRefIds: [
          '1',
        ],
        type: 'a',
      },
    ],
    triggers: [],
  },
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
