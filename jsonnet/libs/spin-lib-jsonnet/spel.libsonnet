/**
    @file Contains commonly used SpEL expressions for Spinnaker.

    More details: {@link https://spinnaker.io/guides/user/pipeline/expressions/}
**/

/**
    Wraps in an expression in SpEL expression tags.

    @example
        local mySpelExpression = spel.expression("true == true");
        // Renders as: ${true == true}

    @param {String} spelExpression - SpEL expression to wrap.
**/
local expression(spelExpression) = (
  '${' + spelExpression + '}'
);

/**
    Part of a SpEL expression, provides a value from the current pipeline's execution context.

    {@link https://spinnaker.io/reference/pipeline/expressions/#helper-properties}

    @example
        local myExecStatusSpel = spel.expression( spel.execution('status') );
        // Renders as: ${execution["status"]}

    @param {String} key - A key/name of an object in the pipeline's execution context.
**/
local execution(key) = (
  'execution["' + key + '"]'
);

/**
    Part of a SpEL expression, retrieves the full stage.

    {@link https://spinnaker.io/reference/pipeline/expressions/#stagestring}

     @example
        local stageSpel = spel.expression( spel.stage('My Test Stage') );
        // Renders as: ${#stage('My Test Stage')}

    @param {String} stageName - Stage name of the stage to retrieve.
**/
local stage(stageName) = (
  "#stage('" + stageName + "')"
);

/**
    Part of a SpEL expression, serilizes an object into JSON.

    {@link https://spinnaker.io/reference/pipeline/expressions/#tojsonobject}

    @example
        local jsonizeObjectSpel = spel.expression( spel.toJson(spel.stage('My Test Stage')) );
        // Renders as: ${#toJson('#stage('My Test Stage')')}

    @param {String} spelObject - Object, from a SpEL expression.
**/
local toJson(spelObject) = (
  "#toJson('" + spelObject + "')"
);

/**
    Part of a SpEL expression, takes a string and encodes to Base64.

    {@link https://github.com/spinnaker/kork/blob/master/kork-expressions/src/main/java/com/netflix/spinnaker/kork/expressions/ExpressionsSupport.java#L256-L274}

    @example
        local parameterSpel = spel.expression( spel.toBase64("encode-me") );
        // Renders as: ${#toBase64("encode-me")}
        //
        // upon spel evaluation will turn to "ZW5jb2RlLW1lCg=="

    @param {String} string - the string to encode.
**/
local toBase64(spelObject) = (
  "#toBase64('" + spelObject + "')"
);

/**
    Part of a SpEL expression, takes a Base64 encoding and decodes it back to String.

    {@link https://github.com/spinnaker/kork/blob/master/kork-expressions/src/main/java/com/netflix/spinnaker/kork/expressions/ExpressionsSupport.java#L256-L274}

    @example
        local parameterSpel = spel.expression( spel.fromBase64("ZW5jb2RlLW1lCg==") );
        // Renders as: ${#fromBase64("ZW5jb2RlLW1lCg==")}
        //
        // upon spel evaluation will turn to "encode-me"

    @param {String} string - the Base64 string to decode.
**/
local fromBase64(spelObject) = (
  "#fromBase64('" + spelObject + "')"
);

/**
    Part of a SpEL expression, retrieves the parameter.

    {@link https://spinnaker.io/reference/pipeline/expressions/#helper-properties}

    @example
        local parameterSpel = spel.expression( spel.parameter("my-pipeline-param") );
        // Renders as: ${parameter["my-pipeline-param"]}

    @param {String} parameterName - Name of the parameter to retrieve.
**/
local parameter(parameterName) = (
  'parameters["' + parameterName + '"]'
);

/**
    Part of a SpEL expression, retrieves an artifact.

    {@link https://spinnaker.io/reference/pipeline/expressions/#helper-properties}

    @example
        local triggerArtifactSearchSpel = spel.expression( spel.getTriggerArtifact("my-artifact-name") );
        // Renders as: ${trigger['artifacts'].?[name == 'my-artifact-name']}

    @param {String} artifactName - Name of the artifact to retrieve.
**/
local getTriggerArtifact(artifactName) = (
  "trigger['artifacts'].?[name == '" + artifactName + "']"
);

