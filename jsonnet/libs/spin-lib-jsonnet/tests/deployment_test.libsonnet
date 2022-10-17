local deployment = import '../deployment.libsonnet';

local tests = [
  deployment.RedBlackStrategy,
  deployment.Moniker { app: 'myApp', stack: 'test' },
  deployment.Capacity,
  deployment.Rollback,
];

local assertions = [
  {
    strategy: 'redblack',
    delayBeforeDisableSec: '0',
    delayBeforeScaleDownSec: '0',
    scaleDown: true,
    maxRemainingAsgs: '3',
    rollback: { onFailure: false },
  },

  { app: 'myApp', detail: '', stack: 'test' },

  { desired: 1, max: 1, min: 1 },

  { onFailure: false },
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
