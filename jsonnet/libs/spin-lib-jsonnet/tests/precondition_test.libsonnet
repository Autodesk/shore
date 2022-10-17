local precondition = import '../precondition.libsonnet';

local tests = [
  precondition.PreCondition { type: 'stageStatus', context: { stageName: 'stage-name', stageStatus: 'SUCCEEDED' } },
  precondition.PreCondition { type: 'stageStatus', context: { stageName: 'stage-name', stageStatus: 'FAILED_CONTINUE' }, failPipeline: false },

  precondition.StageStatusPreCondition { stageName:: 'stage-name', stageStatus:: 'SUCCEEDED' },
  precondition.ExpressionPreCondition { expression:: '#{ true == true }' },
  precondition.ClusterSizePreCondition {
    cluster:: 'potato-dev',
    comparison:: '==',
    credentials:: 'spinnaker-account',
    expected:: 1,
    regions:: ['us-west-2'],
  },
];

local assertions = [
  {
    context: {
      stageName: 'stage-name',
      stageStatus: 'SUCCEEDED',
    },
    failPipeline: true,
    type: 'stageStatus',
  },
  {
    context: {
      stageName: 'stage-name',
      stageStatus: 'FAILED_CONTINUE',
    },
    failPipeline: false,
    type: 'stageStatus',
  },

  {
    context: {
      stageName: 'stage-name',
      stageStatus: 'SUCCEEDED',
    },
    failPipeline: true,
    type: 'stageStatus',
  },
  {
    context: {
      expression: '#{ true == true }',
      failureMessage: '',
    },
    failPipeline: true,
    type: 'expression',
  },
  {
    context: {
      cluster: 'potato-dev',
      comparison: '==',
      credentials: 'spinnaker-account',
      expected: 1,
      regions: [
        'us-west-2',
      ],
    },
    failPipeline: true,
    type: 'clusterSize',
  },
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
