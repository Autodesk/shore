/**
    @file Contains objects for creating an EC2 deploy stage configuration.
    @todo Create the ECS counterpart.
**/

/**
    Server group capacity.

    For EC2 this means how many instances are in a given server group (aka AWS autoscaling group), for a given Cluster.

    @example
        local myCapacity = deployment.Capacity { desired: 3, max: 6 }

        ...

        local myCluster = deployment.Cluster {
            ...
            capacity: myCapacity,
            ...
        }

    @property {int} [desired=1] - Desired deployment capacity.
    @property {int} [max=1] - Maximum deployment capacity.
    @property {int} [min=1] - Minimum deployment capacity.
**/
local Capacity = {
  desired: 1,
  max: 1,
  min: 1,
};

/**
    Rollback configuration object for the given Cluster.

    @example
        local rollbackOnFail = deployment.Rollback { onFailure: true }

        ...

        local myCluster = deployment.Cluster {
            ...
            rollback: rollbackOnFail,
            ...
        }

    @property {Boolean} [onFailure=false] - Whether or not to rollback.
**/
local Rollback = {
  onFailure: false,
};

/**
    A map of all the possible Spinnaker deployment strategies.

    This list is made from:

    {@link https://github.com/spinnaker/orca/blob/master/orca-clouddriver/src/main/groovy/com/netflix/spinnaker/orca/kato/pipeline/strategy/Strategy.java}

    @enum {String}
    @readonly

    @property {String} REDBLACK - Red/Black (aka Blue/Green) deployment strategy.
    @property {String} ROLLINGREDBLACK - Rolling Red/Black deployment strategy.
    @property {String} MONITORED - Monitored strategy.
    @property {String} CFROLLINGREDBLACK - Cloud Foundry Red/Black deployment strategy.
    @property {String} HIGHLANDER - Highlander deployment strategy.
    @property {String} ROLLINGPUSH - Rolling Push deployment strategy.
    @property {String} CUSTOM - Custom deployment strategy.
    @property {String} NONE - No deployment strategy.
**/
local StrategyNames = {
  REDBLACK: 'redblack',
  ROLLINGREDBLACK: 'rollingredblack',
  MONITORED: 'monitored',
  CFROLLINGREDBLACK: 'cfrollingredblack',
  HIGHLANDER: 'highlander',
  ROLLINGPUSH: 'rollingpush',
  CUSTOM: 'custom',
  NONE: 'none',
};

/**
    Red/Black (aka Blue/Green) deployment strategy configuration object.

    @example
        local redBlackStrat = deployment.RedBlackStrategy { delayBeforeDisableSec: '1', delayBeforeScaleDownSec: '1' }

        ...

        local myCluster = deployment.Cluster {
            ...
        } + redBlackStrat

    @property {Boolean} [rollbackFailure=false] - To rollback or not on failure.
    @property {String} [delayBeforeDisableSec='0'] - Delay before disabling the server group.
    @property {String} [delayBeforeScaleDownSec='0'] - Delay before sacling down the server group.
    @property {String} [maxRemainingAsgs='3'] - The maximum number of server groups to keep.
    @property {Boolean} [scaleDown=true] - Whether or not to scale down the previous server group.
**/
local RedBlackStrategy = {
  local this = self,

  rollbackFailure:: false,

  delayBeforeDisableSec: '0',
  delayBeforeScaleDownSec: '0',
  maxRemainingAsgs: '3',
  rollback: {
    onFailure: this.rollbackFailure,
  },
  scaleDown: true,
  strategy: StrategyNames.REDBLACK,
};

/**
    Spinnaker deployment Moniker.

    Formatted as such when used: app-stack-detail

    @example
        local myMoniker = deployment.Moniker { app: 'potatofactory', stack: 'dev', 'detail': 'producer'}

        ...

        local myCluster = deployment.Cluster {
            ...
            moniker: myMoniker,
            ...
        }

    @property {String} app - The app name.
    @property {String} stack - The stack.
    @property {String} [detail=''] - The Spinnaker Detail part of a moniker.
**/
local Moniker = {
  app: error '`app` (String) property is required for Moniker',
  stack: error '`stack` (String) property is required for Moniker',

  detail: '',
};

// Exposed for public use.
{
  Capacity:: Capacity,
  Moniker:: Moniker,
  RedBlackStrategy:: RedBlackStrategy,
  Rollback:: Rollback,
  StrategyNames:: StrategyNames,
}