/**
    Part of a SpEL expression, retrieves an artifact.

    {@link https://spinnaker.io/reference/pipeline/expressions/#helper-properties}

    @example
        local triggerArtifactSearchSpel = spel.expression( spel.getTriggerArtifactWithPrefix
("my-artifact-") );
        // Renders as: ${trigger['artifacts'].?[name.startsWith('my-artifact-')]}

    @param {String} prefix - The prefix of a the artifact name to retrieve.
**/
local getTriggerArtifactWithPrefix(prefix) = (
  "trigger['artifacts'].?[name.startsWith('%s')]" % prefix
);

/**
    Part of a SpEL expression, retrieves an artifact's reference value.

    Artifact class in Spinnaker code:

    {@link https://github.com/spinnaker/kork/blob/master/kork-artifacts/src/main/java/com/netflix/spinnaker/kork/artifacts/model/Artifact.java#L49}

    @example
        local triggerArtifactReferenceSpel = spel.expression( spel.getTriggerArtifactReference("my-artifact-name") );
        // Renders as: ${trigger['artifacts'].?[name == 'my-artifact-name'][0]['reference']}

    @param {String} artifactName - Name of the artifact to retrieve the refence value.
**/
local getTriggerArtifactReference(artifactName) = (
  getTriggerArtifact(artifactName) + "[0]['reference']"
);

/**
    Creates a part of a SpEL expression, provides the stage from the current execution.

    @example
        local triggerArtifactReferenceSpel = spel.expression( spel.executionStages("judgement") );
        // Renders as: ${execution.stages.?[name matches 'judgement']}

    @param {String} search - The stage name to match on.
**/
local executionStages(search) = (
  "execution.stages.?[name matches '" + search + "']"
);

/**
    Creates a part of a SpEL expression, creates a time-date stamp for the current time/date and plus time/date.

    Uses Java's SimpleDateFormat.

    @example
        local datetimeStampSpel = spel.expression( spel.simpleDateFormat("yyyyMMddhhmmss") );
        // Renders as: ${new java.text.SimpleDateFormat("yyyyMMddhhmmss").format(new java.util.Date())}
        local futureDatetimeStampSpel = spel.expression( spel.simpleDateFormat("yyyyMMddhhmmss", ".plusDays(1)") );
        // Renders as: ${new java.text.SimpleDateFormat("yyyyMMddhhmmss").format(new java.util.Date().plusDays(1))}

    @param {String} format - Format string. See Java's SimpleDateFormat.
    @param {String} [plusTime=""] - Used to append methods such as '.plusDays(1)'.

    @returns {String} SpEL expression that finds the deployed server group.
**/
local simpleDateFormat(format, plusTime='') = (
  'new java.text.SimpleDateFormat("' + format + '").format(new java.util.Date()' + plusTime + ')'
);

/**
    Creates a part of a SpEL expression, which converts a given object into a String.

    @example
        local datetimeStampSpel = spel.expression( spel.newString(spel.stage('My Test Stage')) );
        // Renders as: ${new String(#stage('My Test Stage'))}

    @param {String} obj - Object to convert into String.

    @returns {String} String representation a given object.
**/
local newString(obj) = (
  'new String(' + obj + ')'
);

// Exposed for public use.
{
  execution:: execution,
  executionStages:: executionStages,

  expression:: expression,

  newString:: newString,

  parameter:: parameter,

  stage:: stage,

  simpleDateFormat:: simpleDateFormat,

  toJson:: toJson,
  toBase64:: toBase64,
  fromBase64:: fromBase64,
  getTriggerArtifact:: getTriggerArtifact,
  getTriggerArtifactReference:: getTriggerArtifactReference,
  getTriggerArtifactWithPrefix:: getTriggerArtifactWithPrefix,
}
