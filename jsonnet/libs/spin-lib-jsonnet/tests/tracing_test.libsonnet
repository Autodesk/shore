local tracing = import '../tracing.libsonnet';

// Data used for tests - setting to true should create output.
local VERBOSE = false;

local tests = [
  { key: 'potato' } + tracing.traceFields('my-object', { firstName: 'John', lastName: 'Doe' }, {}),
  { key: 'potato' } + tracing.traceFields('my-object', { firstName: 'John', lastName: 'Doe' }, { additionKey: 'tomato' }),

  tracing.log({ key: 'apple' }),
  { key: 'apple' } + tracing.log('hello world', { additionKey: 'orange' }),
];

local assertions = [
  {
    key: 'potato',
  },
  {
    additionKey: 'tomato',
    key: 'potato',
  },

  {
    key: 'apple',
  },
  {
    additionKey: 'orange',
    key: 'apple',
  },
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
