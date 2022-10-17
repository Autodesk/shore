/**
    @file Contains objects for creating various Spinnaker Stages.

    Both `refId` and `requisiteStageRefIds` attributes on all stages are set to '' and [] respectively.

    Connections between stage can be done by hand, at which point those two attributes need to be overwritten.

    Otherwise stage.grapher.libsonnet can be used to connect stages and set those attributes.
**/

/**
    The basic/generic Stage object.

    Contains the bare minimum any type of stage needs.

    Can be used as a basis to make any type of stage not readily avaliable in this library.

    @example
        local myHandMadeWaitStage = Stage {
            name: 'My Wait Stage',
            type: 'wait',
            waitTime: 30
        }

    @property {String} name - Name of the stage.
    @property {String} type - Type of the stage.
    @property {String} [comments] - The comments field for the stage.
    @property {String} [refId=""] - The stage's unique ID.
    @property {Array<String>} [requisiteStageRefIds=[]] - An array of `refId` strings of stages that this stage depends on.
**/
local Stage = {
  name: error '`name` (String) property is required for Stage',
  type: error '`type` (String) property is required for Stage',

  refId: '',
  requisiteStageRefIds: [],
};

/**
    Creates a Find Artifacts From Execution stage.

    This stage attempts to find (match) a Spinnaker Artifact from a Pipeline execution.

    Please use the artifacts.libsonnet for creating Artifacts for this stage.

    @example
        local httpFileArtifact = artifacts.HttpFileArtifact { ... };
        local defaultHttpFileArtifact = artifacts.HttpFileArtifact { ... };
        loccal artifactToFind = artifacts.ExpectedArtifact {
            matchArtifact: httpFileArtifact,
            useDefaultArtifact: true,
            defaultArtifact: defaultHttpFileArtifact,
        };

        local myStage = FindArtifactFromExecutionStage {
            name: 'Find My Artifact',
            application: 'my spinnaker app',
            pipeline: 'my pipeline',
            expectedArtifact: artifactToFind
        }

    @class
    @augments Stage

    @property {String} application - Name of the Spinnaker Application to find the pipeline in.
    @property {Object<Artifact>} expectedArtifact - The expected Spinnaker Artifact.
    @property {String} pipeline - Name of the Spinnaker Pipeline to find the Spinnaker Artifact in.
**/
local FindArtifactFromExecutionStage = Stage {
  application: error '`application` (String) property is required for FindArtifactFromExecutionStage',
  expectedArtifact: error '`expectedArtifact` (Object<Artifact>) property is required for FindArtifactFromExecutionStage',
  pipeline: error '`pipeline` (String) property is required for FindArtifactFromExecutionStage',

  executionOptions: {
    successful: true,
  },
  type: 'findArtifactFromExecution',
};

/**
    Creates a Bake stage.

    While the defaults will render a Bake stage, it is recommened to fine tune them to your environment.

    @example
        local myStage = BakeStage {
            name: 'Bake My Service',
            package: 'potato=1.0.0',
            regions: ['us-west-2', 'ca-central-1'],
            templateFileName: 'custom-debian-packer-template.json',
        }

    @class
    @augments Stage

    @property {String} package - The package to bake.
    @property {Array<String>} regions - The regions to bake in.
    @property {String} [amiName=''] - The AMI name, this will be appended to.
    @property {String} [baseAmi=''] - The base AMI to use.
    @property {String} [baseLabel='release'] - The base label.
    @property {String} [baseOs='ubuntu'] - The base OS.
    @property {String} [cloudProviderType='aws'] - The cloud provider type.
    @property {Object} [extendedAttributes={}] - Extended attributes to provide to the bake.
    @property {Boolean} [rebake=false] - Whether or not to rebake every time or just on changes.
    @property {String} [storeType='ebs'] - Storage type to attach for the AMI.
    @property {String} [templateFileName=''] - Packer tempalte file to use. Must exist inside of the Spinnaker deploy.
    @property {String} [vmType='hvm'] - VM type.
**/
local BakeStage = Stage {
  package: error '`package` (String) property is required for BakeStage',
  regions: error '`regions` (Array<String>) property is required for BakeStage',

  amiName: '',
  baseAmi: '',
  baseLabel: 'release',
  baseOs: 'ubuntu',
  cloudProviderType: 'aws',
  extendedAttributes: {},
  rebake: false,
  storeType: 'ebs',
  templateFileName: '',
  type: 'bake',
  vmType: 'hvm',
};

