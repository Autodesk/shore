/**
    @file Contains objects for creating a Spinnaker Pipeline.
**/
local grapher = import './stage.grapher.libsonnet';

/**
    Creates a Spinnaker Pipeline object.

    @example
        local myConnectedStages =  [ ... ]

        ...

        pipeline.Pipeline {
            application: 'my-app',
            name: 'my-service deployment',
            stages: myConnectedStages
        }

        ------

        local myDisconnectedStages =  [ ... ]

        ...

        pipeline.Pipeline {
            application: 'my-app',
            name: 'my-service deployment',
            Stages:: myDisconnectedStages
        }

    @property {String} application - Name of the Spinnaker Application that this pipeline will be in.
    @property {String} name - Name of the pipeline.
    @property {Array<Parameter>} [parameterConfig=[]] - An array of pipeline parameters.
    @property {Array<Stage>} [Stages=[]] - An array of stages. Setting this one will connect them with stage.grapher.libsonnet.
    @property {Array<Stage>} [stages=[]] - An array of stages. Setting this one will not connect.
    @property {Array<ExpectedArtifact>} [expectedArtifacts=[]] - An array of artifacts that the pipeline exepects.
    @property {Boolean} [keepWaitingPipelines=false] - When set to true, allows pipelines to queue when pipeline concurancy is not enabled.
    @property {Boolean} [limitConcurrent=true] - When set to true, allows pipelines to run concurrantly.
    @property {Array<GenericTrigger>} [triggers=[]] - An array of triggers for the pipeline.

**/
local Pipeline = {
  application: error '`application` (String) property is required for Pipeline',
  name: error '`name` (String) property is required for Pipeline',

  Stages:: [],

  stages: if ($.Stages != null && std.length($.Stages) > 0) then grapher.addRefIdsAndRequisiteRefIds($.Stages) else [],

  expectedArtifacts: [],
  keepWaitingPipelines: false,
  limitConcurrent: true,
  parameterConfig: [],
  // TODO: Figure out how to make this auto-generated.
  spinLibJsonnetVersion: '0.0.0',
  triggers: [],
};


// Exposed for public use.
{
  Pipeline:: Pipeline,
}
