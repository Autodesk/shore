/**
    @file Contains objects for creating Spinnaker Pipeline Triggers.
**/

/**
    Creates a basic trigger.

    This is used to extend into PipelineTrigger, JenkinsTrigger, WebhookTrigger.

    @example
        trigger.Trigger {
            type: trigger.Type.webhook
            source: 'my-webhook-trigger'
        }

    @property {Enum<Type>} type - The type of trigger. See the Type enum.
    @property {Boolean} [enabled=true] - Whether or not this trigger is enabled.
**/
local Trigger = {
  type: error '`type` (String) is a required property of `Trigger`',

  enabled: true,
  expectedArtifactIds: [],
};

/**
    A map of all Trigger types.

    This map was made from:
    {@link https://github.com/spinnaker/echo/blob/master/echo-model/src/main/java/com/netflix/spinnaker/echo/model/Trigger.java}

    @enum {String}
    @readonly

    @property {String} cron - Trigger of the type cron
    @property {String} git - Trigger of the type git
    @property {String} concourse - Trigger of the type concourse
    @property {String} jenkins - Trigger of the type jenkins
    @property {String} docker - Trigger of the type docker
    @property {String} webhook - Trigger of the type webhook
    @property {String} pubsub - Trigger of the type pubsub
    @property {String} dryrun - Trigger of the type dryrun
    @property {String} pipeline - Trigger of the type pipeline
    @property {String} plugin - Trigger of the type plugin
    @property {String} helm - Trigger of the type helm
**/
local Type = {
  cron: 'cron',
  git: 'git',
  concourse: 'concourse',
  jenkins: 'jenkins',
  docker: 'docker',
  webhook: 'webhook',
  pubsub: 'pubsub',
  dryrun: 'dryrun',
  pipeline: 'pipeline',
  plugin: 'plugin',
  helm: 'helm',
};

/**
    Creates a Pipeline trigger.

    Pipeline needs to be a Spinnaker Pipeline ID.

    @example
        trigger.PipelineTrigger {
            application: 'my-external-app'
            pipeline: 'my-other-pipeline'
        }

    @class
    @augments Trigger

    @property {String} application - The Spinnaker Appliaction of the Spinnaker Pipeline that will trigger this pipeline.
    @property {String} pipeline - The Spinnaker Pipeline that will trigger this pipeline.
    @property {Array<String>} [status=['successful']] - The acceptable status of a pipeline to trigger on.
**/
local PipelineTrigger = Trigger {
  application: error '`application` (String) is a required property of `PipelineTrigger`',
  pipeline: error '`pipeline` (String) is a required property of `PipelineTrigger`',

  type: Type.pipeline,
  status: [
    'successful',
  ],
};

/**
    Creates a Jenkins trigger.

    In the Spinnaker UI the "Controller" is the "master" property.

    @example
        trigger.JenkinsTrigger {
            job: 'my-jenkins-job'
            master: 'my-spinnaker-config-jenkins-server'
        }

    @class
    @augments Trigger

    @property {String} job - The Jenkins job that will trigger this pipeline.
    @property {String} master - The Jenkins server that has the Jenkins job.
    @property {String} [propertyFile=''] - The file provided by Jenkins when triggering a Spinnaker pipeline. This could be the build.properties file.
**/
local JenkinsTrigger = Trigger {
  job: error '`job` (String) is a required property of `JenkinsTrigger`',
  master: error '`master` (String) is a required property of `JenkinsTrigger`',

  propertyFile: '',
  type: Type.jenkins,
};

/**
    Creates a Webhook trigger.

    When created, the pipeline will be able to be trigger by an HTTP POST on:

   {@link https://<my.spinnaker.domain>/webhooks/webhook/<`source` property>}

    @example
        trigger.WebhookTrigger {
            source: 'trigger-my-pipeline-here'
            // URL will be: https://<my.spinnaker.domain>/webhooks/webhook/trigger-my-pipeline-here
        }

    @class
    @augments Trigger

    @property {String} source - The name of the webhook by which this pipeline can trigger.
    @property {Object} payloadConstraints - A map/dict of key-values that must be in the payload.
**/
local WebhookTrigger = Trigger {
  source: error '`source` (String) is a required property of `WebhookTrigger`',

  payloadConstraints: {},
  type: Type.webhook,
};

/**
    Creates an object to be used by NewTriggerByType.

    Generic enough to hold the bare minimum, and can have additional properties based on the type.

    `target` is what will be hit (jenkins job, pipeline name, webhook-source).

    `location` is where the target is (jenkins master, application name).

    @example
        trigger.GenericTrigger {
            type: 'pipeline',
            target: 'my-pipeline-name',
            location: 'my-app-name'
        }

        trigger.GenericTrigger {
            type: 'jenkin',
            target: 'my-jenkins-job',
            location: 'my-jenkins-master'
        }

        trigger.GenericTrigger {
            type: 'webhook',
            target: 'my-webhook',
        }

    @property {Enum<Type>} type - The type of trigger. See the Type enum.
    @property {String} target - The "target" to use.
**/
local GenericTrigger = {
  type: error '`type` (String) is a required property of `GenericTrigger`',

  target:: error '`target` (String) is a required property of `GenericTrigger`',

  location:: '',
};

// Exposed for public use.
{
  Trigger:: Trigger,
  Type:: Type,
  GenericTrigger:: GenericTrigger,

  JenkinsTrigger:: JenkinsTrigger,
  PipelineTrigger:: PipelineTrigger,
  WebhookTrigger:: WebhookTrigger,
}
