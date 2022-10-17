local trigger = import '../trigger.libsonnet';

local tests = [
  trigger.Trigger { type: 'generic' },
  trigger.GenericTrigger { type: 'test', target: 'potato-field' },

  trigger.PipelineTrigger { application: 'other-application', pipeline: 'other-pipeline' },
  trigger.JenkinsTrigger { master: 'jenkins-master-server', job: 'my-jenkins-job' },
  trigger.WebhookTrigger { source: 'my-pipelines-webhook' },
];

local assertions = [
  {
    enabled: true,
    expectedArtifactIds: [],
    type: 'generic',
  },
  {
    type: 'test',
  },
  {
    application: 'other-application',
    enabled: true,
    expectedArtifactIds: [],
    pipeline: 'other-pipeline',
    status: [
      'successful',
    ],
    type: 'pipeline',
  },
  {
    enabled: true,
    expectedArtifactIds: [],
    job: 'my-jenkins-job',
    master: 'jenkins-master-server',
    propertyFile: '',
    type: 'jenkins',
  },
  {
    enabled: true,
    expectedArtifactIds: [],
    payloadConstraints: {},
    source: 'my-pipelines-webhook',
    type: 'webhook',
  },
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