/**
    Creates a Find Image stage.

    This is primarily used for extending into FindAmiImageStage or FindContainerImageStage.

    Can be used for creation of new Find Image stages that are missing from the library.

    @example
        local myHandMadeFindAmiStage = FindImageStage {
            name: 'Find My AMI',
            cloudProvider: 'aws',
            cloudProviderType: 'aws',
            packageName: 'potato',
            regions: ['us-west-2', 'ca-central-1'],
        }

    @class
    @augments Stage

    @property {String} cloudProvider - The cloud provider to use for this stage.
    @property {String} cloudProviderType - The cloud provider to use for this stage.
    @property {Boolean} [isNew=true] - Whether or not to pick the newest image if multiple are found.
**/
local FindImageStage = Stage {
  cloudProvider: error '`cloudProvider` (String) property is required for FindImageStage',
  cloudProviderType: error '`cloudProviderType` (String) property is required for FindImageStage',

  isNew: true,
  type: 'findImageFromTags',
};

/**
    Creates a Find Image stage to find AWS AMIs.

    @example
        local myStage = FindAmiImageStage {
            name: 'Find My AMI',
            packageName: 'potato',
            regions: ['us-west-2', 'ca-central-1'],
        }

    @class
    @augments FindImageStage

    @property {String} packageName - The package name to find the AMI with.
    @property {Array<String>} regions - The regions to search in.
**/
local FindAmiImageStage = FindImageStage {
  packageName: error '`packageName` (String) property is required for FindAmiImageStage',
  regions: error '`regions` (Array<String>) property is required for FindAmiImageStage',

  cloudProvider: 'aws',
  cloudProviderType: 'aws',
  tags: {},
};

/**
    Creates a Find Image stage to find Docker Image.

    The imageLabelOrSha is the full Docker Image path.

    It can be set to an ECR or any other Docker registry.

    @example
        local myStage = FindAmiImageStage {
            name: 'Find My Docker Image',
            imageLabelOrSha: '12345678901.dkr.ecr.us-west-2.amazonaws.com/helloworld:latest',
        }

    @class
    @augments FindImageStage

    @property {String} imageLabelOrSha - The image, with a label or SHA256, to find.
**/
local FindContainerImageStage = FindImageStage {
  imageLabelOrSha: error '`imageLabelOrSha` (String) property is required for FindContainerImageStage',

  cloudProvider: 'ecs',
  cloudProviderType: 'ecs',
};

/**
    Creates a Manual Judgement stage.

    @example
        local myStage = ManualJudgmentStage {
            name: 'My Manual Judgement',
            instructions: 'Check the potatoes.'
        };

        local myStageWithOpts = ManualJudgmentStage {
            name: 'My Manual Judgement',
            instructions: 'Check the potatoes.'
            judgmentInputs": [
            {
              "value": "fried"
            },
            {
              "value": "boiled"
            },
            {
              "value": "roasted"
            }
          ],
        };

    @class
    @augments Stage

    @property {String} [instructions=''] - The instructions to provide at the Manual Judgement stage. Supports a limited subset of HTML tags.
    @property {Array<String>} [judgmentInputs=[]] - Additional selections to provide context around the Manual Judgement stage.
**/
local ManualJudgmentStage = Stage {
  instructions: '',
  judgmentInputs: [],
  type: 'manualJudgment',
};

