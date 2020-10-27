// TODO - Refactoring: Consider splitting up into different files - focusing on Stage/Object/Pipeline/...
local stages = {
  Stage:: {
    name: error 'name required',
    type: error 'type required',
    dependsOn: [],
    executionOptions: {
      successful: true,
    },
  },
  Parallel:: {
    parallelStages: error '`parallelStages` is a required property of `Parallel`',
  },
  StageEnabled:: {
    local this = self,
    expression:: error 'expression is required',

    comments: 'This stage only runs if the expression evaluates to true.',
    stageEnabled: {
      expression: this.expression,
      type: 'expression',
    },
  },
  FindArtifactFromExecutionStage:: $.Stage {
    type: 'findArtifactFromExecution',
    application: error 'application required',
    pipeline: error 'pipeline required',
    expectedArtifact: error 'expectedArtifact required',
  },
  BakeStage:: $.Stage {
    type: 'bake',
    amiName: error 'amiName required',
    templateFileName: error 'templateFileName required',
    baseAmi: '',
    vmType: 'hvm',
    storeType: 'ebs',
    regions: ['us-west-2'],
    rebake: false,
    baseLabel: 'release',
    baseOs: 'ubuntu',
    cloudProviderType: 'aws',
    extendedAttributes: {},
  },
  AwsDeployStage:: $.Stage {
    type: 'deploy',
    clusters: error 'clusters required',
  },
  FindImageStage:: $.Stage {
    type: 'findImageFromTags',
    cloudProvider: error 'cloudProvider is required',
    cloudProviderType: error 'cloudProviderType is required',
  },
  FindAmiImageStage:: $.FindImageStage {
    cloudProvider: 'aws',
    cloudProviderType: 'aws',
    packageName: error 'packageName required',
    regions: error 'regions required',
    tags: {},
  },
  FindContainerImageStage:: $.FindImageStage {
    cloudProvider: 'ecs',
    cloudProviderType: 'ecs',
    imageLabelOrSha: error 'imageLabelOrSha required',
  },
  ManualJudgmentStage:: $.Stage {
    type: 'manualJudgment',
    judgmentInputs: [],
    instructions: error 'instructions required',
  },
  DeployStage:: $.Stage {
    type: 'deploy',
    clusters: error 'clusters required',
  },
  PipelineStage:: $.Stage {
    application: error 'application required',
    name: error 'name required',
    pipeline: error 'application required',
    pipelineParameters: {},
    type: 'pipeline',
    waitForCompletion: true,
    failPipeline: true,
  },
  WaitStage:: $.Stage {
    name: 'Wait',
    type: 'wait',
    waitTime: 1,
  },
  RunKubeJobStage:: $.Stage {
    account: 'kubernetes',
    cloudProvider: 'kubernetes',
    credentials: 'kubernetes',
    type: 'runJobManifest',
    alias: 'runJob',
    source: 'text',
    name: error 'property `name` is required for RunKubeJobStage',
    application: error 'property `application` is required for RunKubeJobStage',
    manifestArtifact: {},
    consumeArtifactSource: 'none',
    manifest: {},
  },
};

local artifacts = {
  Artifact:: {
    name: error '`name` property is required for Artifact',
    type: error '`type` property is required for Artifact',
    artifactAccount: error '`artifactAccount` property is required for Artifact',
    kind: error '`kind` property is required for Artifact',
  },
  GithubFileArtifact:: $.Artifact {
    type: 'github/file',
    kind: 'github',
    version: 'master',
    reference: error '`reference` property is required for Artifact',
  },
  GithubRepoArtifact:: $.Artifact {
    type: 'git/repo',
    kind: 'custom',
    version: 'master',
    reference: error '`reference` property is required for Artifact',
  },
  HttpFileArtifact:: $.Artifact {
    type: 'http/file',
    kind: 'http',
    reference: error '`reference` property is required for Artifact',
  },
  EmbeddedArtifact:: $.Artifact {
    type: 'embedded/base64',
    kind: 'embedded',
    reference: error '`reference` property is required for Artifact',
  },
  ExpectedArtifact:: {
    displayName: error 'displayName required',
    useDefaultArtifact: true,
    usePriorArtifact: false,
    defaultArtifact: error 'defaultArtifact must be defined',
    matchArtifact: {},
  },
  NewExpectedArtifact(id):: $.ExpectedArtifact {
    local this = self,
    artifactName:: '',
    artifactAccount:: '',
    artifactType:: '',
    artifactKind:: '',
    artifactReference:: 'property `artifactReference` for `NewExpectedArtifact`',
    matchArtifact: $.Artifact {
      artifactAccount: this.artifactAccount,
      name: id,
      type: this.artifactType,
      kind: this.artifactKind,
    } + (if std.objectHas(this, 'artifactReference') then { reference: this.artifactReference } else {}),
    defaultArtifact: if this.defaultArtifact == null then this.matchArtifact else this.defaultArtifact,
  } + (if std.length(id) > 0 then { id: id } else {}),
};

