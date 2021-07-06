local deployment = import '../main.pipeline.jsonnet';

local tests = [
  deployment({application: 'application', pipeline: 'pipeline', example_value: 'example_value'}),
];

local assertions = [
  {application: 'application', pipeline: 'pipeline', message: "Hello example_value!"}
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
