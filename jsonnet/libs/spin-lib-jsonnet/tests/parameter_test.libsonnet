local parameter = import '../parameter.libsonnet';

local tests = [
  parameter.Parameter {
    name: 'name',
  },
];

local assertions = [
  {
    default: '',
    description: '',
    hasOptions: false,
    label: '',
    name: 'name',
    options: [],
    pinned: false,
    required: false,
  },
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
