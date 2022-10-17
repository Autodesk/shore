local stageGrapher = import '../stage.grapher.libsonnet';
local stage = import '../stage.libsonnet';

// this Stage type is just for testing, but it should work any stage type
local Stage = stage.Stage {
  name: 'stage-grapher-test',
  type: 'test',
};

local validateStage(stage) =
  if !std.objectHas(stage, 'refId') then '%s is missing refId' % [stage]
  else if !std.objectHas(stage, 'name') then '%s is missing name' % [stage]
  else if !std.objectHas(stage, 'requisiteStageRefIds') then '%s is missing requisiteStageRefIds' % [stage]
  else if !std.isArray(stage.requisiteStageRefIds) then '%s is requisiteStageRefIds is not an array' % [stage]
  else if std.member(stage.requisiteStageRefIds, stage.refId) then '%s is self-referencing' % [stage]
  else true;

local validateStages(stages) = [validateStage(stage) for stage in stages];

// checks if this stage has no dependencies that haven't been verified yet
local isLeaf(stage, removedLeafs=[]) =
  local notRemovedLeaf(x) = !std.member(removedLeafs, x);
  validateStage(stage) == true && std.filter(notRemovedLeaf, stage.requisiteStageRefIds) == [];

// checks that there are no cyclic dependencies (A <- B <- A) by means of logical pruning
// Spinnaker builds a dependency tree; if you prune the tree leaves long enough, you'll be be left with an empty array
// if you have (A <- B <- C), then the leaf is A. If you remove A, the leaf is B. If you remove B, the leaf is C. If you remove C, you have verified the tree
// if you have (A <- {B,C} <- D), the the leaf is A. If you remove A, the leafs are both B and C. If you remove them, then the leaf is D. If you remove D, you have verified the tree.
// if you have both (A <- B <- C) and (C <- B), the leaf is A. If you remove A, there are no leafs left), so the algorithm errors out.
local isAcyclic(stages) =
  local errorMessage(stage) = 'found a cycle on stage %s' % [stage];
  local foldFunction(removedLeafs, stage) = if std.isArray(removedLeafs) && isLeaf(stage, removedLeafs) then removedLeafs + [stage.refId] else errorMessage(stage);
  std.isArray(std.foldl(foldFunction, stages, []));

// checks that each element has a unique id and that they're ascending
local isAscendingRefId(stages) =
  local validateConsecutive(previous, next) =
    if std.objectHas(previous, 'refId') && previous.refId + 1 == next.refId then
      next
    else
      {};
  std.objectHas(std.foldl(validateConsecutive, stages, { refId: 0 }), 'refId');

// spinnaker expects strings, but our validation works on numbers
local fromStringtoNumber(stage) = stage {
  refId: std.parseInt(stage.refId),
  requisiteStageRefIds: std.map(std.parseInt, stage.requisiteStageRefIds),
};


// these test cases validate that my helper functions are actually finding issues and NOT just returning true
local badTestStages = {
  repeatingRefIds: [Stage { name: 'refId1', refId: 1 }, Stage { name: 'refId1Too', refId: 1 }],
  cyclic: [Stage { name: 'refId1', refId: 1, requisiteStageRefIds: [2] }, Stage { name: 'refId1Too', refId: 2, requisiteStageRefIds: [1] }],
};

local processedBadTestStages = [{
  testName: field,
  graphedStages: badTestStages[field],
  validStages: validateStages(self.graphedStages),
  ascendingRefIds: isAscendingRefId(self.graphedStages),
  isAcyclic: isAcyclic(self.graphedStages),
} for field in std.objectFields(badTestStages)];

// valid test cases should be added here to be processed down below
local goodTestStages = {
  local this = self,
  singleStage: [Stage { name: 'Apply' }],
  dependantStages: [Stage { name: 'Apply' }, Stage { name: 'Output' }],
  innerArray: [Stage { name: 'Plan' }, [Stage { name: 'Apply' }, Stage { name: 'Output' }]],
  parallelStages: [stage.Parallel { parallelStages: [Stage { name: 'Slack Notification' }, Stage { name: 'Email Notification' }] }],
  parallelInParallel: [Stage { name: 'Brace yourself' }, stage.Parallel { parallelStages: [this.dependantStages, this.parallelStages] }],
  mixedStages: [Stage { name: 'Predeploy infra' }, stage.Parallel { parallelStages: [[Stage { name: 'Deploy Canary' }, Stage { name: 'Wait for Canary' }], [Stage { name: 'Deploy Baseline' }, Stage { name: 'Wait for Baseline' }]] }, Stage { name: 'Verify Canary against Baseline' }],
  mixedStage2: [Stage { name: 'Deploy' }, Stage { name: 'Test' }, stage.Parallel { parallelStages: [Stage { name: 'Rollback' }, Stage { name: 'Manual check' }] }],
  emptyArrayInParallelStage: [stage.Parallel { parallelStages: [Stage { name: 'Slack Notification' }, []] }, Stage { name: 'Email Notification' }],
};

local processedGoodTestStages = [{
  testName: field,
  initialStages: goodTestStages[field],
  graphedStages: std.map(fromStringtoNumber, stageGrapher.addRefIdsAndRequisiteRefIds(goodTestStages[field])),
  validStages: validateStages(self.graphedStages),
  ascendingRefIds: isAscendingRefId(self.graphedStages),
  isAcyclic: isAcyclic(self.graphedStages),
} for field in std.objectFields(goodTestStages)];

local tests = processedGoodTestStages + processedBadTestStages;

local assertions = [
  import 'stage_grapher_results/dependantStages.json',
  import 'stage_grapher_results/emptyArrayInParallelStage.json',
  import 'stage_grapher_results/innerArray.json',
  import 'stage_grapher_results/mixedStage2.json',
  import 'stage_grapher_results/mixedStages.json',
  import 'stage_grapher_results/parallelInParallel.json',
  import 'stage_grapher_results/parallelStages.json',
  import 'stage_grapher_results/singleStage.json',
  import 'stage_grapher_results/cyclic.json',
  import 'stage_grapher_results/repeatingRefIds.json',
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
