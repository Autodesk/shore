/**
    @file Contains objects for creating the preconditions for the CheckPreconditionsStage.
**/

/**
    Creates a base PreCondition object.

    Used mainly for inheretence by other objects.

    @example
        local myPreCondition = precondition.PreCondition {
            type: 'stageStatus',
            context: {
                stageName: 'myStage',
                stageStatus: 'SUCCEEDED',
            }
        }

        ...

        local myPreConStage = stage.CheckPreconditionsStage {
            ...
            preconditions: [
                myPreCondition,
                ...
            ]
            ...
        }
    @property {String} type - The type of pre-condition this is.
    @property {Object} context - The context that will be evaluated.
    @property {Boolean} [failPipeline=true] - Whether or not to fail the pipeline.
**/
local PreCondition = {
  type: error '`type` (String) property is required for PreCondition',
  context: error '`context` (Object) property is required for PreCondition`',
  failPipeline: true,
};

/**
    Creates a StageStatusPreCondition object.

    This precondition checks the status of a given stage in the same pipeline.

    @example
        local myPreCondition = precondition.StageStatusPreCondition {
            stageName: 'stage-name',
            stageStatus: 'SUCCEEDED'
        }

        ...

        local myPreConStage = stage.CheckPreconditionsStage {
            ...
            preconditions: [
                myPreCondition,
                ...
            ]
            ...
        }

    @class
    @augments PreCondition

    @property {String} stageName - Name of the stage to check.
    @property {String} stageStatus - The expected status of the stage that is being checked.
**/
local StageStatusPreCondition = PreCondition {
  local this = self,

  stageName:: error '`stageName` (String) property is required for StageStatusPreCondition',
  stageStatus:: error '`stageStatus` (String) property is required for StageStatusPreCondition',

  context: {
    stageName: this.stageName,
    stageStatus: this.stageStatus,
  },

  type: 'stageStatus',
};

/**
    Creates a ExpressionPreCondition object.

    This precondition checks the given SpEL expression.

    @example
        local myPreCondition = precondition.ExpressionPreCondition {
            expression: '#{ true == true }'
        }

        ...

        local myPreConStage = stage.CheckPreconditionsStage {
            ...
            preconditions: [
                myPreCondition,
                ...
            ]
            ...
        }

    @class
    @augments PreCondition

    @property {String} expression - The SpEL expression to evaluate.
**/
local ExpressionPreCondition = PreCondition {
  local this = self,

  expression:: error '`expression` (String) property is required for StageStatusPreCondition',

  context: {
    expression: this.expression,
    failureMessage: this.failureMessage,
  },

  failureMessage:: '',
  type: 'expression',
};

/**
    Creates a ClusterSizePreCondition object.

    This precondition checks if a Spinnaker Cluster is of the expected size.

    It can only check a Spinnaker Cluster that exists in the same Spinnaker Application as the pipeline.

    The comparison can be - ==, >=, <=, >, <

    @example
        local myPreCondition = precondition.ClusterSizePreCondition {
            cluster: 'potato-dev',
            comparison: '==',
            credentials: 'spinnaker-account',
            expected: 1,
            regions: ['us-west-2'],
        }

        ...

        local myPreConStage = stage.CheckPreconditionsStage {
            ...
            preconditions: [
                myPreCondition,
                ...
            ]
            ...
        }

    @class
    @augments PreCondition

    @property {String} cluster - Cluster to check.
    @property {String} comparison - Comparison operator to use - ==, >=, <=, >, <.
    @property {String} credentials - Spinnaker account to use.
    @property {int} expected - Expected size of the cluster.
    @property {Array<String>} regions - Regions to check in.
**/
local ClusterSizePreCondition = PreCondition {
  local this = self,

  cluster:: error '`cluster` (String) property is required for StageStatusPreCondition',
  comparison:: error '`comparison` (String) property is required for StageStatusPreCondition',
  credentials:: error '`credentials` (String) property is required for StageStatusPreCondition',
  expected:: error '`expectedSize` (int) property is required for StageStatusPreCondition',
  regions:: error '`regions` (Array<String>) property is required for StageStatusPreCondition',

  context: {
    cluster: this.cluster,
    comparison: this.comparison,
    credentials: this.credentials,
    expected: this.expected,
    regions: this.regions,
  },

  type: 'clusterSize',
};

// Exposed for public use.
{
  PreCondition:: PreCondition,

  ClusterSizePreCondition:: ClusterSizePreCondition,
  ExpressionPreCondition:: ExpressionPreCondition,
  StageStatusPreCondition:: StageStatusPreCondition,
}
