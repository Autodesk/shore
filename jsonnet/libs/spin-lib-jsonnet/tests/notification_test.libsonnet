local notification = import '../notification.libsonnet';

// Data used for tests
local completeStageMessage = notification.StateMessage {
  message: 'Bake Potato stage completed running',
  state: notification.State.complete,
};
local startingStageMessage = notification.StateMessage {
  message: 'Bake Potato stage starting to running',
  state: notification.State.starting,
};
local manualJudgementAwaiting = notification.StateMessage {
  message: 'Waiting manual judgement decision',
  state: notification.State.manualJudgment,
};
local manualJudgementStopping = notification.StateMessage {
  message: 'Manual judgement decision was to stop',
  state: notification.State.manualJudgmentStop,
};

local tests = [
  notification.Notification {
    address: 'potato-channel',
    level: notification.Level.stage,
    type: notification.Type.SLACK,
    when: [notification.Level.stage + '.' + notification.State.complete],
  },
  completeStageMessage,
  notification.StageNotificationsConfiguration {
    notifications: [
      notification.Notification {
        address: 'potato-channel',
        level: notification.Level.stage,
        type: notification.Type.SLACK,
        when: [notification.Level.stage + '.' + notification.State.complete],
      },
    ],
  },
  notification.NewNotification {
    address: 'potato-channel',
    type: notification.Type.SLACK,
    Level:: notification.Level.stage,
    StateMessages:: [completeStageMessage, startingStageMessage],
  },
  notification.NewNotification {
    address: 'potato-channel',
    type: notification.Type.SLACK,
    Level:: notification.Level.stage,
    StateMessages:: [manualJudgementAwaiting, manualJudgementStopping],
  },
];

local assertions = [
  {
    address: 'potato-channel',
    level: 'stage',
    type: 'slack',
    when: [
      'stage.complete',
    ],
  },
  {
    message: 'Bake Potato stage completed running',
    state: 'complete',
  },
  {
    notifications: [
      {
        address: 'potato-channel',
        level: 'stage',
        type: 'slack',
        when: [
          'stage.complete',
        ],
      },
    ],
    sendNotifications: true,
  },
  {
    address: 'potato-channel',
    level: 'stage',
    message: {
      'stage.complete': {
        text: 'Bake Potato stage completed running',
      },
      'stage.starting': {
        text: 'Bake Potato stage starting to running',
      },
    },
    type: 'slack',
    when: [
      'stage.complete',
      'stage.starting',
    ],
  },
  {
    address: 'potato-channel',
    level: 'stage',
    message: {
      manualJudgment: {
        text: 'Waiting manual judgement decision',
      },
      manualJudgmentStop: {
        text: 'Manual judgement decision was to stop',
      },
    },
    type: 'slack',
    when: [
      'manualJudgment',
      'manualJudgmentStop',
    ],
  },
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