/**
    Creates a Deploy stage.

    @example
        local myFirstCluster = { ... }
        local mySecondCluster = { ... }

        local myStage = DeployStage {
            name: 'Deploy My Two Clusters',
            clusters [
                myFirstCluster,
                mySecondCluster
            ],
        }

    @class
    @augments Stage

    @property {Array<Object>} clusters - An array of objects representing clusters
**/
local DeployStage = Stage {
  clusters: error '`clusters` (Array<Object>) property is required for DeployStage',

  type: 'deploy',
};

/**
    Creates a Pipeline stage.

    This triggers another pipeline.

    The `pipeline` attribute needs to be the Spinnaker Pipeline ID or the JSON representation of the full pipeline.

    If it is the JSON representation of the full pipeline, that pipeline will be created/updated when this stage is
    rendered.

    Please use the pipeline.libsonnet to create Spinnaker Pipelines.

    @example
        local myPipeline = pipeline.Pipeline { ... }

        local myStage = PipelineStage {
            name: 'Call Another Pipeline',
            application: 'another spinnaker app',
            pipeline: myPipeline,
            pipelineParameters: {
                'parameter-1': 'value-1',
                'parameter-2': 'value-2',
            },
        }

    @class
    @augments Stage

    @property {String} application - The Spinnaker Application that contains the Spinnaker Pipeline to call.
    @property {String} pipeline - The Spinnaker Pipeline to call.
    @property {Boolean} [failPipeline=true] - To fail this, calling pipeline, if the called pipeline fails.
    @property {Object} [pipelineParameters={}] - A map/dict of parameters to pass to the called pipeline.
    @property {Boolean} [waitForCompletion=true] - To wait for completion or not.
**/
local PipelineStage = Stage {
  application: error '`application` (String) property is required for PipelineStage',
  pipeline: error '`pipeline` (String) property is required for PipelineStage',

  pipelineJSON:: '',

  failPipeline: true,
  pipelineParameters: {},
  type: 'pipeline',
  waitForCompletion: true,

};

/**
    Creates a Pipeline stage.

    This stage requires supplying both the Parent Pipeline that will do the calling and the Child Pipeline that will be
    called.

    Please use the pipeline.libsonnet to create Spinnaker Pipelines.


    @example
        local myParentPipeline = pipeline.Pipeline { ... }
        local myChildPipeline = pipeline.Pipeline { ... }

        local myStage = NestedPipelineStage {
            name: 'Call the Child Pipeline',
            Parent: myParentPipeline,
            Pipeline: myChildPipeline,
        }

    @example
        local myParentPipeline = pipeline.Pipeline {
            ...
            stages: [
                NestedPipelineStage {
                    name: "child pipeline",
                    Parent: $,
                    Pipeline: pipeline.Pipeline { ... }
                }
            ]
        }

    @class
    @augments PipelineStage

    @property {Pipeline} Parent - The parent pipeline.
    @property {Pipeline} Pipeline - The child pipeline.
**/
local NestedPipelineStage = PipelineStage {
  Parent:: error '`Parent` (Object<Pipeline>) property is required for NestedPipelineStage',
  Pipeline:: error '`Pipeline` (Object<Pipeline>) property is required for NestedPipelineStage',

  local this = self,
  local pipeline = this.Pipeline,
  local innerPipeline = pipeline { application: this.Parent.application },

  application: this.Parent.application,
  pipeline: innerPipeline,
};

/**
    Creates a Wait stage.
    Time is in seconds.
    Default waitTime is 1 second.

    @example
        local myStage = WaitStage {
            name: 'Wait 30 seconds',
            waitTime: 30
        }

    @class
    @augments Stage

    @property {int} [waitTime=1] - The time to wait.
**/
local WaitStage = Stage {
  type: 'wait',
  waitTime: 1,
};