local specs = {
  KubeJobSpec:: {
    apiVersion: 'batch/v1',
    kind: 'Job',
    metadata: error '`metadata` property is required for KubeJobSpec',
    spec: error '`spec` property is required for KubeJobSpec',
  },
};

local triggers = {
  Trigger:: {
    type: error 'type required',
    enabled: true,
  },
  PipelineTrigger:: $.Trigger {
    type: 'pipeline',
    application: error 'application required',
    pipeline: error 'pipeline required',
    status: [
      'successful',
    ],
  },
  NewPipelineTrigger(trigger): {
    fields:: std.objectFields(trigger),
    data: $.PipelineTrigger { application: trigger.application, pipeline: trigger.pipeline },
  },
  JenkinsTrigger:: $.Trigger {
    type: 'jenkins',
    master: error 'master required',
    job: error 'job required',
  },
  WebhookTrigger:: $.Trigger {
    type: 'webhook',
    enabled: true,
    source: error 'property `source` is required for WebhookTrigger',
  },
  NewWebhookTrigger(service): $.WebhookTrigger {
    source: service,
  },
  NewJenkinsTrigger(trigger): {
    fields:: std.objectFields(trigger),
    data: $.JenkinsTrigger {
      master: trigger.server,
      job: trigger.job,
      propertyFile: if std.objectHas(trigger, 'propertyFile') then trigger.propertyFile else '',
    },
  },
  NewTriggerByType(trigger): {
    fields:: std.objectFields(trigger),
    template:
      if trigger.type == 'pipeline' then
        $.PipelineTrigger(trigger)
      else if trigger.type == 'jenkins' then
        $.JenkinsTrigger(trigger)
      else
        {},
  },
};

local parameters = {
  Parameter:: {
    name: error 'property `name` is required for Parameter',
    required: false,
    hasOptions: false,
    pinned: false,
  },
};

{
  stages:: stages,
  artifacts:: artifacts,
  notifications:: {},
  triggers:: triggers,
  parameters:: parameters,
  specs:: specs,
  Pipeline:: {
    application: error 'application required',
    name: error 'name required',
    keepWaitingPipelines: false,
    limitConcurrent: true,
    parallel: false,
    stages: [],
    triggers: [],
    parameterConfig: [],
    expectedArtifacts: [],
    connectStages(stages):: [
      stages[i] {
        refId: '' + (i + 1),
        requisiteStageRefIds: if i > 0 then ['' + i] else [],
      }
      for i in std.range(0, std.length(stages) - 1)
    ],
  },
  NewCluster:: {
    account: error 'account required',
    application: error 'application required',
    associatePublicIpAddress: null,
    availabilityZones: error 'availabilityZones required',
    capacity: $.NewCapacity,
    cloudProvider: 'aws',
    cooldown: 10,
    copySourceCustomBlockDeviceMappings: false,
    delayBeforeDisableSec: 0,
    delayBeforeScaleDownSec: 0,
    ebsOptimized: false,
    enabledMetrics: [],
    freeFormDetails: '',
    healthCheckGracePeriod: 600,
    healthCheckType: 'EC2',
    iamRole: 'SPNKR-C-UW2-BaseRole',
    instanceMonitoring: false,
    instanceType: 't2.micro',
    keyPair: 'default',
    loadBalancers: [],
    maxRemainingAsgs: '2',
    moniker: error 'moniker required',
    preferSourceCapacity: true,
    provider: 'aws',
    rollback: $.NewRollback,
    scaleDown: true,
    securityGroups: error 'securityGroups required',
    spotPrice: '',
    stack: error 'stack required',
    strategy: 'redblack',
    subnetType: 'app',
    suspendedProcesses: [],
    tags: {},
    targetGroups: error 'targetGroups required',
    targetHealthyDeployPercentage: 100,
    terminationPolicies: ['Default'],
    useAmiBlockDeviceMappings: false,
    useSourceCapacity: true,
  },
  NewRedBlackStrategy:: {
    local this = self,
    rollbackFailure:: false,
    // rendered
    strategy: 'redblack',
    delayBeforeDisableSec: '0',
    delayBeforeScaleDownSec: '0',
    scaleDown: true,
    maxRemainingAsgs: '3',
    rollback: {
      onFailure: this.rollbackFailure,
    },
  },
  NewMoniker:: {
    app: error 'app required',
    detail: '',
    stack: error 'stack required',
  },
  NewCapacity:: {
    desired: 1,
    max: 1,
    min: 1,
  },
  NewRollback:: {
    onFailure: false,
  },
  NewOption(value): {
    data: { value: value },
  },
}
