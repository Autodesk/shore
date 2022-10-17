local main = import '../main.pipeline.jsonnet';
local assertions = import './assertions.json';

local tests = [
  main({application: 'test', pipeline_name: 'test'}),
];

{
  pass: true,
  assertions: assertions,
  tests: tests,
}