/**
    Creates a Webhook stage.

    Webhook stages can have all the properties listed at the link below:

    {@link https://github.com/spinnaker/orca/blob/master/orca-webhook/src/main/groovy/com/netflix/spinnaker/orca/webhook/pipeline/WebhookStage.groovy#L81}


    name, type, url, and method are required.

    @example
        local myParentPipeline = pipeline.Pipeline { ... }
        local myChildPipeline = pipeline.Pipeline { ... }

        local myStage = WebhookStage {
            name: 'Get My README',
            url: "http://localhost/README.md",
            method: "GET",
        }

    @class
    @augments Stage

    @property {String} url - The URL to query.
    @property {String} method - The HTTP method to use.
**/
local WebhookStage = Stage {
  url: error '`url` (String) property is required for WebhookStage',
  method: error '`method` (String) property is required for WebhookStage',

  type: 'webhook',

  // Other fields
  // customHeaders: {},
  // failFastStatusCodes: [],
  // payload: {},
};

/**
    Creates a Run Job (Manifest) stage, that runs a Kubernetes Job/Container.

    Please use the kube.libsonnet for creating Manifest object.

    Please use the artifact.libsonnet for creating required (stage uses), expected (stage produces) artifacts, and/or
    the manifest object.

    @example
        local myManifest = kube.Manifest { ... };
        // local myManifestArtifact = artifact.GitHubFileArtifact { ... };

        local myStage = RunKubeJobStage {
            name: 'Run My Kube Job',
            account: 'kube-account',
            credentials: 'kube-account',
            manifest: myManifest,
            // manifestArtifact: myManifestArtifact
        }

    @class
    @augments Stage

    @property {String} account - The Spinnaker account to run in.
    @property {String} application - The Spinnaker Application which this stage is in.s
    @property {String} credentials - The Kuberenets account to use.
**/
local RunKubeJobStage = Stage {
  account: error '`account` (String) property is required for RunKubeJobStage',
  application: error '`application` (String) property is required for RunKubeJobStage',
  credentials: error '`credentials` (String) property is required for RunKubeJobStage',

  alias: 'runJob',
  cloudProvider: 'kubernetes',
  consumeArtifactSource: 'none',
  manifest: {},
  manifestArtifact: {},
  source: 'text',
  type: 'runJobManifest',

  // Other fields
  // expectedArtifacts: [],
  // requiredArtifacts: [],
};

/**
    Creates a Check Preconditions stage.

    Please use the precondition.libsonnet for creating PreCondition objects.

    @example
        local successfulStagePreCon = precondition.StageStatusPreCondition { ... };
        local atLeastOneClusterPreCon = precondition.ClusterSizePreCondition { ... };

        local myStage = CheckPreconditionsStage {
            name: 'Check if Stage Passed and Cluster Deployed',
            preconditions: [
                successfulStagePreCon,
                atLeastOneClusterPreCon
            ],
        }

    @class
    @augments Stage

    @property {Array<PreCondition>} preconditions - The pre-conditions for this stage to check.
**/
local CheckPreconditionsStage = Stage {
  preconditions: error '`preconditions` (Array<PreCondition>) property is required for CheckPreconditionsStage',

  type: 'checkPreconditions',
  requisiteStageRefIds: 'none',
};

/**
    Creates a Rollback Cluster stage.

    The cluster must be in the same application as the pipeline that this stage is in.

    Please use the deployment.libsonnet for creating the Moniker object.

    @example
        local clusterMoniker = deployment.Moniker { ... };

        local myStage = RollbackClusterStage {
            name: 'Rollback My Cluster',
            cluster: 'potato-dev-producer',
            credentials: 'spinnaker-account',
            moniker: clusterMoniker,
            regions: [ 'us-west-2' ]
        }

    @class
    @augments Stage

    @property {String} cluster - Cluster to rollback.
    @property {String} credentials - Spinnaker account to use.
    @property {Moniker} moniker - Moniker to rollback.
    @property {Array<String>} regions - Regions to rollback in.
    @property {int} [targetHealthyRollbackPercentage=100] - The target must meet this precentage to not rollback.
**/
local RollbackClusterStage = Stage {
  cluster: error '`cluster` (String) property is required for RollbackClusterStage',
  credentials: error '`credentials` (String) property is required for RollbackClusterStage',
  moniker: error '`moniker` (Object<Moniker>) property is required for RollbackClusterStage',
  regions: error '`regions` (Array<String>) property is required for RollbackClusterStage',

  cloudProvider: 'aws',
  cloudProviderType: 'aws',
  targetHealthyRollbackPercentage: 100,
  type: 'rollbackCluster',
};

