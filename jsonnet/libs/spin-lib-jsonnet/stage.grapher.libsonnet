/**
    @file Holds methods for connecting stages, in a graph-like manner.

    See `addRefIdsAndRequisiteRefIds` for more details.

    @author Santiago Gomez
**/

/**
    Makes an array of refIds from an array of stages.

    @author Santiago Gomez

    @param {Array} array - An array of stages.

    @returns {Array} An array of refIds from the given array of stages.
**/
local getRefIds(array) = std.map(function(object) (object.refId), array);

/**
    Gets the numerically last (greatest) refId from a given array of stages.

    @author Santiago Gomez

    @param {Array} stages - An array of stages.

    @returns {int} The last refId.
**/
local getLastRefId(stages) = std.foldl(std.max, getRefIds(stages), 0);

/**
    Checks if a given object is a stage.

    @author Santiago Gomez

    @param {Object} object - The object to check.

    @returns {Boolean} True if the object is a stage, false if otherwise.
**/
local isStage(object) = !std.isArray(object) && std.objectHas(object, 'name') && std.objectHas(object, 'type');

/**
    Checks if a given object is a stage.

    @author Santiago Gomez

    @param {Object} object  - The object to check.

    @returns {Boolean} True if the object is a stage, false if otherwise.
**/
local AddRefIdToObject(object, lastRefId, requisiteRefIds) =
  object {
    refId: lastRefId + 1,
    requisiteStageRefIds: if (requisiteRefIds == [] && std.objectHas(object, 'requisiteStageRefIds')) then object.requisiteStageRefIds else requisiteRefIds,
  };

/**
    An object used to keep track of the state/results while recursing.

    @author Santiago Gomez

    @param {Array} requisiteRefIds - An array of the "current" stage's required refIds.
    @param {Array} stages - Current array of stages, as a flat array.
    @param {int} lastRefId - The last refId of the current array of stages.
**/
local AccumulationResult(requisiteRefIds, stages, lastRefId=-1) = {
  requisiteRefIds: requisiteRefIds,
  lastRefId: if lastRefId == -1 then getLastRefId(stages) else lastRefId,
  stages: stages,
};

/**
    Merges two arrays of stages togather, producces a flat array of stages where their refIds don't overlap.

    @author Santiago Gomez

    @param {Array} first - First array of stages that will be merged.
    @param {Array} second - Second array of staegs that will be merged.

    @param {Array} A flat array of stages, merged from two arrays.
**/
local mergeIndependentAccumulationResults(first, second) =
  // shift all the stages based on the length of the previous parallel branch stages
  local shiftSecondRefIds(id) = if (std.member(getRefIds(second.stages), id)) then id + std.length(first.stages) else id;
  local shiftedStages = std.map(function(entry) entry { refId: shiftSecondRefIds(entry.refId), requisiteStageRefIds: std.map(shiftSecondRefIds, entry.requisiteStageRefIds) }, second.stages);
  // shift also the requisiteRefIds, since we've just changed all of their ids
  local shiftedSecondRequisiteRefIds = std.map(shiftSecondRefIds, second.requisiteRefIds);
  // return the merged result
  AccumulationResult(std.set(first.requisiteRefIds + shiftedSecondRequisiteRefIds), first.stages + shiftedStages);

/**
    Resurvively connects connects stages.

    @author Santiago Gomez

    More details/description see the `addRefIdsAndRequisiteRefIds` method which is the public interface.

    @param {AccumulationResult} previousAccumulationResult - The result being built up during recursion.
    @param {Object} object - Current object - Stage/Array/Parallel - that will be operated on.

    @returns {Array} A flat array of connected stages.
**/
local accumulateStages(previousAccumulationResult, object) =
  // if Stage, add refId and requisiteStageRefIds
  if isStage(object) then
    local stage = AddRefIdToObject(object, previousAccumulationResult.lastRefId, previousAccumulationResult.requisiteRefIds);
    AccumulationResult([stage.refId], previousAccumulationResult.stages + [stage])
  // if Array, iterate through all the items from left to right
  else if (std.isArray(object)) then
    std.foldl(accumulateStages, object, previousAccumulationResult)
  // if not an Array or a Stage, then it must be 2 or more parallel branches
  else if (std.objectHas(object, 'parallelStages')) && std.isArray(object.parallelStages) then
    if std.length(object.parallelStages) == 0 then
      previousAccumulationResult
    else
      // for each parallel branch, we accumulateStages, therefore creating some parallel and independent branches
      local parallelAccumulationResults = [std.foldl(accumulateStages, [stage], previousAccumulationResult) for stage in object.parallelStages];
      // each parallelAccumulationResult is duplicating the previousAccumulationResult.stages, so we must remove them by slicing the stages and removed what what's there before
      local parallelAccumulationResultsWithNoPreviousStages = [AccumulationResult(parallelAccumulationResult.requisiteRefIds, parallelAccumulationResult.stages[std.length(previousAccumulationResult.stages):]) for parallelAccumulationResult in parallelAccumulationResults];
      // merge the results into a single result, which means shifting the refIds of the stages so they don't overlap
      local mergedAccumulationResult = std.foldl(mergeIndependentAccumulationResults, parallelAccumulationResultsWithNoPreviousStages, AccumulationResult([], []));
      // put back all the previous stages onto the final result
      AccumulationResult(mergedAccumulationResult.requisiteRefIds, previousAccumulationResult.stages + mergedAccumulationResult.stages)
  else
    error 'object is not a Stage with name and type, nor an Array, nor has parallelStages';

/**
    Converts integer `refIds` and `requisiteStageRefIds` to Strings for a Stage - so that Spinnaker is happy.

    @author Santiago Gomez

    @param {Stage} stage - The stage for which to convert the refIds into strings.
**/
local fromNumbertoString(stage) = stage {
  refId: std.toString(stage.refId),
  requisiteStageRefIds: std.map(std.toString, stage.requisiteStageRefIds),
};

/**
	Adds RefIds and RequisiteStageRefIds to a series of stages.


	Stage dependencies go from left to right (StageA -> StageB -> StageC) based on the order in the array [StageA, StageB, StageC]


	If the array includes an object { x: StageX , y: StageY}, then all stages in the object are put at the end as independent dependencies
	[StageA, { x: StageX , y: StageY}, StageB] becomes (StageA -> {StageX, StageY} -> StageB)


	Also, remember that before the first stage, there's always an implicit Configuration stage:
	[{ x: StageX , y: StageY}] becomes (Configuration -> {StageX, StageY})


	See verify_stages.libsonnet for more examples


	*** Important: *** doesn't create a pipeline, just adds RefIds and RequisiteStageRefIds; you can combine it with any code that creates pipelines

	@example
	    local bakeStage = stage.BakeStage { ... }
	    local deployStage = stage.DeployStage { ... }

	    ...

	    local myStages = [ bakeStage, deployStage, ... ]

	    ...

	    pipeline.Pipeline {
	        ...
	        stages: stageGrapher.addRefIdsAndRequisiteRefIds(myStages),
	        ...
	    }

    @author Santiago Gomez

    @param {Array} stages - An array of Stage objects that will be connected.

    @returns {Array} An array of Stage objects that are now connected.
**/
local addRefIdsAndRequisiteRefIds(stages) = (
  std.map(fromNumbertoString, std.foldl(accumulateStages, stages, AccumulationResult([], [])).stages)
);

// Exposed for public use.
{
  addRefIdsAndRequisiteRefIds:: addRefIdsAndRequisiteRefIds,
}
