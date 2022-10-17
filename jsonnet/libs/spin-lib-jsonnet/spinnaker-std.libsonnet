local application = './application.libsonnet';
local artifact = './artifact.libsonnet';
local deployment = 'deployment.libsonnet';
local kube = './kube.libsonnet';
local notification = './notification.libsonnet';
local parameter = './parameter.libsonnet';
local pipeline = './pipeline.libsonnet';
local preconidition = './preconidition.libsonnet';
local stage = './stage.libsonnet';
local trigger = './trigger.libsonnet';

{
  application:: application,
  artifact:: artifact,
  deployment:: deployment,
  kube:: kube,
  notification:: notification,
  parameter:: parameter,
  pipeline:: pipeline,
  preconidition:: preconidition,
  stage:: stage,
  trigger:: trigger,
}