/**
    Creates a Destroy Server Group stage.

    The cluster must be in the same application as the pipeline that this stage is in.

    @example
        local myStage = DestroyServerGroupStage {
            name: 'Destroy My Cluster',
            cluster: 'potato-dev-producer',
            credentials: 'spinnaker-account',
            regions: [ 'us-west-2' ],
            target: 'current_asg_dynamic'
        }

    @class
    @augments Stage

    @property {String} cluster - Cluster to destroy.
    @property {String} credentials - Spinnaker account to use.
    @property {Array<String>} regions - Regions the server group in.
    @property {String} target - The target about which server group will be destroyed.
    @property {String} cloudProvider - The cloud driver provider to use - can be either "aws" or "ecs", must match `cloudProviderType`.
    @property {String} cloudProviderType - The cloud driver provider to use - can be either "aws" or "ecs", must match `cloudProvider`.
**/
local DestroyServerGroupStage = Stage {
  cluster: error '`cluster` (String) property is required for DestroyServerGroupStage',
  credentials: error '`credentials` (String) property is required for DestroyServerGroupStage',
  regions: error '`regions` (Array<String>) property is required for DestroyServerGroupStage',
  target: error '`target` (String) property is required for DestroyServerGroupStage',
  cloudProvider: error '`cloudProvider` (String) property is required for DestroyServerGroupStage',
  cloudProviderType: error '`cloudProviderType` (String) property is required for DestroyServerGroupStage',

  type: 'destroyServerGroup',
};

/**
    Holds stages that should be ran in parallel.

    This is used by the stage.grapher.libsonnet.

    @example
        local myFirstStage = ...
        local mySecondStage = ...
        local myThirdStage = ...
        local myFourthStage = ...


        local myParalleStages = stage.Parallel {
            parallelStages: [
                mySecondStage,
                myThirdStage
            ]
        }

        local myStages = [
            myFirstStage,
            myParalleStages,
            myFourthStage
        ]

        stageGrapher.addRefIdsAndRequisiteRefIds(myStages)
**/
local Parallel = {
  parallelStages: error '`parallelStages` (Array<Stage>) is a required property of `Parallel`',
};

/**
    Creates an object to attach to a Stage, making the stage execute only when the supplied SpEL expression is evaluated
    to true.

    In the Spinnaker UI, this is the checkbox and field called "Conditional on Expression".

    @example
        local conditionalDeploy = stage.StageEnabled {
            expression: '${ true == true }'
        }
        local myDeployStage = stage.Deploy { ... } + conditionalDeploy
**/
local StageEnabled = {
  local this = self,
  expression:: error '`expression` (String) is a required property of `StageEnabled`',

  comments: 'This stage only runs if the expression evaluates to true.',
  stageEnabled: {
    expression: this.expression,
    type: 'expression',
  },
};

// Exposed for public use.
{
  Parallel:: Parallel,
  Stage:: Stage,
  StageEnabled:: StageEnabled,

  BakeStage:: BakeStage,
  CheckPreconditionsStage:: CheckPreconditionsStage,
  DeployStage:: DeployStage,
  FindAmiImageStage:: FindAmiImageStage,
  FindArtifactFromExecutionStage:: FindArtifactFromExecutionStage,
  FindContainerImageStage:: FindContainerImageStage,
  FindImageStage:: FindImageStage,
  ManualJudgmentStage:: ManualJudgmentStage,
  NestedPipelineStage:: NestedPipelineStage,
  PipelineStage:: PipelineStage,
  RollbackClusterStage:: RollbackClusterStage,
  DestroyServerGroupStage:: DestroyServerGroupStage,
  RunKubeJobStage:: RunKubeJobStage,
  WaitStage:: WaitStage,
  WebhookStage:: WebhookStage,
}
