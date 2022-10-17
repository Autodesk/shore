/**
    @file Contains logging functions.
**/

/**
    Prints a given message.

    If a defValue is provided, it will be used as the output value for the std.trace.

    Otherwise the message will be set as the out value for the std.trace.

    @example
        local config = { test: "value" }
        local myTestConfig = tracing.log("Setting my test config to: "+config, config)
        {
            myConfig: tracing.log(myTestConfig),
        }

    @param {String} objName - Name of the object.
    @param {String} obj - Object to output the fields of.
    @param {String} defValue - Object to return from the trace. Default: null

    @returns {Object} - The default object passed in.
**/
local log(message, defValue=null) = (
  if defValue != null then
    std.trace(std.toString(message), defValue)
  else
    std.trace(std.toString(message), message)
);

/**
    Outputs the fields of a given object.

    @example
        local myString = "something" + tracing.traceFields("object name", myObject, " this will be appended");
        local myOtherObject = tracing.traceFields("object name", myObject);

    @param {String} objName - Name of the object.
    @param {String} obj - Object to output the fields of.
    @param {String} defValue - Object to return from the trace. Default: null

    @returns {Object} - The default object passed in.
**/
local traceFields(objName, obj, defValue=null) = (
  if defValue != null then
    std.trace(objName + ' fields are ' + std.manifestJsonEx(std.objectFieldsAll(obj), ' '), defValue)
  else
    std.trace(objName + ' fields are ' + std.manifestJsonEx(std.objectFieldsAll(obj), ' '), obj)
);


// Exposed for public use.
{
  log:: log,
  traceFields:: traceFields,
}
