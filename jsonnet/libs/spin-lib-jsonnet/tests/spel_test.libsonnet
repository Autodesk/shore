local spel = import '../spel.libsonnet';

// Data used for tests.
local stageName = 'my-stage-name';
local parameterName = 'parameter-name';
local artifactName = 'artifact-name';
local artifactPrefix = 'artifact-prefix';

local tests = [
  spel.simpleDateFormat('yyMMdd'),

  spel.executionStages(stageName),

  spel.getTriggerArtifactReference(artifactName),
  spel.getTriggerArtifact(artifactName),

  spel.parameter(parameterName),
  spel.stage(stageName),
  spel.execution('status'),

  spel.toJson('SpEL-object'),

  spel.expression(spel.parameter(parameterName)),

  spel.newString('new Date()'),
  spel.toBase64('cow'),
  spel.fromBase64('deadbeef=='),

  spel.getTriggerArtifactWithPrefix(artifactPrefix),
];

local assertions = [
  'new java.text.SimpleDateFormat("yyMMdd").format(new java.util.Date())',

  "execution.stages.?[name matches 'my-stage-name']",

  "trigger['artifacts'].?[name == 'artifact-name'][0]['reference']",
  "trigger['artifacts'].?[name == 'artifact-name']",

  'parameters["parameter-name"]',
  "#stage('my-stage-name')",
  'execution["status"]',

  "#toJson('SpEL-object')",

  '${parameters["parameter-name"]}',

  'new String(new Date())',
  "#toBase64('cow')",
  "#fromBase64('deadbeef==')",
  "trigger['artifacts'].?[name.startsWith('artifact-prefix')]",
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
